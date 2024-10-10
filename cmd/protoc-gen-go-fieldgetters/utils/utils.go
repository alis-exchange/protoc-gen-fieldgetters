package utils

import "strings"

// Contains function is a nifty tool! It returns true if it finds the value you're looking for in the array.
// And guess what? It can handle any type, all thanks to the power of Go's Generics!
func Contains[T comparable](s []T, searchTerm T) bool {
	for _, item := range s {
		if item == searchTerm {
			return true
		}
	}
	return false
}

// Transform is a utility function for transforming the elements of a slice.
//
// For example, you can use it to double the values of a slice of integers:
//
//	ints := []int{1, 2, 3}
//	doubled := Transform(ints, func(i int) int { return i * 2 })
//	// doubled = [2, 4, 6]
func Transform[T, U any](arr []T, fn func(T) U) []U {
	result := make([]U, len(arr))
	for i, v := range arr {
		result[i] = fn(v)
	}
	return result
}

// ToLowerFirst transforms the first letter of a string to lowercase
func ToLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func ToUpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func SnakeCaseToCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}

	words := strings.Split(s, "_")
	if len(words) == 1 {
		return ToLowerFirst(words[0])
	}

	var b strings.Builder
	for i, word := range words {
		if i == 0 {
			b.WriteString(strings.ToLower(word))
			continue
		}

		b.WriteString(ToUpperFirst(strings.ToLower(word)))
	}
	return b.String()
}
