package gamemap

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nx23/final-path/internal/config"
	"github.com/nx23/final-path/internal/utils"
)

// Path is a segment of the map path (can be vertical or horizontal)
type Path struct {
	StartX float32
	StartY float32
	EndX   float32
	EndY   float32
}

// Map is the complete path formed by multiple segments
type Map []Path

// DefaultMap returns the default game map
func DefaultMap() Map {
	return Map{
		{StartX: 350, StartY: 0, EndX: 350, EndY: 150},
		{StartX: 350, StartY: 150, EndX: 550, EndY: 150},
		{StartX: 550, StartY: 150, EndX: 550, EndY: 350},
		{StartX: 550, StartY: 350, EndX: 150, EndY: 350},
		{StartX: 150, StartY: 350, EndX: 150, EndY: 600},
	}
}

// Draw renders the map on the screen
func (m Map) Draw(screen *ebiten.Image) {
	for _, path := range m {
		width := path.EndX - path.StartX
		height := path.EndY - path.StartY

		// Vertical path has width=0, so we use PathWidth
		if width == 0 {
			width = config.PathWidth
		}

		// Horizontal path has height=0, so we use PathWidth
		if height == 0 {
			height = config.PathWidth
		}

		vector.FillRect(screen, path.StartX, path.StartY, width, height, color.White, false)
	}
}

// IsPositionOnPath checks if a position is on the path.
// Uses a 30px margin to make tower validation easier.
func IsPositionOnPath(x, y float32, m Map) bool {
	const margin float32 = 30

	for _, path := range m {
		minX := utils.Min(path.StartX, path.EndX) - margin
		maxX := utils.Max(path.StartX, path.EndX) + config.PathWidth + margin
		minY := utils.Min(path.StartY, path.EndY) - margin
		maxY := utils.Max(path.StartY, path.EndY) + config.PathWidth + margin

		if x >= minX && x <= maxX && y >= minY && y <= maxY {
			return true
		}
	}

	return false
}
