;; I used "printed-buffer-chunk" as an example to get this working,
;; but I have not been able to do exactly what I want.

;; Right now this produces something like:

;;   goal: GOAL-CHUNK0
;;     NUM1  3
;;     NUM2  1
;;     COUNT  EMPTY
;;     SUM  EMPTY

;; What I want is more like this:

;;   goal: add(count=empty, num1=3, num2=1, sum=empty)

;; 1. I don't know how to get the "isa" from the chunk ("GOAL-CHUNK0" -> "add" in this case).

;; 2. Displaying each of the slots is actually happening in "printed-chunk", but I don't
;; understand Lisp enough to pick it apart & format the way I want it.

(defun vanilla-print-buffer (&rest buffer-names-list)
   (verify-current-model
    "vanilla-print-buffer called with no current model."
    (let ((s (make-string-output-stream)))
      (dolist (buffer-name (if buffer-names-list
                               buffer-names-list
                             (model-buffers)))
        (let ((buffer (buffer-instance buffer-name)))
          (when buffer
            (bt:with-recursive-lock-held ((act-r-buffer-lock buffer))
              (let ((chunk (act-r-buffer-chunk buffer)))
                (when buffer-names-list
                  (format s "~a" (printed-chunk chunk))))))))
      (get-output-stream-string s))))