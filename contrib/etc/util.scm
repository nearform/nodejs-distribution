(define set-major-minor (lambda (v)
  (let ((versions (string-split v #\.)))
    (gmk-eval (string-join (list "MAJOR:=" (car versions))))
    (gmk-eval (string-join (list "MINOR:=" (cdr versions)))))))
