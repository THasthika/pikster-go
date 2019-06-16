package pikster

// PixelPoint is a structure that keeps a pixel
type PixelPoint struct {
	x     uint32
	y     uint32
	color *ColorPoint
}

// NewPixelPoint creates a new pixel point
func NewPixelPoint(x uint32, y uint32, color *ColorPoint) *PixelPoint {
	return &PixelPoint{
		x:     x,
		y:     y,
		color: color,
	}
}
