package entity

import (
	"fmt"

	"github.com/nx23/final-path/internal/config"
	"github.com/nx23/final-path/internal/gamemap"
	"github.com/nx23/final-path/internal/utils"
)

// Enemy is an enemy that follows the map path.
// X/Y coordinates always represent the enemy's center.
type Enemy struct {
	PositionX        float32 // Center X
	PositionY        float32 // Center Y
	Speed            float32
	CurrentPathIndex int
	Life             int
}

// NewEnemy creates an enemy positioned at the map's starting point
func NewEnemy(m gamemap.Map) *Enemy {
	if len(m) == 0 {
		return &Enemy{}
	}
	
	firstPath := m[0]
	return &Enemy{
		PositionX:        utils.CenterInPath(firstPath.StartX, config.PathWidth),
		PositionY:        utils.CenterInPath(firstPath.StartY, config.PathWidth),
		Speed:            2,
		CurrentPathIndex: 0,
		Life:             40,
	}
}

// IsAlive checks if the enemy is still alive
func (e *Enemy) IsAlive() bool {
	return e.Life > 0
}

// TakeDamage applies damage to the enemy
func (e *Enemy) TakeDamage(damage int) {
	e.Life -= damage
	if e.Life < 0 {
		e.Life = 0
	}
}

// FollowPath makes the enemy follow the map path.
// It automatically moves through the current path and advances to the next one when complete.
func (e *Enemy) FollowPath(m gamemap.Map) {
	if e.CurrentPathIndex >= len(m) {
		return
	}

	path := m[e.CurrentPathIndex]

	// Vertical movement
	if path.StartX == path.EndX {
		// Center the enemy on X axis for this vertical path
		targetCenterX := utils.CenterInPath(path.StartX, config.PathWidth)
		e.PositionX = targetCenterX

		targetCenterY := utils.CenterInPath(path.EndY, config.PathWidth)
		if e.PositionY < targetCenterY {
			e.PositionY += e.Speed
		} else {
			e.CurrentPathIndex++
			fmt.Printf("Path %d completed\n", e.CurrentPathIndex)
		}
	} else if path.StartY == path.EndY { // Horizontal movement
		// Center the enemy on Y axis for this horizontal path
		targetCenterY := utils.CenterInPath(path.StartY, config.PathWidth)
		e.PositionY = targetCenterY

		// Move right or left depending on direction
		if path.EndX > path.StartX {
			// Move right
			targetCenterX := utils.CenterInPath(path.EndX, config.PathWidth)
			if e.PositionX < targetCenterX {
				e.PositionX += e.Speed
			} else {
				e.CurrentPathIndex++
				fmt.Printf("Path %d completed\n", e.CurrentPathIndex)
			}
		} else {
			// Move left
			targetCenterX := utils.CenterInPath(path.EndX, config.PathWidth)
			if e.PositionX > targetCenterX {
				e.PositionX -= e.Speed
			} else {
				e.CurrentPathIndex++
				fmt.Printf("Path %d completed\n", e.CurrentPathIndex)
			}
		}
	}
}
