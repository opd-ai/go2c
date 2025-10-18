; Test for multiple return values represented as a struct
; This simulates what TinyGo would generate for: func foo(x int) (int, int)

define {i32, i32} @foo(i32 %x) {
entry:
  %result = insertvalue {i32, i32} undef, i32 %x, 0
  %result2 = insertvalue {i32, i32} %result, i32 42, 1
  ret {i32, i32} %result2
}

define void @main() {
entry:
  %0 = call {i32, i32} @foo(i32 10)
  %val1 = extractvalue {i32, i32} %0, 0
  %val2 = extractvalue {i32, i32} %0, 1
  call void @runtime.printint(i32 %val1)
  call void @runtime.printint(i32 %val2)
  ret void
}

declare void @runtime.printint(i32)

define i32 @_start() {
entry:
  call void @main()
  ret i32 0
}
