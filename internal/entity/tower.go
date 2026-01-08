package entity

import (
	"github.com/nx23/final-path/internal/config"
	"github.com/nx23/final-path/internal/gamemap"
)

// Tower is a defense tower that attacks enemies within range.
// X/Y coordinates always represent the tower's center.
type Tower struct {
	PositionX    float32 // Center X
	PositionY    float32 // Center Y
	Range        float32
	Damage       int
	FireRate     float32 // Shots per second
	LastFireTime int
}

// NewTower creates a tower at the specified position
func NewTower(x, y float32) Tower {
	return Tower{
		PositionX: x,
		PositionY: y,
		Range:     100,
		Damage:    10,
		FireRate:  1,
	}
}

// IsEnemyInRange checks if the enemy is within range.
// Uses squared distance to avoid sqrt (faster).
func (t *Tower) IsEnemyInRange(enemy *Enemy) bool {
	dx := t.PositionX - enemy.PositionX
	dy := t.PositionY - enemy.PositionY
	distanceSquared := dx*dx + dy*dy
	return distanceSquared <= t.Range*t.Range
}

// CanFire checks if enough time has passed since the last shot.
// Respects the tower's FireRate (shots per second).
func (t *Tower) CanFire(currentTick int) bool {
	ticksPerShot := int(60.0 / t.FireRate)
	return currentTick-t.LastFireTime >= ticksPerShot
}

// Attack creates a projectile that will chase the enemy
func (t *Tower) Attack(enemy *Enemy) Projectile {
	return NewProjectile(t.PositionX, t.PositionY, enemy)
}

// CanPlaceTower validates if a tower can be placed at the position.
// Checks all corners and center of the tower to ensure it's not on the path.
func CanPlaceTower(centerX, centerY float32, m gamemap.Map) bool {
	const halfSize = config.TowerSize / 2

	// Verify all four corners and center of the tower area
	points := []struct{ x, y float32 }{
		{centerX, centerY},                       // Center
		{centerX - halfSize, centerY - halfSize}, // Top left corner
		{centerX + halfSize, centerY - halfSize}, // Top right corner
		{centerX - halfSize, centerY + halfSize}, // Bottom left corner
		{centerX + halfSize, centerY + halfSize}, // Bottom right corner
	}

	// All points must be outside the path
	for _, point := range points {
		if gamemap.IsPositionOnPath(point.x, point.y, m) {
			return false
		}
	}

	return true
}
