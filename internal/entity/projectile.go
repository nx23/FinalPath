package entity

import "math"

// Projectile is a projectile fired by a tower.
// It automatically chases the enemy until it hits.
type Projectile struct {
	PositionX float32 // Center X
	PositionY float32 // Center Y
	Speed     int
	Target    *Enemy
}

// NewProjectile creates a projectile that will aim at the enemy
func NewProjectile(x, y float32, target *Enemy) Projectile {
	return Projectile{
		PositionX: x,
		PositionY: y,
		Speed:     10,
		Target:    target,
	}
}

// Hit moves the projectile towards the target.
// Returns true when it hits the enemy (distance < speed).
func (p *Projectile) Hit() bool {
	if p.Target == nil || !p.Target.IsAlive() {
		return false
	}

	// Calculate direction and move in a straight line
	dx := p.Target.PositionX - p.PositionX
	dy := p.Target.PositionY - p.PositionY
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance < float32(p.Speed) {
		// Reached the target
		return true
	}
	
	// Move towards the target
	p.PositionX += (dx / distance) * float32(p.Speed)
	p.PositionY += (dy / distance) * float32(p.Speed)

	return false
}
