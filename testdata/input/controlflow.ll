; ModuleID = 'controlflow.go'
source_filename = "controlflow.go"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-wasi"

; Test if-else pattern
define i32 @max(i32 %a, i32 %b) {
entry:
  %cmp = icmp sgt i32 %a, %b
  br i1 %cmp, label %if.then, label %if.else

if.then:
  ret i32 %a

if.else:
  ret i32 %b
}

; Test while loop pattern
define i32 @sum_to_n(i32 %n) {
entry:
  %sum = alloca i32
  %i = alloca i32
  store i32 0, i32* %sum
  store i32 0, i32* %i
  br label %while.cond

while.cond:
  %i.val = load i32, i32* %i
  %cmp = icmp slt i32 %i.val, %n
  br i1 %cmp, label %while.body, label %while.end

while.body:
  %sum.val = load i32, i32* %sum
  %i.val2 = load i32, i32* %i
  %add = add i32 %sum.val, %i.val2
  store i32 %add, i32* %sum
  %inc = add i32 %i.val2, 1
  store i32 %inc, i32* %i
  br label %while.cond

while.end:
  %result = load i32, i32* %sum
  ret i32 %result
}

; Test simple if pattern (no else)
define void @print_positive(i32 %x) {
entry:
  %cmp = icmp sgt i32 %x, 0
  br i1 %cmp, label %if.then, label %if.end

if.then:
  call void @runtime.printint(i32 %x)
  br label %if.end

if.end:
  ret void
}

declare void @runtime.printint(i32)
