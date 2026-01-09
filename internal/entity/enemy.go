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

type NewEnemyParams struct {
	Map   gamemap.Map
	Speed float32
	Life  int
}

func NewEnemy(params NewEnemyParams) *Enemy {
	if len(params.Map) == 0 {
		return &Enemy{}
	}

	firstPath := params.Map[0]
	return &Enemy{
		PositionX:        utils.CenterInPath(firstPath.StartX, config.PathWidth),
		PositionY:        utils.CenterInPath(firstPath.StartY, config.PathWidth),
		Speed:            params.Speed,
		CurrentPathIndex: 0,
		Life:             params.Life,
	}
}

func (e *Enemy) IsAlive() bool {
	return e.Life > 0
}

func (e *Enemy) TakeDamage(damage int) {
	e.Life -= damage
	if e.Life < 0 {
		e.Life = 0
	}
}

// FollowPath moves enemy along the map path
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
		// Horizontal movement
	} else if path.StartY == path.EndY {
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
