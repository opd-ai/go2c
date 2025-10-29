; ModuleID = 'dowhile.go'
source_filename = "dowhile.go"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-wasi"

; Test do-while loop pattern
; C equivalent:
; int sum = 0;
; int i = 0;
; do {
;     sum += i;
;     i++;
; } while (i < n);
; return sum;
define i32 @sum_do_while(i32 %n) {
entry:
  %sum = alloca i32
  %i = alloca i32
  store i32 0, i32* %sum
  store i32 0, i32* %i
  br label %do.body

do.body:
  %sum.val = load i32, i32* %sum
  %i.val = load i32, i32* %i
  %add = add i32 %sum.val, %i.val
  store i32 %add, i32* %sum
  %inc = add i32 %i.val, 1
  store i32 %inc, i32* %i
  br label %do.cond

do.cond:
  %i.val2 = load i32, i32* %i
  %cmp = icmp slt i32 %i.val2, %n
  br i1 %cmp, label %do.body, label %do.end

do.end:
  %result = load i32, i32* %sum
  ret i32 %result
}

; Test simple do-while with single statement
define void @print_countdown(i32 %n) {
entry:
  %counter = alloca i32
  store i32 %n, i32* %counter
  br label %do.body

do.body:
  %val = load i32, i32* %counter
  call void @runtime.printint(i32 %val)
  %dec = sub i32 %val, 1
  store i32 %dec, i32* %counter
  br label %do.cond

do.cond:
  %val2 = load i32, i32* %counter
  %cmp = icmp sgt i32 %val2, 0
  br i1 %cmp, label %do.body, label %do.end

do.end:
  ret void
}

declare void @runtime.printint(i32)
