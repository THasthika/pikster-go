package pikster

import (
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"os"

	// Package image/{jpeg,gif,png} is not used explicitly in the code below
	_ "image/gif"
	_ "image/png"

	"image/jpeg"
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
	colorMap     map[ColorPoint]*colorMapVal
	colorList    []*ColorPoint
	clusterList  []*ColorPoint
	pixelList    []*PixelPoint
}

func NewPImageFromFile(filename string, clusterCount uint) *PImage {
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

	return NewPImageFromImage(img, imgType, clusterCount)
}

func NewPImageFromImage(img image.Image, imageType ImageType, clusterCount uint) *PImage {
	bounds := img.Bounds()

	pimage := &PImage{
		imgType:      imageType,
		colorMap:     map[ColorPoint]*colorMapVal{},
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
				cmv = &colorMapVal{
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

// SaveFile saves the image to a file
func (p *PImage) SaveFile(name string) {
	filename := name + "." + getImageTypeExt(p.imgType)

	imgRect := image.Rect(0, 0, p.width, p.height)
	img := image.NewRGBA(imgRect)

	// color the image
	for _, pp := range p.pixelList {
		cmv := p.colorMap[*pp.color]
		cp := p.clusterList[cmv.index]
		img.SetRGBA(int(pp.x), int(pp.y), cp.ToColor())
	}

	saveJPEG(filename, img)
}

func (p *PImage) RunClusteringStep() int {
	for _, cp := range p.colorList {
		closest := cp.Distance(p.clusterList[0])
		closestCluster := 0

		for i := 1; i < len(p.clusterList); i++ {
			t := cp.Distance(p.clusterList[i])
			if t < closest {
				closest = t
				closestCluster = i
			}
		}

		cmv := p.colorMap[*cp]
		cmv.index = closestCluster
	}

	var totalChange = 0

	for i, cp := range p.clusterList {
		var sr float64
		var sg float64
		var sb float64
		var n uint64

		prevR := cp.r
		prevG := cp.g
		prevB := cp.b

		for k, v := range p.colorMap {
			if v.index != i {
				continue
			}
			sr += float64(k.r)
			sg += float64(k.g)
			sb += float64(k.b)
			n++
		}

		sr, sg, sb = sr/float64(n), sg/float64(n), sb/float64(n)
		cp.r = uint8(sr)
		cp.g = uint8(sg)
		cp.b = uint8(sb)

		totalChange += int(math.Abs(float64(prevR-cp.r)) + math.Abs(float64(prevG-cp.g)) + math.Abs(float64(prevB-cp.b)))
	}

	fmt.Printf("change: %d\n", totalChange)

	fmt.Println("Clustering Step Done!")

	return totalChange
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

// // Run runs
// func Run() {

// 	pimage := NewPImage("./img2.png", 10)
// 	// fmt.Println(len(pimage.colorMap))
// 	// fmt.Println(len(pimage.pixelList))
// 	// fmt.Println(len(pimage.colorList))
// 	// fmt.Println(len(pimage.clusterList))

// 	for i := 0; i < 100; i++ {
// 		if pimage.runClusteringStep() == 0 {
// 			break
// 		}
// 	}

// 	pimage.SaveFile("xxx")
// }
