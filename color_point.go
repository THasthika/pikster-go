package pikster

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
