;; Libraries import
(use readline loops matchable)

;; Read the first line of server input
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
