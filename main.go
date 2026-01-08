package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Constants for entity sizes
const (
	EnemySize      float32 = 25
	TowerSize      float32 = 25
	PathWidth      float32 = 50
	ProjectileSize float32 = 5
)

// Centered position helper - returns the top-left corner coordinates
// for an entity of given size that should be centered at (x, y)
type CenteredPosition struct {
	X    float32
	Y    float32
	Size float32
}

// TopLeft returns the top-left corner coordinates for drawing
func (cp CenteredPosition) TopLeft() (float32, float32) {
	halfSize := cp.Size / 2
	return cp.X - halfSize, cp.Y - halfSize
}

// Center returns the center coordinates
func (cp CenteredPosition) Center() (float32, float32) {
	return cp.X, cp.Y
}

// Helper to get center from top-left position and size
func getCenterFromTopLeft(topLeftX, topLeftY, size float32) (float32, float32) {
	halfSize := size / 2
	return topLeftX + halfSize, topLeftY + halfSize
}

// Helper to center an entity in a path
// pathPos is the start position of the path (e.g., StartX or StartY)
// pathWidth is the width of the path (e.g., 50)
// Returns the center position of the path
func centerInPath(pathPos float32, pathWidth float32) float32 {
	return pathPos + pathWidth/2
}

type Game struct {
	maps         []Map
	towers       []Tower
	projectiles  []Projectile
	towerLimit   int
	mousePressed bool
	tick         int // Frame counter
	errorMessage string
	errorTimer   int
}

type Window struct {
	Width  int
	Height int
	Title  string
}

type Path struct {
	StartX float32
	StartY float32
	EndX   float32
	EndY   float32
}

type Map []Path

type Enemy struct {
	PositionX        float32 // This is the CENTER position
	PositionY        float32 // This is the CENTER position
	Speed            float32
	currentPathIndex int
	life             int
}

type Tower struct {
	PositionX    float32 // This is the CENTER position
	PositionY    float32 // This is the CENTER position
	Range        float32
	Damage       int
	FireRate     float32
	lastFireTime int
}

type Projectile struct {
	PositionX float32 // This is the CENTER position
	PositionY float32 // This is the CENTER position
	Speed     int
	Target    *Enemy
}

var firstMap = Map{
	{StartX: 350, StartY: 0, EndX: 350, EndY: 150},
	{StartX: 350, StartY: 150, EndX: 550, EndY: 150},
	{StartX: 550, StartY: 150, EndX: 550, EndY: 350},
	{StartX: 550, StartY: 350, EndX: 150, EndY: 350},
	{StartX: 150, StartY: 350, EndX: 150, EndY: 600},
}

var newWindow = Window{
	Width:  800,
	Height: 600,
	Title:  "Final Path v1.0",
}

var enemy = &Enemy{
	PositionX:        centerInPath(350, PathWidth), // Automatically centered in the path
	PositionY:        centerInPath(0, PathWidth),   // Also center in Y
	Speed:            2,
	currentPathIndex: 0,
	life:             40,
}

func createTower(x, y float32) Tower {
	return Tower{
		PositionX: x,
		PositionY: y,
		Range:     100,
		Damage:    10,
		FireRate:  1,
	}
}

func isPositionOnPath(x, y float32, m Map) bool {
	const margin float32 = 30

	for _, path := range m {
		minX := min(path.StartX, path.EndX) - margin
		maxX := max(path.StartX, path.EndX) + PathWidth + margin
		minY := min(path.StartY, path.EndY) - margin
		maxY := max(path.StartY, path.EndY) + PathWidth + margin

		if x >= minX && x <= maxX && y >= minY && y <= maxY {
			return true
		}
	}

	return false
}

func canPlaceTower(centerX, centerY float32, m Map) bool {
	const halfSize = TowerSize / 2

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
		if isPositionOnPath(point.x, point.y, m) {
			return false
		}
	}

	return true
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func createProjectile(x, y float32, target *Enemy) Projectile {
	return Projectile{
		PositionX: x,
		PositionY: y,
		Speed:     10,
		Target:    target,
	}
}

func createMap(screen *ebiten.Image, m Map) {
	for _, path := range m {
		width := path.EndX - path.StartX
		height := path.EndY - path.StartY

		// If vertical path (width is 0), add PathWidth
		if width == 0 {
			width = PathWidth
		}

		// If horizontal path (height is 0), add PathWidth
		if height == 0 {
			height = PathWidth
		}

		vector.FillRect(screen, path.StartX, path.StartY, width, height, color.White, false)
	}
}

func (g *Game) drawBuildableAreas(screen *ebiten.Image) {
	const gridSize float32 = 40
	screenWidth := float32(screen.Bounds().Dx())
	screenHeight := float32(screen.Bounds().Dy())

	// Draw grid of buildable areas
	for x := float32(0); x < screenWidth; x += gridSize {
		for y := float32(0); y < screenHeight; y += gridSize {
			centerX := x + gridSize/2
			centerY := y + gridSize/2

			if !isPositionOnPath(centerX, centerY, g.maps[0]) {
				// Valid buildable area
				vector.FillRect(screen, x, y, gridSize, gridSize, color.RGBA{0, 100, 0, 30}, false)
				// Draw grid border
				vector.StrokeRect(screen, x, y, gridSize, gridSize, 1, color.RGBA{0, 150, 0, 50}, false)
			}
		}
	}
}

func (projectile *Projectile) hit() bool {
	// Simple straight-line movement towards target
	if projectile.Target == nil || !projectile.Target.isAlive() {
		return false
	}

	// Both positions are already center coordinates, no need to adjust
	dx := projectile.Target.PositionX - projectile.PositionX
	dy := projectile.Target.PositionY - projectile.PositionY
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance < float32(projectile.Speed) {
		// Reached the target
		return true
	} else {
		// Move towards the target
		projectile.PositionX += (dx / distance) * float32(projectile.Speed)
		projectile.PositionY += (dy / distance) * float32(projectile.Speed)
	}

	return false
}

func (tower *Tower) isEnemyInRange(enemy *Enemy) bool {
	dx := tower.PositionX - enemy.PositionX
	dy := tower.PositionY - enemy.PositionY
	distanceSquared := dx*dx + dy*dy
	return distanceSquared <= tower.Range*tower.Range
}

func (tower *Tower) canFire(currentTick int) bool {
	ticksPerShot := int(60.0 / tower.FireRate)
	return currentTick-tower.lastFireTime >= ticksPerShot
}

func (tower *Tower) attack(enemy *Enemy) Projectile {
	return createProjectile(tower.PositionX, tower.PositionY, enemy)
}

func (enemy *Enemy) isAlive() bool {
	return enemy.life > 0
}

func (enemy *Enemy) followPath(m Map) {
	if enemy.currentPathIndex >= len(m) {
		return
	}

	path := m[enemy.currentPathIndex]

	// Vertical movement
	if path.StartX == path.EndX {
		// Center the enemy on X axis for this vertical path
		targetCenterX := centerInPath(path.StartX, PathWidth)
		enemy.PositionX = targetCenterX

		targetCenterY := centerInPath(path.EndY, PathWidth)
		if enemy.PositionY < targetCenterY {
			enemy.PositionY += enemy.Speed
		} else {
			enemy.currentPathIndex++
			fmt.Printf("Path %d completed\n", enemy.currentPathIndex)
		}
	} else if path.StartY == path.EndY { // Horizontal movement
		// Center the enemy on Y axis for this horizontal path
		targetCenterY := centerInPath(path.StartY, PathWidth)
		enemy.PositionY = targetCenterY

		// Move right or left depending on direction
		if path.EndX > path.StartX {
			// Move right
			targetCenterX := centerInPath(path.EndX, PathWidth)
			if enemy.PositionX < targetCenterX {
				enemy.PositionX += enemy.Speed
			} else {
				enemy.currentPathIndex++
				fmt.Printf("Path %d completed\n", enemy.currentPathIndex)
			}
		} else {
			// Move left
			targetCenterX := centerInPath(path.EndX, PathWidth)
			if enemy.PositionX > targetCenterX {
				enemy.PositionX -= enemy.Speed
			} else {
				enemy.currentPathIndex++
				fmt.Printf("Path %d completed\n", enemy.currentPathIndex)
			}
		}
	}
}

func (g *Game) Update() error {
	g.tick++

	if enemy.isAlive() {
		enemy.followPath(firstMap)
	}

	// Check for tower attacks
	for i := range g.towers {
		tower := &g.towers[i]
		if tower.isEnemyInRange(enemy) && enemy.isAlive() && tower.canFire(g.tick) {
			g.projectiles = append(g.projectiles, tower.attack(enemy))
			tower.lastFireTime = g.tick
		}
	}

	// Update projectiles
	var activeProjectiles []Projectile
	for i := range g.projectiles {
		projectile := &g.projectiles[i]
		if projectile.hit() {
			// Projectile hit the target
			if projectile.Target != nil && projectile.Target.isAlive() {
				projectile.Target.life -= 10
				fmt.Printf("Enemy hit! Life: %d\n", projectile.Target.life)
			}
		} else if projectile.Target != nil && projectile.Target.isAlive() {
			// Projectile still moving
			activeProjectiles = append(activeProjectiles, *projectile)
		}
	}
	g.projectiles = activeProjectiles

	// Handle mouse input for placing towers
	mousePressedCurrent := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if mousePressedCurrent && !g.mousePressed && len(g.towers) < g.towerLimit {
		mx, my := ebiten.CursorPosition()
		// Verify if the entire tower can be placed at the position
		if canPlaceTower(float32(mx), float32(my), g.maps[0]) {
			// Add tower at mouse position
			g.towers = append(g.towers, createTower(float32(mx), float32(my)))
		} else {
			g.errorMessage = "Cannot place tower on path!"
			g.errorTimer = 120 // Display only for 120 frames
		}
	}

	// Decrement error message timer
	if g.errorTimer > 0 {
		g.errorTimer--
		if g.errorTimer == 0 {
			g.errorMessage = ""
		}
	}

	// Update mouse pressed state
	g.mousePressed = mousePressedCurrent
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	vector.FillRect(screen, 0, 0, float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy()), color.Black, false)

	// Draw buildable areas (areas where towers can be placed)
	g.drawBuildableAreas(screen)

	// Map
	createMap(screen, g.maps[0])

	// Enemy (draw using center position)
	if enemy.isAlive() {
		topLeftX, topLeftY := CenteredPosition{X: enemy.PositionX, Y: enemy.PositionY, Size: EnemySize}.TopLeft()
		vector.FillRect(screen, topLeftX, topLeftY, EnemySize, EnemySize, color.RGBA{255, 0, 0, 255}, false)
	}

	// Towers (draw using center position)
	for _, tower := range g.towers {
		// Draw range circle centered on tower
		vector.StrokeCircle(screen, tower.PositionX, tower.PositionY, tower.Range, 2, color.RGBA{0, 0, 255, 20}, false)
		topLeftX, topLeftY := CenteredPosition{X: tower.PositionX, Y: tower.PositionY, Size: TowerSize}.TopLeft()
		vector.FillRect(screen, topLeftX, topLeftY, TowerSize, TowerSize, color.RGBA{0, 255, 255, 255}, false)
	}

	// Projectiles (draw centered)
	for _, projectile := range g.projectiles {
		vector.FillCircle(screen, projectile.PositionX, projectile.PositionY, ProjectileSize, color.RGBA{255, 255, 0, 255}, false)
	}

	// Error message
	if g.errorMessage != "" {
		ebitenutil.DebugPrintAt(screen, g.errorMessage, 10, 10)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return newWindow.Width, newWindow.Height
}

func main() {
	g := &Game{
		maps:       []Map{firstMap},
		towerLimit: 3,
	}
	ebiten.SetWindowSize(newWindow.Width, newWindow.Height)
	ebiten.SetWindowTitle(newWindow.Title)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
