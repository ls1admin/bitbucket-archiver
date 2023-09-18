package utils

// function that creates chunks of length n from a slice of strings
func Chunks[S any](s []S, n int) [][]S {
	var chunks [][]S
	for n < len(s) {
		s, chunks = s[n:], append(chunks, s[0:n:n])
	}
	return append(chunks, s)
}
