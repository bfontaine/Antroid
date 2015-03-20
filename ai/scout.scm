#!/usr/bin/csi

;; Be sure that readline, loops and matchable librairies are installed
;; (e.g. with `[sudo] chicken-install readline loops matchable list-utils`)

;; Run this AI with: csi scout.scm

;; This AI is meant to explore the map
;; If it encounter a enemy ant, it will run away from it
;; If possible, this ant avoid going on a cell it knows

;; Note: global variables names are UPPERCASE.

;; Libraries import
(use readline loops matchable list-utils)

;; Read a line on stdin, split using ' ' as delim
(define (input-line)
  (string-split (read-line) ) )

;; Same as input-line, but also cast resulting list to int list
(define (input-int-line)
  (map (lambda (x) (string->number x))
       (string-split (read-line) ) ) )

;; Read the first line of server input
;; [T, A, P and S] are global variables used be other reading functions
(define (read-header)
  (match-let (( (t a p s) (input-int-line) ))
             (set! T t) ;; Turn number
             (set! A a) ;; Ant / player
             (set! P p) ;; Number of players
             (set! S s) ;; Game status. 1: playing, 0: over
             )
  )

(define (<> a b) (not (= a b) ) )

;; Read [A] lines and fill the [ANTS] list reading next lines
;; [ANT] associate an ID to a list of parameters.
(define (read-ants)
  ( set! ANTS '() ) ;; reset ANTS
  ( do-times _ A (match-let (( (id x y dx dy e a b) (input-int-line) ))
                            ( set! ANTS (list (list id (list x y dx dy e a b))
                                              ANTS ) ) ) ) )

;; Read the header line (the number of lines to parse)
;; Fill the [ENEMIES] list.
(define (read-enemies)
  (let ((e (string->number (read-line) ) ))
    ( set! ENEMIES '() ) ;; reset ENEMIES
    ( do-times _ e (set! ENEMIES (list (input-int-line) ENEMIES) ) ) ) )

;; Read the header line containing [W H N] and then
;; read the next lines according the header,
;; filling [MAP].
(define (read-map)
  (match-let (( (w h n) (input-int-line) ))
             ( set! MAP '() ) ;; reset MAP
             ( do-times _ n (match-let (( (x y c s) (input-int-line) ))
                                       ( set! MAP (list (list (list x y)
                                                              (list c s) )
                                                        MAP ) ) ) ) ) )

;; Helpers to test a cell content
(define (rock? c) (= 2))
(define (water? c) (= 4))

;; Helpers to write an ant move
(define (rest a) (print a ":rest") )
(define (forward a) (print a ":forward") )
(define (right a) (print a ":right") )
(define (left a) (print a ":left") )

;; Return the cell fetched at [(x, y)]
;; or [#f] if no cell is known here
(define (cell x y) (assoc-def (list x y) MAP equal? #f) )

;; [#f] if a cell pointed by [x y] is already known
;; list-utils is needed because primitive key comparison is not
;; working with key which are pairs
(define (unknown-cell? x y) (not (cell x y) ) )

;; Test if the cell at [(x, y)] is safe or not
;; FIXME: should also test if an ant is already on this cell or not
(define (walkable-cell? x y)
  (match-let (( (c _) (cell x y) )) (not (or (rock? c) (water? c) ) ) ) )

;; For a given [ant] ant, choose the next move.
;; FIXME: This function will keep an ant turning right / left
;; without moving if it come to a point where it can not find any
;; unknown cell close.
(define (choose-move ant)
  (match-let (( (x y dx dy e a b) (cadr (assq ant ANTS) ) ))
             (cond

              ;; brain is not controlled: send dummy [rest] instruction?
              ((<> 1 b) => (rest ant))

              ;; two cells forward is an unkown cell and the forward cell
              ;; can be walked on: go [forward]
              ((and (unknown-cell? (+ (+ x dx) dx) (+ (+ y dy) dy) )
                    (walkable-cell? (+ x dx) (+ y dy) ) )
               => (forward ant) )

              ;; Randomly change the ant orientation if forward is not
              ;; a suitable move
              (else => (if (= 0 (random 2)) (left a) (right a) ) )
              ) ) )

;;;; Main loop

(do-forever (read-header)
            (if (= S 0) (exit) (print "turn:" T
                                      " ant/player: " A
                                      " players: " P
                                      " status: " S))
            (read-enemies)
            (print "enemies:" ENEMIES)
            (read-ants)
            (print "ants:" ANTS)
            (read-map)
            (print "map:" MAP)
            (print "moves:")
            (choose-move 4)
            )
