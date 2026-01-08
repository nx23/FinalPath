package gameover

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// GameOver manages the game over screen and restart functionality
type GameOver struct {
	Active              bool
	RestartButtonX      float32
	RestartButtonY      float32
	RestartButtonWidth  float32
	RestartButtonHeight float32
	mousePressed        bool
}

// NewGameOver creates a new game over screen
func NewGameOver() *GameOver {
	return &GameOver{
		Active:              false,
		RestartButtonX:      300,
		RestartButtonY:      360,
		RestartButtonWidth:  200,
		RestartButtonHeight: 60,
		mousePressed:        false,
	}
}

// Draw renders the game over screen with restart button
func (go_screen *GameOver) Draw(screen *ebiten.Image, enemiesDefeated int, drawTextFunc func(*ebiten.Image, string, float64, float64, float64)) {
	if !go_screen.Active {
		return
	}

	// Semi-transparent dark overlay
	vector.FillRect(screen, 0, 0, float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy()),
		color.RGBA{0, 0, 0, 200}, false)

	// Game Over text (large)
	drawTextFunc(screen, "GAME OVER!", 280, 250, 4.0)

	// Final score
	scoreText := fmt.Sprintf("Enemies Defeated: %d", enemiesDefeated)
	drawTextFunc(screen, scoreText, 240, 310, 2.5)

	// Restart button
	buttonColor := color.RGBA{0, 200, 0, 255}
	vector.FillRect(screen, go_screen.RestartButtonX, go_screen.RestartButtonY,
		go_screen.RestartButtonWidth, go_screen.RestartButtonHeight, buttonColor, false)

	// Button border
	vector.StrokeRect(screen, go_screen.RestartButtonX, go_screen.RestartButtonY,
		go_screen.RestartButtonWidth, go_screen.RestartButtonHeight, 3, color.RGBA{255, 255, 255, 255}, false)

	// Button text (centered)
	drawTextFunc(screen, "RESTART", float64(go_screen.RestartButtonX+50), float64(go_screen.RestartButtonY+15), 2.5)
}

// Update handles input for the game over screen
// Returns true if the restart button was clicked
func (go_screen *GameOver) Update() bool {
	if !go_screen.Active {
		return false
	}

	// Handle mouse input
	mousePressedCurrent := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Check for mouse click (press and release)
	if !mousePressedCurrent && go_screen.mousePressed {
		mx, my := ebiten.CursorPosition()

		// Check if clicking the restart button
		if go_screen.isRestartButtonClicked(mx, my) {
			go_screen.mousePressed = false
			return true
		}
	}

	go_screen.mousePressed = mousePressedCurrent
	return false
}

// isRestartButtonClicked checks if the restart button was clicked
func (go_screen *GameOver) isRestartButtonClicked(x, y int) bool {
	fx, fy := float32(x), float32(y)
	return fx >= go_screen.RestartButtonX && fx <= go_screen.RestartButtonX+go_screen.RestartButtonWidth &&
		fy >= go_screen.RestartButtonY && fy <= go_screen.RestartButtonY+go_screen.RestartButtonHeight
}

// Activate triggers the game over screen
func (go_screen *GameOver) Activate() {
	go_screen.Active = true
}

// Reset deactivates the game over screen
func (go_screen *GameOver) Reset() {
	go_screen.Active = false
	go_screen.mousePressed = false
}
