package utils

// CenteredPosition helps work with entities that use centered coordinates.
// Makes it easy to convert between center and top-left for screen drawing.
type CenteredPosition struct {
	X    float32
	Y    float32
	Size float32
}

// TopLeft converts center position to top-left coordinates.
// Useful for drawing rectangles that need the upper-left corner position.
func (cp CenteredPosition) TopLeft() (float32, float32) {
	halfSize := cp.Size / 2
	return cp.X - halfSize, cp.Y - halfSize
}

// Center returns the center coordinates (already centered, just for consistency)
func (cp CenteredPosition) Center() (float32, float32) {
	return cp.X, cp.Y
}

// getCenterFromTopLeft does the reverse: top-left -> center
func GetCenterFromTopLeft(topLeftX, topLeftY, size float32) (float32, float32) {
	halfSize := size / 2
	return topLeftX + halfSize, topLeftY + halfSize
}

// centerInPath calculates where the center of a path is.
// Example: path at X=350 with width 50 -> center at X=375
func CenterInPath(pathPos float32, pathWidth float32) float32 {
	return pathPos + pathWidth/2
}

// min returns the minimum of two float32 values
func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two float32 values
func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
