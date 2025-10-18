; ModuleID = 'forloop.go'
source_filename = "forloop.go"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-wasi"

; Test for loop pattern with initialization, condition, body, and increment
define i32 @sum_range(i32 %start, i32 %end) {
entry:
  br label %for.init

for.init:
  %i = alloca i32
  %sum = alloca i32
  store i32 %start, i32* %i
  store i32 0, i32* %sum
  br label %for.cond

for.cond:
  %i.val = load i32, i32* %i
  %cmp = icmp slt i32 %i.val, %end
  br i1 %cmp, label %for.body, label %for.end

for.body:
  %sum.val = load i32, i32* %sum
  %i.val2 = load i32, i32* %i
  %add = add i32 %sum.val, %i.val2
  store i32 %add, i32* %sum
  br label %for.inc

for.inc:
  %i.val3 = load i32, i32* %i
  %inc = add i32 %i.val3, 1
  store i32 %inc, i32* %i
  br label %for.cond

for.end:
  %result = load i32, i32* %sum
  ret i32 %result
}

; Test classic for loop: for i := 0; i < n; i++
define i32 @count_to_n(i32 %n) {
entry:
  br label %for.init

for.init:
  %i = alloca i32
  store i32 0, i32* %i
  br label %for.cond

for.cond:
  %i.val = load i32, i32* %i
  %cmp = icmp slt i32 %i.val, %n
  br i1 %cmp, label %for.body, label %for.end

for.body:
  ; Loop body - could have work here
  br label %for.inc

for.inc:
  %i.val2 = load i32, i32* %i
  %inc = add i32 %i.val2, 1
  store i32 %inc, i32* %i
  br label %for.cond

for.end:
  %result = load i32, i32* %i
  ret i32 %result
}
