package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"warlockxins/texturepack/space"
)

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

type BoxContainer struct {
	box       space.Box
	newBounds space.Bounds
}

type ByHeight []BoxContainer

func (a ByHeight) Len() int           { return len(a) }
func (a ByHeight) Less(i, j int) bool { return a[i].box.Height > a[j].box.Height }
func (a ByHeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func main() {
	s := space.NewSpace(space.Bounds{X: 0, Y: 0, Width: 512, Height: 512}, false)

	boxes := []BoxContainer{

		{space.Box{Width: 90, Height: 21}, space.Bounds{}},
		{space.Box{Width: 80, Height: 215}, space.Bounds{}},
		{space.Box{Width: 80, Height: 115}, space.Bounds{}},

		{space.Box{Width: 200, Height: 215}, space.Bounds{}},
		{space.Box{Width: 100, Height: 21}, space.Bounds{}},

		{space.Box{Width: 50, Height: 90}, space.Bounds{}},

		{space.Box{Width: 90, Height: 24}, space.Bounds{}},
		{space.Box{Width: 200, Height: 100}, space.Bounds{}},
		{space.Box{Width: 80, Height: 100}, space.Bounds{}},
		{space.Box{Width: 80, Height: 50}, space.Bounds{}},
		{space.Box{Width: 100, Height: 24}, space.Bounds{}},

		{space.Box{Width: 20, Height: 84}, space.Bounds{}},
		{space.Box{Width: 80, Height: 100}, space.Bounds{}},
		{space.Box{Width: 80, Height: 20}, space.Bounds{}},
		{space.Box{Width: 100, Height: 104}, space.Bounds{}},
	}

	// important to sort
	if 1 == 1 {
		sort.Sort(ByHeight(boxes))
	}

	for i := 0; i < len(boxes); i++ {
		s.InsertSpace(&boxes[i].box, &boxes[i].newBounds)
	}

	var svgAccumulator bytes.Buffer
	// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/rect
	svgAccumulator.WriteString("<svg viewBox=\"0 0 715 512\" xmlns=\"http://www.w3.org/2000/svg\">")

	// box for presenting existing boxes
	svgAccumulator.WriteString("<rect x=\"512\" y=\"0\" width=\"210\" height=\"512\" stroke=\"red\"/>")

	var y = 0
	currentColor = 0
	scale := 5
	for _, box := range boxes {
		svgAccumulator.WriteString(
			fmt.Sprintf(
				"<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" rx=\"5\" fill=\"%s\"/>",
				513, y, box.box.Width/scale, box.box.Height/scale, nextColor(),
			),
		)

		y = y + box.box.Height/scale
	}

	svgAccumulator.WriteString("<rect x=\"0\" y=\"0\" width=\"512\" height=\"512\" fill=\"#ccd5ae\"/>")

	currentColor = 0
	for _, box := range boxes {
		svgAccumulator.WriteString(
			fmt.Sprintf(
				"<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" rx=\"5\" fill=\"%s\"/>",
				box.newBounds.X, box.newBounds.Y, box.newBounds.Width, box.newBounds.Height, nextColor(),
			),
		)
	}

	// buildSVG(s, &svgAccumulator)
	svgAccumulator.WriteString("</svg>")
	// fmt.Println("----------")
	// fmt.Println(svgAccumulator.String())

	os.WriteFile("./svgRects.svg", svgAccumulator.Bytes(), 0644)

	fmt.Println("===========")
}
