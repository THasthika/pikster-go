package pikster

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"

	// Package image/{jpeg,gif,png} is not used explicitly in the code below
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
)

// ImageType is the type of image
type ImageType int

// ImageType constants
const (
	JPEG    ImageType = 0
	GIF     ImageType = 1
	PNG     ImageType = 2
	UNKNOWN ImageType = 10
)

type colorMapVal struct {
	colorPoint *ColorPoint
	index      int
}

// PImage holds all the details needed to cluster the image
type PImage struct {
	imgType      ImageType
	height       int
	width        int
	clusterCount uint
	colorMap     map[ColorPoint]colorMapVal
	colorList    []*ColorPoint
	clusterList  []*ColorPoint
	pixelList    []*PixelPoint
}

// NewPImage creates a new PImage structure given an image path
func NewPImage(filename string, clusterCount uint) *PImage {
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, name, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	imgType := getImageType(name)
	if imgType == UNKNOWN {
		log.Fatal("Unknown image Type")
	}

	bounds := img.Bounds()

	pimage := &PImage{
		imgType:      imgType,
		colorMap:     map[ColorPoint]colorMapVal{},
		colorList:    make([]*ColorPoint, 0),
		clusterList:  make([]*ColorPoint, 0),
		pixelList:    make([]*PixelPoint, 0),
		height:       bounds.Max.Y,
		width:        bounds.Max.X,
		clusterCount: clusterCount,
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r, g, b = r/0x101, g/0x101, b/0x101
			p := NewColorPoint(uint8(r), uint8(g), uint8(b))

			cmv, ok := pimage.colorMap[*p]
			if !ok {
				pimage.colorList = append(pimage.colorList, p)
				cmv = colorMapVal{
					colorPoint: p,
					index:      -1,
				}
				pimage.colorMap[*p] = cmv
			}

			pimage.pixelList = append(pimage.pixelList, NewPixelPoint(uint32(x), uint32(y), cmv.colorPoint))
		}
	}

	if len(pimage.colorList) < int(clusterCount) {
		log.Fatal("Error more clusters")
	}

	clusterPointIndexes := getNfromM(int(clusterCount), len(pimage.colorList))
	for _, i := range clusterPointIndexes {
		pimage.clusterList = append(pimage.clusterList, pimage.colorList[i].Copy())
	}

	return pimage
}

func getNfromM(n int, m int) []int {

	pickedColors := map[int]bool{}
	ret := make([]int, 0)

	for len(ret) < n {
		r := rand.Intn(m)
		_, v := pickedColors[r]
		if !v {
			pickedColors[r] = true
			ret = append(ret, r)
		}
	}

	return ret
}

// SaveFile saves the image to a file
func (p *PImage) SaveFile(name string) {
	filename := name + "." + getImageTypeExt(p.imgType)

	imgRect := image.Rect(0, 0, p.width, p.height)
	img := image.NewRGBA(imgRect)

	// color the image

	saveJPEG(filename, img)
}

func saveJPEG(filename string, img image.Image) {
	var opt jpeg.Options
	opt.Quality = 80

	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = jpeg.Encode(out, img, &opt)
	if err != nil {
		log.Fatal(err)
	}
}

func getImageTypeExt(t ImageType) string {
	switch t {
	case JPEG:
		return "jpeg"
	case PNG:
		return "png"
	case GIF:
		return "gif"
	}
	return "unknown"
}

func getImageType(t string) ImageType {
	switch t {
	case "jpeg":
		fallthrough
	case "jpg":
		return JPEG
	case "png":
		return PNG
	case "gif":
		return GIF
	}

	return UNKNOWN
}

// Run runs
func Run() {

	pimage := NewPImage("./img.jpg", 3)
	fmt.Println(len(pimage.colorMap))
	fmt.Println(len(pimage.pixelList))
	fmt.Println(len(pimage.colorList))
	fmt.Println(len(pimage.clusterList))

	// pimage.SaveFile("xxx")
}
