: fib dup 1 > if dup 2 - fib swap 1 - fib + then ;
: print-fib-numbers 10 0 do i fib . loop ;
print-fib-numbers
