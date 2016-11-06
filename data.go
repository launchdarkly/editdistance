package editdistance

import (
	"fmt"
	"io"
)

type (
	StringSequencePair interface {
		SequencePair
		SourceAt(i int) string
		TargetAt(i int) string
	}

	Runes struct {
		Source []rune
		Target []rune
	}

	Words struct {
		Source []string
		Target []string
	}
)

var (
	_ StringSequencePair = Runes{}
	_ StringSequencePair = Words{}
)


func NewRunes(source, target string) Runes {
	return Runes{
		Source: []rune(source),
		Target: []rune(target),
	}
}

func (r Runes) SourceLen() int {
	return len(r.Source)
}
func (r Runes) TargetLen() int {
	return len(r.Target)
}
func (r Runes) Equal(i, j int) bool {
	return r.Source[i] == r.Target[j]
}
func (r Runes) SourceAt(i int) string {
	return string(r.Source[i])
}
func (r Runes) TargetAt(i int) string {
	return string(r.Target[i])
}


func (w Words) SourceLen() int {
	return len(w.Source)
}
func (w Words) TargetLen() int {
	return len(w.Target)
}
func (w Words) Equal(i, j int) bool {
	return w.Source[i] == w.Target[j]
}
func (w Words) SourceAt(i int) string {
	return w.Source[i]
}
func (w Words) TargetAt(i int) string {
	return w.Target[i]
}


func Write(matrix Matrix, p StringSequencePair, writer io.Writer) {
	fmt.Fprint(writer, "    ")
	for j := 0; j < p.TargetLen(); j++ {
		targetRune := p.TargetAt(j)
		fmt.Fprintf(writer, "  %s", targetRune)
	}
	fmt.Fprintln(writer)
	fmt.Fprintf(writer, "  %2d", matrix[0][0])
	for j := 0; j < p.TargetLen(); j++ {
		fmt.Fprintf(writer, " %2d", matrix[0][j+1])
	}
	fmt.Fprintln(writer)
	for i := 0; i < p.SourceLen(); i++ {
		sourceRune := p.SourceAt(i)
		fmt.Fprintf(writer, "%s %2d", sourceRune, matrix[i+1][0])
		for j := 0; j < p.TargetLen(); j++ {
			fmt.Fprintf(writer, " %2d", matrix[i+1][j+1])
		}
		fmt.Fprintln(writer)
	}
}
