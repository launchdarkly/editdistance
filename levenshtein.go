package editdistance

import (
	"fmt"
	"math"
)

type (
	TemporalIndices struct {
		Source int
		Target int
	}

	EditScript []EditOperation

	EditOperations []EditOperation

	EditOperation interface {
		fmt.Stringer
		Apply(seqPair SequencePair, matrix Matrix, i, j int) (int, bool)
		Backtrack(matrix Matrix, i, j int) (int, int)
	}

	Matrix [][]int

	SequencePair interface {
		SourceLen() int
		TargetLen() int
		Equal(i, j int) bool
	}
)

// Apply computes the lowest cost of applying one of the operations
func (ops EditOperations) Apply(seqPair SequencePair, matrix Matrix, row, col int) (int, bool) {
	found := false
	lowestCost := math.MaxInt32
	for _, op := range ops {
		if cost, ok := op.Apply(seqPair, matrix, row, col); ok && cost < lowestCost {
			lowestCost = cost
			found = true
		}
	}
	return lowestCost, found
}

// Find finds the corresponding operation
func (ops EditOperations) Find(seqPair SequencePair, matrix Matrix, row, col int) EditOperation {
	for _, op := range ops {
		if cost, ok := op.Apply(seqPair, matrix, row, col); ok && cost == matrix[row][col] {
			return op
		}
	}
	return nil
}

// NewMatrix creates the new Levenshtein matrix
func NewMatrix(seqPair SequencePair, ops EditOperations) Matrix {
	// Make a 2-D matrix. Rows correspond to prefixes of source, columns to
	// prefixes of target. Cells will contain edit distances.
	// Cf. http://www.let.rug.nl/~kleiweg/lev/levenshtein.html
	height := seqPair.SourceLen() + 1
	width := seqPair.TargetLen() + 1
	matrix := make(Matrix, height)

	for i := 0; i < height; i++ {
		matrix[i] = make([]int, width)
	}

	// Fill in the remaining cells: for each prefix pair, choose the
	// (edit history, operation) pair with the lowest cost.
	matrix[0][0] = 0
	for j := 1; j < width; j++ {
		matrix[0][j], _ = ops.Apply(seqPair, matrix, 0, j)
	}
	for i := 1; i < height; i++ {
		for j := 0; j < width; j++ {
			matrix[i][j], _ = ops.Apply(seqPair, matrix, i, j)
		}
	}
	return matrix
}

// Distance reads the edit distance off the Levenshtein matrix.
func (m Matrix) Distance() int {
	return m[len(m)-1][len(m[0])-1]
}

// EditScript returns an optimal edit script based on the given Levenshtein matrix.
func (m Matrix) EditScript(seqPair SequencePair, ops EditOperations) (EditScript, []TemporalIndices) {
	i := len(m) - 1
	j := len(m[0]) - 1
	script := EditScript{}
	indices := []TemporalIndices{}
	op := ops.Find(seqPair, m, i, j)
	for op != nil {
		indices = append(indices, TemporalIndices{i, j})
		script = append(script, op)
		i, j = op.Backtrack(m, i, j)
		op = ops.Find(seqPair, m, i, j)
	}

	for k := 0; k < len(script)/2; k++ {
		script[k], script[len(script)-k-1] = script[len(script)-k-1], script[k]
		indices[k], indices[len(indices)-k-1] = indices[len(indices)-k-1], indices[k]
	}

	return script, indices
}

type Operation struct {
	op          EditOperation
	sourceIndex int
	targetIndex int
}
