;; Run this AI with: csi -s scout.scm ID
;; where ID is the number affected to the ant you want this IA to control

;; Be sure that loops and matchable librairies are installed
;; (e.g. with `[sudo] chicken-install loops matchable list-utils`)

;; This AI is meant to explore the map
;; If it encounter a enemy ant, it will run away from it
;; If possible, this ant avoid going on a cell it knows

;; Note: global variables names are UPPERCASE.

;; Libraries import
(use loops matchable list-utils)

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

;; We use it **a lot**
(define li list)

(define (<> a b) (not (= a b) ) )

;; Read [A] lines and fill the [ANTS] list reading next lines
;; [ANT] associate an ID to a list of parameters.
(define (read-ants)
  ( set! ANTS (li) ) ;; reset ANTS
  ( do-times
    _ A (match-let (( (id x y dx dy e a b) (input-int-line) ))
                   (set! ANTS (append ANTS (li (li id (li x y dx dy e a b) ) )
                                      ) ) ) ) )

;; Read the header line (the number of lines to parse)
;; Fill the [ENEMIES] list.
(define (read-enemies)
  (let ((e (string->number (read-line) ) ))
    ( set! ENEMIES '() ) ;; reset ENEMIES
    ( do-times _ e (set! ENEMIES (li (input-int-line) ENEMIES) ) ) ) )

;; Read the header line containing [W H N] and then
;; read the next lines according the header,
;; filling [MAP].
(define (read-map)
  (match-let (( (w h n) (input-int-line) ))
             ( set! MAP (li) ) ;; reset MAP
             ( do-times
               _ n (match-let (( (x y c s) (input-int-line) ))
                              ( set! MAP (append (li (li (li x y)
                                                         (li c s) ) )
                                                 MAP ) ) ) ) ) )

;; Helpers to test a cell content
(define (rock? c) (= 2 c))
(define (water? c) (= 4 c))

;; Helpers to write an ant move
(define (action a act) (string-append (number->string a) ":" act) )
(define (rest a) (action a "rest") )
(define (forward a) (action a "forward") )
(define (right a) (action a "right") )
(define (left a) (action a "left") )

;; Return the cell fetched at [(x, y)]
;; or [#f] if no cell is known here
(define (cell x y) (assoc-def (li x y) MAP equal? #f) )

(define (ant a) (cadr (assoc a ANTS equal?) ) )

;; [#f] if a cell pointed by [x y] is already known
;; list-utils is needed because primitive key comparison is not
;; working with key which are pairs
(define (unknown-cell? x y) (not (cell x y) ) )

;; Test if the cell at [(x, y)] is safe or not
;; FIXME: should also test if an ant is already on this cell or not
(define (walkable-cell? x y)
  (match-let (( (_ (c _)) (cell x y) ))
             (not (or (rock? c) (water? c) ) ) ) )

;; For a given [ant] ant, choose the next move.
;; FIXME: This function will keep an ant turning right / left
;; without moving if it come to a point where it can not find any
;; unknown cell close.
(define (choose-move antid)
  (match-let (( (x y dx dy e a b) (ant antid) ))
             (cond

              ;; brain is not controlled: send dummy [rest] instruction?
              ((<> 1 b) (rest antid))

              ;; two cells forward is an unkown cell and the forward cell
              ;; can be walked on: go [forward]
              ((and (unknown-cell? (+ (+ x dx) dx) (+ (+ y dy) dy) )
                    (walkable-cell? (+ x dx) (+ y dy) ) )
               (forward antid) )

              ;; Randomly change the ant orientation if forward is not
              ;; a suitable move
              (else (if (= 0 (random 2)) (left antid) (right antid) ) )
              ) ) )

;;;; Main loop

(define ANT-ID (string->number (car (command-line-arguments) ) ) )

(do-forever (read-header)
            (if (= S 0) (exit) '())
            (read-ants)
            (read-enemies)
            (read-map)
            (print (choose-move ANT-ID))
            )
