package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	maps         []Map
	towers       []Tower
	towerLimit   int
	mousePressed bool
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
	PositionX        float32
	PositionY        float32
	Speed            float32
	currentPathIndex int
	life             int
}

type Tower struct {
	PositionX float32
	PositionY float32
	Range     float32
	Damage    int
	FireRate  float32
}

var firstMap = Map{
	{StartX: 350, StartY: 0, EndX: 350, EndY: 150},
	{StartX: 350, StartY: 150, EndX: 550, EndY: 150},
	{StartX: 550, StartY: 150, EndX: 550, EndY: 350},
	{StartX: 600, StartY: 350, EndX: 150, EndY: 350},
	{StartX: 150, StartY: 350, EndX: 150, EndY: 600},
}

var newWindow = Window{
	Width:  800,
	Height: 600,
	Title:  "Final Path v1.0",
}

var enemy = &Enemy{
	PositionX:        362.5, // Centered in the path (350 + 12.5)
	PositionY:        0,
	Speed:            2,
	currentPathIndex: 0,
	life:             100,
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

func createMap(screen *ebiten.Image, m Map) {
	const pathWidth float32 = 50
	for _, path := range m {
		width := path.EndX - path.StartX
		height := path.EndY - path.StartY

		// If vertical path (width is 0), add pathWidth
		if width == 0 {
			width = pathWidth
		}

		// If horizontal path (height is 0), add pathWidth
		if height == 0 {
			height = pathWidth
		}

		vector.FillRect(screen, path.StartX, path.StartY, width, height, color.White, false)
	}
}

func (enemy *Enemy) followPath(m Map) {
	if enemy.currentPathIndex >= len(m) {
		return
	}

	const pathWidth float32 = 50
	const enemySize float32 = 25
	const offset = (pathWidth - enemySize) / 2 // Center offset: 12.5

	path := m[enemy.currentPathIndex]

	// Vertical movement
	if path.StartX == path.EndX {
		if enemy.PositionY-offset < path.EndY {
			enemy.PositionY += enemy.Speed
		} else {
			enemy.currentPathIndex++
			fmt.Printf("Path %d completed\n", enemy.currentPathIndex)
		}
	} else if path.StartY == path.EndY { // Horizontal movement
		// Move right or left depending on direction
		if path.EndX > path.StartX {
			// Move right
			if enemy.PositionX-offset < path.EndX {
				enemy.PositionX += enemy.Speed
			} else {
				enemy.currentPathIndex++
				fmt.Printf("Path %d completed\n", enemy.currentPathIndex)
			}
		} else {
			// Move left
			if enemy.PositionX-offset > path.EndX {
				enemy.PositionX -= enemy.Speed
			} else {
				enemy.currentPathIndex++
				fmt.Printf("Path %d completed\n", enemy.currentPathIndex)
			}
		}
	}
}

func (g *Game) Update() error {
	enemy.followPath(firstMap)

	mousePressedCurrent := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if mousePressedCurrent && !g.mousePressed && len(g.towers) < g.towerLimit {
		mx, my := ebiten.CursorPosition()
		// Add tower at mouse position
		g.towers = append(g.towers, createTower(float32(mx), float32(my)))
	}

	// Update mouse pressed state
	g.mousePressed = mousePressedCurrent
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	vector.FillRect(screen, 0, 0, float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy()), color.Black, false)

	// Map
	createMap(screen, g.maps[0])

	// Enemy
	vector.FillRect(screen, enemy.PositionX, enemy.PositionY, 25, 25, color.RGBA{255, 0, 0, 255}, false)

	// Towers
	for _, tower := range g.towers {
		// Draw range circle centered on tower
		vector.FillCircle(screen, tower.PositionX, tower.PositionY, tower.Range, color.RGBA{0, 0, 255, 20}, false)
		// Draw tower square (25x25) centered on position
		vector.FillRect(screen, tower.PositionX-12.5, tower.PositionY-12.5, 25, 25, color.RGBA{0, 255, 255, 255}, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return newWindow.Width, newWindow.Height
}

func main() {
	g := &Game{
		maps:       []Map{firstMap},
		towerLimit: 1,
	}
	ebiten.SetWindowSize(newWindow.Width, newWindow.Height)
	ebiten.SetWindowTitle(newWindow.Title)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
