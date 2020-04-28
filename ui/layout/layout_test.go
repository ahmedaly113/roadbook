package layout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRows(t *testing.T) {
	l := New(10, 10)
	r1 := l.Row(2)
	r2 := l.Row(6)
	r3 := l.Row(2)

	assert.Equal(t, 10, r1.W())
	assert.Equal(t, 10, r2.W())
	assert.Equal(t, 10, r3.W())

	assert.Equal(t, 2, r1.H())
	assert.Equal(t, 6, r2.H())
	assert.Equal(t, 2, r3.H())

	assert.Equal(t, 0, l.remainingHeight)
	assert.Equal(t, 10, l.remainingWidth)
}
func TestCols(t *testing.T) {
	l := New(10, 10)
	c1 := l.Col(2)
	c2 := l.Col(6)
	c3 := l.Col(2)

	assert.Equal(t, 10, c1.H())
	assert.Equal(t, 10, c2.H())
	assert.Equal(t, 10, c3.H())

	assert.Equal(t, 2, c1.W())
	assert.Equal(t, 6, c2.W())
	assert.Equal(t, 2, c3.W())

	assert.Equal(t, 10, l.remainingHeight)
	assert.Equal(t, 0, l.remainingWidth)
}

func TestRowsAndCols(t *testing.T) {
	l := New(10, 10)

	// consume top row
	l.Row(5)
	bottomRow := l.Row(5)
	assert.Equal(t, 0, bottomRow.x)
	assert.Equal(t, 0, bottomRow.absX())
	assert.Equal(t, 5, bottomRow.absY())

	// consume 1 column of 1 width, expect 9 maining
	bottomRow.Col(1)
	assert.Equal(t, 9, bottomRow.remainingWidth)

	// consume 1 column of width 8
	bottomRow.Col(8)

	// final col should be 1x5 at x=9,y=5
	testCol := bottomRow.Col(1)
	assert.Equal(t, 1, testCol.W())
	assert.Equal(t, 5, testCol.H())
	assert.Equal(t, 9, testCol.absX())
	assert.Equal(t, 5, testCol.absY())
}

func TestDeepNesting(t *testing.T) {
	l := New(100, 100)
	current := l
	depth := 0

	// alternate row and column, discarding 10 pixels; 18 is 9 discards in
	// each dimension, leaving a single 10x10 square
	for depth < 18 {
		if depth%2 == 0 {
			// discard 10 pixel row
			current.Row(10)
			// remainder becomes current
			current = current.Row(current.remainingHeight)
		} else {
			// discard 10 pixel col
			current.Col(10)
			// remainder becomes current
			current = current.Row(current.remainingWidth)
		}
		depth++
	}

	assert.Equal(t, 10, current.W())
	assert.Equal(t, 10, current.H())
	assert.Equal(t, 90, current.absX())
	assert.Equal(t, 90, current.absY())

	parentCount := 1
	parent := current.parent
	for parent.parent != nil {
		parent = parent.parent
		parentCount++
	}

	assert.Equal(t, 18, parentCount)
}
