;; Be sure that readline, loops and matchable librairies are installed
;; (e.g. with `[sudo] chicken-install readline loops matchable`)

;; Run this AI with: csi scout.scm

;; This AI is meant to explore the map
;; If it encounter a enemy ant, it will run away from it
;; If possible, this ant avoid going on a cell it knows

;; Note: global variables names are UPPERCASE.

;; Libraries import
(use readline loops matchable)

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

;; Read [A] lines and fill the [ANTS] list reading next lines
;; [ANT] associate an ID to a list of parameters.
(define (read-ants)
  ( set! ANTS '() ) ;; reset ANTS
  ( do-times _ A (match-let (( (id x y dx dy e a b) (input-int-line) ))
                            ( set! ANTS (cons (list id (list x y dx dy e a b))
                                              ANTS ) ) ) ) )

;; Read the header line containing [W H N] and then
;; read the next lines according the header,
;; filling [MAP].
(define (read-map)
  (match-let (( (w h n) (input-int-line) ))
             ( set! MAP '() ) ;; reset MAP
             ( do-times _ n (set! MAP (cons (input-int-line) MAP) ) ) ) )

;;;; Main loop

(do-forever (read-header)
            (if (= S 0) (exit) (print "turn:" T
                                      " ant/player: " A
                                      " players: " P
                                      " status: " S))
            (read-ants)
            (print "ants:" ANTS)
            (read-map)
            (print "map:" MAP)
            )
