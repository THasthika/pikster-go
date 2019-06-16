package pikster_test

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/tharindu96/pikster-go"
)

func Test_Main(t *testing.T) {

	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	pikster.Run()
}
