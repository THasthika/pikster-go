package pikster_test

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"testing"

	"github.com/THasthika/pikster-go"
)

func createTestImage(width int, height int, colorCount uint) image.Image {
	topLeft := image.Point{0, 0}
	bottomRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})

	colors := make([]color.Color, 0)

	for i := 0; i < int(colorCount); i++ {
		color := color.RGBA{uint8(rand.Uint32() % 255), uint8(rand.Uint32() % 255), uint8(rand.Uint32() % 255), 0xff}
		colors = append(colors, color)
	}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, colors[rand.Uint32()%uint32(colorCount)])
		}
	}

	return img
}

func Test_In_Memory_Image_Clustering(t *testing.T) {

	// image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	// image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	// image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	img := createTestImage(3000, 3000, 10)

	pImg := pikster.NewPImageFromImage(img, pikster.PNG, 5)

	for i := 0; i < 100; i++ {
		if pImg.RunClusteringStep() == 0 {
			break
		}
	}

	// pikster.Run()
}

func Test_Image_File_Clustering(t *testing.T) {

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	img := createTestImage(3000, 3000, 10)
	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)

	pImg := pikster.NewPImageFromFile("image.png", 3)

	for i := 0; i < 100; i++ {
		if pImg.RunClusteringStep() == 0 {
			break
		}
	}

	pImg.SaveFile("image.out")

	os.Remove("image.png")
	os.Remove("image.out.png")

}
