; ModuleID = 'hello.go'
source_filename = "hello.go"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-wasi"

@main_string = internal constant [14 x i8] c"Hello, World!\00", align 1

declare void @runtime_printstring(i8*, i32)

define void @main() {
entry:
  call void @runtime_printstring(i8* @main_string, i32 13)
  ret void
}

define i32 @_start() {
entry:
  call void @main()
  ret i32 0
}
