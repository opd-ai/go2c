package codegen

// GetRuntimeLibrary returns the C code for the runtime library
// This library provides implementations of TinyGo/Go runtime functions
// that are commonly called from generated LLVM IR
func GetRuntimeLibrary() string {
	return `// ===== Go Runtime Library =====
// Provides implementations of common Go runtime functions
// This library is automatically included when runtime functions are detected

// runtime.printint - Print an integer to stdout
static inline void runtime_printint(int32_t value) {
    printf("%d\n", value);
}

// runtime.printstring - Print a string with explicit length to stdout
// This matches Go's string semantics (data pointer + length)
static inline void runtime_printstring(const char* data, int32_t len) {
    if (data != NULL && len > 0) {
        fwrite(data, 1, len, stdout);
    }
}

// runtime.printuint32 - Print an unsigned 32-bit integer
static inline void runtime_printuint32(uint32_t value) {
    printf("%u\n", value);
}

// runtime.printuint64 - Print an unsigned 64-bit integer
static inline void runtime_printuint64(uint64_t value) {
    printf("%llu\n", (unsigned long long)value);
}

// runtime.printint64 - Print a signed 64-bit integer
static inline void runtime_printint64(int64_t value) {
    printf("%lld\n", (long long)value);
}

// runtime.printfloat32 - Print a 32-bit float
static inline void runtime_printfloat32(float value) {
    printf("%f\n", value);
}

// runtime.printfloat64 - Print a 64-bit float (double)
static inline void runtime_printfloat64(double value) {
    printf("%f\n", value);
}

// runtime.printbool - Print a boolean value
static inline void runtime_printbool(bool value) {
    printf("%s\n", value ? "true" : "false");
}

// runtime.printpointer - Print a pointer address
static inline void runtime_printpointer(void* ptr) {
    printf("%p\n", ptr);
}

// runtime.printnl - Print a newline
static inline void runtime_printnl(void) {
    putchar('\n');
}

// runtime.printspace - Print a space
static inline void runtime_printspace(void) {
    putchar(' ');
}

// runtime.alloc - Simple memory allocation wrapper
// Note: In real Go, this would interact with the garbage collector
static inline void* runtime_alloc(size_t size) {
    return malloc(size);
}

// runtime.free - Memory deallocation (not typically in Go, but useful for C)
static inline void runtime_free(void* ptr) {
    if (ptr != NULL) {
        free(ptr);
    }
}

// runtime.memcpy - Memory copy wrapper
static inline void* runtime_memcpy(void* dest, const void* src, size_t n) {
    return memcpy(dest, src, n);
}

// runtime.memset - Memory set wrapper
static inline void* runtime_memset(void* s, int c, size_t n) {
    return memset(s, c, n);
}

// runtime.slicecopy - Copy slice data (simplified version)
// In real Go, this handles slice metadata; here we just copy memory
static inline int runtime_slicecopy(void* dst, void* src, size_t len, size_t elemsize) {
    size_t bytes = len * elemsize;
    memcpy(dst, src, bytes);
    return (int)len;
}

// runtime.strcmp - String comparison for C strings
static inline int runtime_strcmp(const char* a, const char* b) {
    return strcmp(a, b);
}

// runtime.strlen - String length for C strings
static inline size_t runtime_strlen(const char* s) {
    return strlen(s);
}

// ===== End of Go Runtime Library =====

`
}
