package space

import (
	"bytes"
	"fmt"
	"os"
)

type Bounds struct {
	X, Y, Width, Height int
}

type Box struct {
	Width, Height int
}

type Space struct {
	Ocupied bool
	Bounds  Bounds
	Spaces  []Space
}

func NewSpace(b Bounds, ocupied bool) *Space {
	return &Space{
		Ocupied: ocupied,
		Bounds:  b,
	}
}

var gap = 2

func (space *Space) InsertSpace(box *Box, newBounds *Bounds) bool {
	if len(space.Spaces) == 0 {
		if space.Bounds.X+box.Width <= space.Bounds.X+space.Bounds.Width &&
			space.Bounds.Y+box.Height <= space.Bounds.Y+space.Bounds.Height {

			newBounds.X = space.Bounds.X
			newBounds.Y = space.Bounds.Y
			newBounds.Width = box.Width
			newBounds.Height = box.Height

			// left ocupied space that actually used by box
			space.Spaces = append(
				space.Spaces,
				// new ocupied section
				*NewSpace(
					Bounds{
						space.Bounds.X, space.Bounds.Y, box.Width, box.Height,
					},
					true,
				),
				// new internal available space to the right
				*NewSpace(
					Bounds{
						space.Bounds.X + box.Width + gap, space.Bounds.Y,
						space.Bounds.Width - box.Width - gap, box.Height,
					},
					false,
				),
				// row below
				*NewSpace(
					Bounds{
						space.Bounds.X, space.Bounds.Y + box.Height + gap,
						space.Bounds.Width, space.Bounds.Height - box.Height - gap,
					},
					false,
				),
			)
			// fmt.Println("inserted?", space.spaces)

			return true
		}
		return false
	}

	// iterate existing spaces
	for i := 0; i < len(space.Spaces); i++ {
		subSpace := &space.Spaces[i]
		if subSpace.Ocupied == true {
			continue
		}

		inserted := subSpace.InsertSpace(box, newBounds)
		if inserted {
			return true
		}
	}

	return false
}

// --------- Visualize for debug purposes

var colors = []string{"#8ecae6", "#219ebc", "#023047", "#ffb703", "#fb8500", "#606c38", "#264653", "#cdb4db"}
var currentColor = 0

func nextColor() string {
	c := colors[currentColor]
	currentColor++
	if currentColor >= len(colors) {
		currentColor = 0
	}
	return c
}

func (s *Space) SaveToSVG() {
	var svgAccumulator bytes.Buffer
	// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/rect
	svgAccumulator.WriteString(
		fmt.Sprintf(
			"<svg viewBox=\"0 0 %d %d\" xmlns=\"http://www.w3.org/2000/svg\">", s.Bounds.Width, s.Bounds.Height,
		),
	)

	s.buildSVG(&svgAccumulator)

	svgAccumulator.WriteString("</svg>")

	os.WriteFile("./svgRects.svg", svgAccumulator.Bytes(), 0644)

}

func (s *Space) buildSVG(accumulator *bytes.Buffer) {
	// https://coolors.co/palettes/trending
	if s.Ocupied {
		accumulator.WriteString(
			fmt.Sprintf(
				"<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" rx=\"5\" fill=\"%s\"/>",
				s.Bounds.X, s.Bounds.Y, s.Bounds.Width, s.Bounds.Height, nextColor(),
			),
		)
	} else {
		for i := 0; i < len(s.Spaces); i++ {
			subSpace := &s.Spaces[i]
			subSpace.buildSVG(accumulator)
		}
	}
}
