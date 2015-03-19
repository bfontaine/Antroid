;; Be sure that readline, loops and matchable librairies are installed
;; (e.g. with `[sudo] chicken-install readline loops matchable`)

;; Run this AI with: csi scout.scm

;; This AI is meant to explore the map
;; If it encounter a enemy ant, it will run away from it
;; If possible, this ant avoid going on a cell it knows

;; Libraries import
(use readline loops matchable)

;; Read the first line of server input
;; [T, A, P and S] are global variables used be other reading functions
(define (read-header)
  (match-let (( (t a p s) (string-split (read-line) ) ))
             (set! T (string->number t)) ;; Turn number
             (set! A (string->number a)) ;; Ant / player
             (set! P (string->number p)) ;; Number of players
             (set! S (string->number s)) ;; Game status. 1: playing, 0: over
             )
  )

;;;; Main loop

(define S 1) ;; initialize game status to 'playing' for cold start.

(do-until (= S 0)
          (read-header)
          (print T A P S)
          )

(exit)
