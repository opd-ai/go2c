; ModuleID = 'simple.go'
source_filename = "simple.go"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-wasi"

define i32 @add(i32 %a, i32 %b) {
entry:
  %result = add i32 %a, %b
  ret i32 %result
}

define void @main() {
entry:
  %0 = call i32 @add(i32 5, i32 3)
  call void @runtime.printint(i32 %0)
  ret void
}

declare void @runtime.printint(i32)

define i32 @_start() {
entry:
  call void @main()
  ret i32 0
}
