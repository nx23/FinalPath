package instructions

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Instructions manages the game instructions screen
type Instructions struct {
	Active       bool
	mousePressed bool
}

// NewInstructions creates a new instructions screen
func NewInstructions() *Instructions {
	return &Instructions{
		Active:       true, // Show on game start
		mousePressed: false,
	}
}

// Update handles input for the instructions screen
// Returns true if the user clicked to close instructions
func (i *Instructions) Update() bool {
	if !i.Active {
		return false
	}

	mousePressedCurrent := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Check for click to close instructions
	if mousePressedCurrent && !i.mousePressed {
		i.mousePressed = false
		i.Active = false
		return true
	}

	i.mousePressed = mousePressedCurrent
	return false
}

// Draw renders the instructions overlay
func (i *Instructions) Draw(screen *ebiten.Image, drawTextFunc func(*ebiten.Image, string, float64, float64, float64)) {
	if !i.Active {
		return
	}

	// Semi-transparent dark overlay
	vector.FillRect(screen, 0, 0, float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy()),
		color.RGBA{0, 0, 0, 220}, false)

	// Instructions panel
	panelX := float32(100)
	panelY := float32(80)
	panelWidth := float32(600)
	panelHeight := float32(560)

	vector.FillRect(screen, panelX, panelY, panelWidth, panelHeight, color.RGBA{30, 30, 30, 255}, false)
	vector.StrokeRect(screen, panelX, panelY, panelWidth, panelHeight, 3, color.RGBA{0, 120, 255, 255}, false)

	// Title
	drawTextFunc(screen, "FINAL PATH - HOW TO PLAY", 170, 100, 3.0)

	// Game objective
	drawTextFunc(screen, "OBJECTIVE:", 120, 160, 2.2)
	drawTextFunc(screen, "Defend your path from waves of enemies!", 140, 190, 1.8)
	drawTextFunc(screen, "Don't let them reach the end!", 140, 215, 1.8)

	// Controls
	drawTextFunc(screen, "CONTROLS:", 120, 265, 2.2)
	drawTextFunc(screen, "LEFT CLICK: Place towers on green areas", 140, 295, 1.8)
	drawTextFunc(screen, "RIGHT CLICK: Remove towers", 140, 320, 1.8)
	drawTextFunc(screen, "SHOP BUTTON: Buy upgrades with coins", 140, 345, 1.8)
	drawTextFunc(screen, "NEXT WAVE: Start the next enemy wave", 140, 370, 1.8)

	// Game mechanics
	drawTextFunc(screen, "MECHANICS:", 120, 420, 2.2)
	drawTextFunc(screen, "- Towers auto-attack enemies in range", 140, 450, 1.8)
	drawTextFunc(screen, "- Earn 10 coins per enemy defeated", 140, 475, 1.8)
	drawTextFunc(screen, "- Lose 1 life if enemy reaches the end", 140, 500, 1.8)
	drawTextFunc(screen, "- Game over when lives reach 0", 140, 525, 1.8)

	// Start button
	buttonX := float32(250)
	buttonY := float32(570)
	buttonWidth := float32(300)
	buttonHeight := float32(50)

	vector.FillRect(screen, buttonX, buttonY, buttonWidth, buttonHeight, color.RGBA{0, 200, 0, 255}, false)
	vector.StrokeRect(screen, buttonX, buttonY, buttonWidth, buttonHeight, 3, color.RGBA{255, 255, 255, 255}, false)
	drawTextFunc(screen, "CLICK TO START", float64(buttonX+50), float64(buttonY+6), 2.5)
}

// Show displays the instructions screen
func (i *Instructions) Show() {
	i.Active = true
}

// Hide closes the instructions screen
func (i *Instructions) Hide() {
	i.Active = false
	i.mousePressed = false
}
