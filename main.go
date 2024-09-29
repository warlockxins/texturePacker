// starter code to create tiled image
// https://yourbasic.org/golang/create-image/
package main

import (
	// "fmt"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"warlockxins/texturepack/utils"
)

func readPngNamesAt(dirName string) []string {
	files, err := os.ReadDir(dirName)
	if err != nil {
		panic(err)
	}

	var imageFileNames []string
	for _, file := range files {
		// fmt.Println(file.Name(), file.IsDir())
		if file.IsDir() {
			continue
		}
		extension := filepath.Ext(file.Name())
		if extension == ".png" {
			imageFileNames = append(imageFileNames, dirName+"/"+file.Name())
		}
	}

	return imageFileNames
}

func readAnimFolderNamesAt(dirName *string) []string {
	files, err := os.ReadDir(*dirName)
	if err != nil {
		panic(err)
	}

	var animFloderNames []string
	for _, file := range files {
		// fmt.Println(file.Name(), file.IsDir())
		if !file.IsDir() {
			continue
		}

		animFloderNames = append(animFloderNames, file.Name())
	}

	return animFloderNames
}

func doSample() {
	config, err := utils.GetAnimationConfig("./testSamples/")
	if err != nil {
		panic(err)
	}

	imagesWithBounds := config.ToImagesWithBounds("./testSamples/")

	imagesWithBounds.ToSpritesheet("./testSamples/spriteSheet.png")

}

func main() {
	doSample()
	inFolderPtr := flag.String("in", "", "an input folder path string")
	outFolderPtr := flag.String("out", "", "an output folder path string for texture atlas")
	flag.Parse()

	if *inFolderPtr == "" {
		println("input folder not provided")
		os.Exit(1)
	}

	if *outFolderPtr == "" {
		println("output folder not provided")
		os.Exit(1)
	}

	fmt.Println("input:", *inFolderPtr)
	fmt.Println("output:", *outFolderPtr)
}

func process(inFolderPtr *string, outFolderPtr *string) {
	animationFolders := readAnimFolderNamesAt(inFolderPtr)

	fmt.Println("input folders:", animationFolders)
	// animationFolders := [...]string{"armActionTake", "idle", "run", "walk", "walkCrouch", "death"}

	directionFolders := [...]string{"N", "NE", "E", "SE", "S"}

	// animationFolders = []string{}

	for _, anim := range animationFolders {
		for _, direction := range directionFolders {
			images := readPngNamesAt(*inFolderPtr + "/" + anim + "/" + direction)
			// fmt.Println("iterate over", images)
			imageSize := 128
			width := len(images) * imageSize
			height := imageSize

			upLeft := image.Point{0, 0}
			lowRight := image.Point{width, height}

			img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

			// -------------------------------------------------------------------------------------------
			// now combine them together

			for imageIndex := 0; imageIndex < len(images); imageIndex++ {
				// -------------------------------------------------------------------------------------------
				filename := images[imageIndex]
				src := *utils.GetFile(filename)
				// src.Bounds()

				// Create a new grayscale image
				bounds := src.Bounds()
				w, h := bounds.Max.X, bounds.Max.Y
				if w != imageSize || h != imageSize {
					panic("image not the same size: " + filename)
				}

				xOfset := imageIndex * imageSize

				for x := 0; x < imageSize; x++ {
					for y := 0; y < imageSize; y++ {
						img.Set(x+xOfset, y, src.At(x, y))
					}
				}

			}
			// encode as png
			f, err := os.Create(*outFolderPtr + "/" + anim + "-" + direction + ".png")
			if err != nil {
				panic(err)
			}
			png.Encode(f, img)
		}
	}
}
