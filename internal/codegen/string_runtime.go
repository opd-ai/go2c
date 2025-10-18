package codegen

// GetStringRuntimeLibrary returns the C code for the string runtime library
// This library provides string operations compatible with Go string semantics
func GetStringRuntimeLibrary() string {
	return `// ===== Go String Runtime Library =====
// Provides string operations compatible with Go string semantics
// This library is automatically included when string operations are detected

typedef struct {
    const char* data;  // Pointer to string data
    size_t len;        // Length of string (not including null terminator)
} go_string_t;

// Create a new go_string from a C string
static inline go_string_t go_string_new(const char* cstr) {
    go_string_t s;
    s.data = cstr;
    s.len = strlen(cstr);
    return s;
}

// Create a go_string from bytes with explicit length
static inline go_string_t go_string_from_bytes(const char* data, size_t len) {
    go_string_t s;
    s.data = data;
    s.len = len;
    return s;
}

// Compare two go_strings (returns 0 if equal, <0 if a<b, >0 if a>b)
static inline int go_string_compare(go_string_t a, go_string_t b) {
    if (a.len != b.len) {
        return (int)(a.len - b.len);
    }
    return memcmp(a.data, b.data, a.len);
}

// Check if two go_strings are equal
static inline bool go_string_equals(go_string_t a, go_string_t b) {
    if (a.len != b.len) {
        return false;
    }
    return memcmp(a.data, b.data, a.len) == 0;
}

// Concatenate two go_strings (allocates new memory)
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

// Extract substring (allocates new memory)
// Parameters: s - source string, start - start index, end - end index (exclusive)
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

// Convert go_string to null-terminated C string (allocates new memory)
// Caller is responsible for freeing the returned memory
static inline char* go_string_to_cstr(go_string_t s) {
    char* cstr = (char*)malloc(s.len + 1);
    if (cstr == NULL) {
        return NULL;
    }
    memcpy(cstr, s.data, s.len);
    cstr[s.len] = '\0';
    return cstr;
}

// Free string memory (only for strings created with concat/substring/to_cstr)
// Do NOT use this for static string constants
static inline void go_string_free(char* data) {
    if (data != NULL) {
        free(data);
    }
}

// Get character at index (no bounds checking)
static inline char go_string_at(go_string_t s, size_t index) {
    return s.data[index];
}

// Check if string contains a substring
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

`
}
