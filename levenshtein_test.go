package editdistance

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var testCases = []struct {
	source          string
	target          string
	distance        int
	damerauDistance int
	lcsDistance     int
}{
	{"", "a", 1, 1, 1},
	{"a", "aa", 1, 1, 1},
	{"a", "aaa", 2, 2, 2},
	{"", "", 0, 0, 0},
	{"a", "b", 1, 1, 2},
	{"aaa", "aba", 1, 1, 2},
	{"aaa", "ab", 2, 2, 3},
	{"a", "a", 0, 0, 0},
	{"ab", "ab", 0, 0, 0},
	{"a", "", 1, 1, 1},
	{"aa", "a", 1, 1, 1},
	{"aaa", "a", 2, 2, 2},
	{"ab", "ba", 2, 1, 2},
	{"typo", "tyop", 2, 1, 2},
}

func checkDistance(t *testing.T, name, source, target string, correctDistance int, ops []EditOperation) {
	pair := NewRunes(source, target)
	matrix := NewMatrix(pair, ops)
	distance := matrix.Distance()

	if distance != correctDistance {
		t.Log(
			name,
			"distance between",
			source,
			"and",
			target,
			"computed as",
			distance,
			", should be",
			correctDistance)
		t.Fail()
	}
}

func TestDistanceForStrings(t *testing.T) {
	for _, testCase := range testCases {
		checkDistance(t, "Levenshtein", testCase.source, testCase.target, testCase.distance, DefaultLevenshtein)
		checkDistance(t, "Damerau", testCase.source, testCase.target, testCase.damerauDistance, DefaultDamerau)
		checkDistance(t, "LCS", testCase.source, testCase.target, testCase.lcsDistance, DefaultLCS)
	}
}

func ExampleNewMatrix() {
	source := "a"
	target := "aa"

	pair := NewRunes(source, target)
	distance := NewMatrix(pair, DefaultLevenshtein).Distance()
	fmt.Printf(`Distance between "%s" and "%s" computed as %d`, source, target, distance)
	// Output: Distance between "a" and "aa" computed as 1
}

func ExampleWrite() {
	pair := NewRunes("neighbor", "Neighbour")

	ops := []EditOperation{
		Match{},
		Insertion{1},
		Deletion{1},
		Substitution{2},
	}

	matrix := NewMatrix(pair, ops)
	Write(matrix, pair, os.Stdout)
	// Output:
	//       N  e  i  g  h  b  o  u  r
	//    0  1  2  3  4  5  6  7  8  9
	// n  1  2  3  4  5  6  7  8  9 10
	// e  2  3  2  3  4  5  6  7  8  9
	// i  3  4  3  2  3  4  5  6  7  8
	// g  4  5  4  3  2  3  4  5  6  7
	// h  5  6  5  4  3  2  3  4  5  6
	// b  6  7  6  5  4  3  2  3  4  5
	// o  7  8  7  6  5  4  3  2  3  4
	// r  8  9  8  7  6  5  4  3  4  3
}

var wordTestCases = []struct {
	source   []string
	target   []string
	distance int
}{
	{[]string{"Hello", "Word"}, []string{"Hello", "Word"}, 0},
	{[]string{"Hello", "Word"}, []string{"Hello"}, 1},
	{[]string{"Hello", "Word"}, []string{"Hello", "Beautiful", "Word"}, 1},
	{[]string{"Hello", "Word"}, []string{"My", "Word"}, 1},
	{[]string{"Hello", "Word"}, []string{"What", "a", "Beautiful", "Word"}, 3},
	{[]string{"What", "a", "dangerous", "Word"}, []string{"What", "a", "Beautiful", "Word"}, 1},
}

func checkWordDist(t *testing.T, name string, source, target []string, correctDistance int) {
	words := Words{Source: source, Target: target}
	matrix := NewMatrix(words, DefaultLevenshtein)
	distance := matrix.Distance()

	if distance != correctDistance {
		t.Log(
			name,
			"distance between",
			source,
			"and",
			target,
			"computed as",
			distance,
			", should be",
			correctDistance)
		t.Fail()
	}
	editScript := matrix.EditScript(words, DefaultLevenshtein)
	log.Printf("  %v, %v: script: %+v", source, target, editScript)
}

func TestWords(t *testing.T) {
	ExampleWrite()
	for _, testCase := range wordTestCases {
		checkWordDist(t, "word-Levenshtein", testCase.source, testCase.target, testCase.distance)
	}

}
