package pikster

import (
	"image/color"
	"math"
)

// ColorPoint is a structure that keeps a color
type ColorPoint struct {
	r uint8
	g uint8
	b uint8
}

// NewColorPoint creates a new color point
func NewColorPoint(r uint8, g uint8, b uint8) *ColorPoint {
	return &ColorPoint{
		r: r,
		g: g,
		b: b,
	}
}

// Copy colorpoint
func (cp *ColorPoint) Copy() *ColorPoint {
	return &ColorPoint{
		r: cp.r,
		g: cp.g,
		b: cp.b,
	}
}

// Distance compute distance
func (cp *ColorPoint) Distance(other *ColorPoint) float32 {
	dr := math.Abs(float64(cp.r) - float64(other.r))
	dg := math.Abs(float64(cp.g) - float64(other.g))
	db := math.Abs(float64(cp.b) - float64(other.b))
	return float32(dr + dg + db)
}

// ToColor xxx
func (cp *ColorPoint) ToColor() color.RGBA {
	return color.RGBA{
		R: cp.r,
		G: cp.g,
		B: cp.b,
		A: 0xFF,
	}
}
