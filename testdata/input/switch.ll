; ModuleID = 'switch.go'
source_filename = "switch.go"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-wasi"

define i32 @testSwitch(i32 %x) {
entry:
  switch i32 %x, label %default [
    i32 0, label %case0
    i32 1, label %case1
    i32 2, label %case2
  ]

case0:
  ret i32 10

case1:
  ret i32 20

case2:
  ret i32 30

default:
  ret i32 -1
}

define void @main() {
entry:
  %0 = call i32 @testSwitch(i32 1)
  call void @runtime.printint(i32 %0)
  ret void
}

declare void @runtime.printint(i32)
