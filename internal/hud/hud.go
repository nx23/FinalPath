package hud

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nx23/final-path/internal/config"
)

// HUD represents the game's heads-up display
type HUD struct {
	TowersBuilt     int
	TowersLimit     int
	EnemiesDefeated int
	CurrentWave     int
	WaveActive      bool
	buttonX         float32
	buttonY         float32
	buttonWidth     float32
	buttonHeight    float32
}

// NewHUD creates a new HUD instance
func NewHUD(towerLimit int) *HUD {
	return &HUD{
		TowersBuilt:     0,
		TowersLimit:     towerLimit,
		EnemiesDefeated: 0,
		CurrentWave:     0,
		WaveActive:      false,
		buttonX:         550,
		buttonY:         30,
		buttonWidth:     200,
		buttonHeight:    60,
	}
}

// Draw renders the HUD on the screen
func (h *HUD) Draw(screen *ebiten.Image) {
	// Draw background panel
	screenWidth := float32(screen.Bounds().Dx())
	vector.FillRect(screen, 0, 0, screenWidth, config.HUDHeight, color.RGBA{0, 0, 0, 255}, false)

	// Tower info - larger text
	towerText := fmt.Sprintf("Towers: %d/%d", h.TowersBuilt, h.TowersLimit)
	h.drawLargeText(screen, towerText, 20, 20, 3.0)

	// Enemy info - larger text
	enemyText := fmt.Sprintf("Enemies Defeated: %d", h.EnemiesDefeated)
	h.drawLargeText(screen, enemyText, 20, 65, 3.0)

	// Draw Next Wave button
	h.drawButton(screen)
}

// drawButton draws the "Next Wave" button
func (h *HUD) drawButton(screen *ebiten.Image) {
	var buttonColor color.RGBA
	var buttonText string

	if h.WaveActive {
		// Wave in progress - gray button
		buttonColor = color.RGBA{100, 100, 100, 200}
		buttonText = fmt.Sprintf("Wave %d", h.CurrentWave)
	} else {
		// Ready for next wave - green button
		buttonColor = color.RGBA{0, 200, 0, 200}
		buttonText = "Start Wave"
		if h.CurrentWave > 0 {
			buttonText = "Next Wave"
		}
	}

	// Draw button background
	vector.DrawFilledRect(screen, h.buttonX, h.buttonY, h.buttonWidth, h.buttonHeight, buttonColor, false)

	// Draw button border
	vector.StrokeRect(screen, h.buttonX, h.buttonY, h.buttonWidth, h.buttonHeight, 3, color.RGBA{255, 255, 255, 255}, false)

	// Draw button text (centered)
	h.drawLargeText(screen, buttonText, float64(h.buttonX)+20, float64(h.buttonY)+15, 2.5)
}

// IsButtonClicked checks if the button was clicked at the given coordinates
func (h *HUD) IsButtonClicked(x, y int) bool {
	fx, fy := float32(x), float32(y)
	return fx >= h.buttonX && fx <= h.buttonX+h.buttonWidth &&
		fy >= h.buttonY && fy <= h.buttonY+h.buttonHeight
}

// drawLargeText draws text with actual scaling for better readability
func (h *HUD) drawLargeText(screen *ebiten.Image, text string, x, y, scale float64) {
	// Create a temporary image to render text
	bounds := image.Rect(0, 0, 400, 30)
	textImg := ebiten.NewImage(bounds.Dx(), bounds.Dy())

	// Draw text on temporary image
	ebitenutil.DebugPrint(textImg, text)

	// Scale and draw the text image to the screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)

	screen.DrawImage(textImg, op)
}
