package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"sort"
	"warlockxins/texturepack/space"
)

type AnimationConfigDirectionList struct {
	N  []string
	NE []string
	E  []string
	SE []string
	S  []string
}

type AnimationConfig map[string]AnimationConfigDirectionList

func GetAnimationConfig(folderPathWithConfig string) (*AnimationConfig, error) {
	content, err := os.ReadFile(filepath.Join(folderPathWithConfig, "out.json"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading file: %v ", folderPathWithConfig))
	}
	animationConfig := &AnimationConfig{}
	if err := json.Unmarshal(content, animationConfig); err != nil {
		return nil, errors.New(fmt.Sprintf("Not valid json at: %v ", folderPathWithConfig))
	}

	return animationConfig, nil
}

type ImageWithBounds struct {
	ImageName           string
	Image               image.Image
	NonAlphaBounds      image.Rectangle
	NonAlphaSize        space.Box
	TargetTextureBounds *space.Bounds
}

type ImagesWithBounds []ImageWithBounds

type SortByHeight ImagesWithBounds

func (a SortByHeight) Len() int { return len(a) }
func (a SortByHeight) Less(i, j int) bool {
	return a[i].NonAlphaSize.Height > a[j].NonAlphaSize.Height
}
func (a SortByHeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func newBounds() *space.Bounds {
	return &space.Bounds{X: 0, Y: 0, Width: 0, Height: 0}
}
func (ac *AnimationConfig) ToImagesWithBounds(folderPathWithConfig string) *ImagesWithBounds {
	var images []string = []string{}
	for _, value := range *ac {
		images = append(images, value.N...)
		images = append(images, value.NE...)
		images = append(images, value.E...)
		images = append(images, value.SE...)
		images = append(images, value.S...)
	}

	imageSpaces := ImagesWithBounds{}

	for _, imageName := range images {

		imageContent := GetFile(filepath.Join(folderPathWithConfig, "images", imageName))
		nonAlphaBounds := GetImageNonAlphaBounds(imageContent)
		targetBounds := newBounds()

		imageSpaces = append(imageSpaces, ImageWithBounds{
			ImageName:      imageName,
			Image:          *imageContent,
			NonAlphaBounds: nonAlphaBounds,
			NonAlphaSize: space.Box{
				Width:  nonAlphaBounds.Dx(),
				Height: nonAlphaBounds.Dy(),
			},
			TargetTextureBounds: targetBounds,
		})

		// fmt.Println("---", imageName, nonAlphaBounds.Dx(), nonAlphaBounds.Dy())

	}

	sort.Sort(SortByHeight(imageSpaces))

	sheetSpaces := space.NewSpace(space.Bounds{X: 0, Y: 0, Width: 1024, Height: 1024}, false)

	for i := 0; i < len(imageSpaces); i++ {
		imageForSpace := imageSpaces[i]
		inserted := sheetSpaces.InsertSpace(&imageForSpace.NonAlphaSize, imageForSpace.TargetTextureBounds)

		if inserted == false {
			fmt.Println(imageForSpace.NonAlphaBounds)
			panic("image not inserted: " + imageForSpace.ImageName)
		}
		// fmt.Println("===", imageForSpace.TargetTextureBounds)

	}

	sheetSpaces.SaveToSVG()

	return &imageSpaces
}

func GetImageNonAlphaBounds(src *image.Image) image.Rectangle {
	b := (*src).Bounds()

	min := image.Point{
		X: 10000,
		Y: 10000,
	}
	max := image.Point{
		X: 0,
		Y: 0,
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			color := (*src).At(x, y)
			_, _, _, a := color.RGBA()
			if a > 0 {
				min.X = int(math.Min(float64(min.X), float64(x)))
				min.Y = int(math.Min(float64(min.Y), float64(y)))

				max.X = int(math.Max(float64(max.X), float64(x)))
				max.Y = int(math.Max(float64(max.Y), float64(y)))
			}

		}
	}

	return image.Rectangle{
		Min: min,
		Max: max,
	}
}

func GetFile(filename string) *image.Image {
	infile, err := os.Open(filename)
	defer infile.Close()
	if err != nil {
		// replace this with real error handling
		panic(err.Error())
	}

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(infile)
	if err != nil {
		// replace this with real error handling
		panic(err.Error())
	}

	return &src
}
