// Example: Using the go2c string runtime library
// This demonstrates the string operations available in generated C code

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdbool.h>

// ===== Go String Runtime Library =====
// (This is automatically included when strings are detected)

typedef struct {
    const char* data;
    size_t len;
} go_string_t;

static inline go_string_t go_string_new(const char* cstr) {
    go_string_t s;
    s.data = cstr;
    s.len = strlen(cstr);
    return s;
}

static inline go_string_t go_string_from_bytes(const char* data, size_t len) {
    go_string_t s;
    s.data = data;
    s.len = len;
    return s;
}

static inline int go_string_compare(go_string_t a, go_string_t b) {
    if (a.len != b.len) {
        return (int)(a.len - b.len);
    }
    return memcmp(a.data, b.data, a.len);
}

static inline bool go_string_equals(go_string_t a, go_string_t b) {
    if (a.len != b.len) {
        return false;
    }
    return memcmp(a.data, b.data, a.len) == 0;
}

static inline go_string_t go_string_concat(go_string_t a, go_string_t b) {
    char* new_data = (char*)malloc(a.len + b.len + 1);
    if (new_data == NULL) {
        go_string_t empty = {NULL, 0};
        return empty;
    }
    memcpy(new_data, a.data, a.len);
    memcpy(new_data + a.len, b.data, b.len);
    new_data[a.len + b.len] = '\0';
    
    go_string_t result;
    result.data = new_data;
    result.len = a.len + b.len;
    return result;
}

static inline go_string_t go_string_substring(go_string_t s, size_t start, size_t end) {
    if (start > end || end > s.len) {
        go_string_t empty = {NULL, 0};
        return empty;
    }
    
    size_t sub_len = end - start;
    char* new_data = (char*)malloc(sub_len + 1);
    if (new_data == NULL) {
        go_string_t empty = {NULL, 0};
        return empty;
    }
    
    memcpy(new_data, s.data + start, sub_len);
    new_data[sub_len] = '\0';
    
    go_string_t result;
    result.data = new_data;
    result.len = sub_len;
    return result;
}

static inline char* go_string_to_cstr(go_string_t s) {
    char* cstr = (char*)malloc(s.len + 1);
    if (cstr == NULL) {
        return NULL;
    }
    memcpy(cstr, s.data, s.len);
    cstr[s.len] = '\0';
    return cstr;
}

static inline void go_string_free(char* data) {
    if (data != NULL) {
        free(data);
    }
}

static inline char go_string_at(go_string_t s, size_t index) {
    return s.data[index];
}

static inline bool go_string_contains(go_string_t s, go_string_t substr) {
    if (substr.len > s.len) {
        return false;
    }
    
    for (size_t i = 0; i <= s.len - substr.len; i++) {
        if (memcmp(s.data + i, substr.data, substr.len) == 0) {
            return true;
        }
    }
    return false;
}

// ===== End of Go String Runtime Library =====

// Example usage
int main() {
    printf("=== Go String Operations Example ===\n\n");
    
    // 1. Create strings from C strings
    printf("1. Creating strings:\n");
    go_string_t hello = go_string_new("Hello");
    go_string_t world = go_string_new("World");
    printf("   hello: \"%.*s\" (len=%zu)\n", (int)hello.len, hello.data, hello.len);
    printf("   world: \"%.*s\" (len=%zu)\n\n", (int)world.len, world.data, world.len);
    
    // 2. String concatenation
    printf("2. String concatenation:\n");
    go_string_t hello_space = go_string_new("Hello ");
    go_string_t message = go_string_concat(hello_space, world);
    printf("   \"Hello \" + \"World\" = \"%.*s\"\n\n", (int)message.len, message.data);
    
    // 3. String comparison
    printf("3. String comparison:\n");
    go_string_t hello2 = go_string_new("Hello");
    printf("   \"Hello\" == \"Hello\": %s\n", go_string_equals(hello, hello2) ? "true" : "false");
    printf("   \"Hello\" == \"World\": %s\n", go_string_equals(hello, world) ? "true" : "false");
    int cmp = go_string_compare(hello, world);
    printf("   compare(\"Hello\", \"World\"): %d (negative = less than)\n\n", cmp);
    
    // 4. Substring extraction
    printf("4. Substring extraction:\n");
    go_string_t full = go_string_new("Hello, World!");
    go_string_t sub1 = go_string_substring(full, 0, 5);    // "Hello"
    go_string_t sub2 = go_string_substring(full, 7, 12);   // "World"
    printf("   \"%.*s\"[0:5] = \"%.*s\"\n", (int)full.len, full.data, (int)sub1.len, sub1.data);
    printf("   \"%.*s\"[7:12] = \"%.*s\"\n\n", (int)full.len, full.data, (int)sub2.len, sub2.data);
    
    // 5. Substring search
    printf("5. Substring search:\n");
    go_string_t search = go_string_new("World");
    printf("   \"%.*s\" contains \"World\": %s\n", (int)full.len, full.data,
           go_string_contains(full, search) ? "true" : "false");
    go_string_t not_found = go_string_new("xyz");
    printf("   \"%.*s\" contains \"xyz\": %s\n\n", (int)full.len, full.data,
           go_string_contains(full, not_found) ? "true" : "false");
    
    // 6. Character access
    printf("6. Character access:\n");
    for (size_t i = 0; i < hello.len; i++) {
        printf("   hello[%zu] = '%c'\n", i, go_string_at(hello, i));
    }
    printf("\n");
    
    // 7. Binary-safe strings (with embedded null bytes)
    printf("7. Binary-safe strings:\n");
    const char binary_data[] = {'H', 'e', 'l', '\0', 'l', 'o'};
    go_string_t binary = go_string_from_bytes(binary_data, 6);
    printf("   Binary string with null byte: len=%zu\n", binary.len);
    printf("   Data (hex): ");
    for (size_t i = 0; i < binary.len; i++) {
        printf("%02X ", (unsigned char)go_string_at(binary, i));
    }
    printf("\n\n");
    
    // 8. Convert to C string
    printf("8. Convert to C string:\n");
    char* cstr = go_string_to_cstr(message);
    printf("   go_string -> C string: \"%s\"\n\n", cstr);
    
    // Cleanup dynamically allocated strings
    go_string_free((char*)message.data);
    go_string_free((char*)sub1.data);
    go_string_free((char*)sub2.data);
    go_string_free(cstr);
    
    printf("=== Example complete ===\n");
    return 0;
}
