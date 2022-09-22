package scan

// SPDX-License-Identifier: Apache-2.0

import (
	"io"
)

// Point is an x,y coordinate
type Point struct {
	X int
	Y int
}

// Vector contains series of straight lines that sometimes join with rounded corners.
// A single Vector contains the following, which are drawn in the order listed:
// - zero or more Point
// - an optional rounded corner
// - an optional Vector
type Vector struct {
	Lines         []Point
	RoundedCorner bool
	*Vector
}

// ScanVector scans simple ASCII vector art:
//
// - = horizontal line
// | = vertical line
// / = diagonal line
// \ = diagonal line
// + = join - and |, which could be a corner or a T junction
// Slashes can be used to make corners instead of +, in which case they are rounded
//
// Horizontal and vertical lines are in the middle of a character, and corners take a 1/4 of a character
//
// You don't actually have to use + for corners and T junctions, the result is the same if a - and | meet.
// Leading and trailing blank lines are ignored.
// A trailing blank line or EOF is considered the end of the diagram.
// A diagram can contain blank lines and disjointed pieces.
//
// Each disjointed piece of the diagram is always drawn as a polygon using a series of lines and rounded corners,
// no portion is ever drawn as a rectangle.
//
// Example 1: simple one piece diagram of a box with one square corner and 3 rounded corners.
// It is 8 chars wide and 3 chars high, which means when scaled to x pixels horizontally, it will be x * 3/8 pixels high.
//
// +------\
// |      |
// \------/
//
// Example 2: two piece diagram of open boxes (note lack of + in corners)
//
// |----\
// |    |
// \    /
//
//	/    |
//	|  -/|
//	-----/
func ScanVector(src io.RuneScanner) {

}
