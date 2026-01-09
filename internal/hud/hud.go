package hud

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nx23/final-path/internal/config"
	"github.com/nx23/final-path/internal/utils"
)

type HUD struct {
	TowersBuilt         int
	TowersLimit         int
	TowerCost           int
	TowerRefund         int
	EnemiesDefeated     int
	CurrentWave         int
	WaveActive          bool
	EnemiesInWave       int
	EnemiesKilledInWave int
	Lives               int
	Coins               int
	buttonX             float32
	buttonY             float32
	buttonWidth         float32
	buttonHeight        float32
	shopButtonX         float32
	shopButtonY         float32
	shopButtonWidth     float32
	shopButtonHeight    float32
}

func NewHUD(towerLimit int, towerCost int, towerRefund int, initialLives int, initialCoins int) *HUD {
	return &HUD{
		TowersBuilt:         0,
		TowersLimit:         towerLimit,
		TowerCost:           towerCost,
		TowerRefund:         towerRefund,
		EnemiesDefeated:     0,
		CurrentWave:         0,
		WaveActive:          false,
		EnemiesInWave:       3,
		EnemiesKilledInWave: 0,
		Lives:               initialLives,
		Coins:               initialCoins,
		buttonX:             620,
		buttonY:             35,
		buttonWidth:         150,
		buttonHeight:        50,
		shopButtonX:         620,
		shopButtonY:         90,
		shopButtonWidth:     150,
		shopButtonHeight:    25,
	}
}

func (h *HUD) Draw(screen *ebiten.Image) {
	// Draw background panel
	screenWidth := float32(screen.Bounds().Dx())
	vector.FillRect(screen, 0, 0, screenWidth, config.HUDHeight, color.RGBA{0, 0, 0, 255}, false)

	// Tower info
	towerText := fmt.Sprintf("Towers Placed: %d/%d", h.TowersBuilt, h.TowersLimit)
	utils.DrawLargeText(screen, towerText, 20, 10, 2.0)

	// Tower costs info
	costText := fmt.Sprintf("Tower Cost: %d coins", h.TowerCost)
	utils.DrawLargeText(screen, costText, 20, 45, 2.0)

	// Tower refund info
	refundText := fmt.Sprintf("Tower Refund: %d coins", h.TowerRefund)
	utils.DrawLargeText(screen, refundText, 20, 80, 2.0)

	// Wave progress info (when active)
	if h.WaveActive {
		waveProgressText := fmt.Sprintf("Wave %d: %d/%d", h.CurrentWave, h.EnemiesKilledInWave, h.EnemiesInWave)
		utils.DrawLargeText(screen, waveProgressText, 350, 10, 2.0)
	}

	// Next wave preview (when not active)
	if !h.WaveActive && h.EnemiesInWave > 0 {
		nextWaveText := fmt.Sprintf("Next Wave: %d enemies", h.EnemiesInWave)
		utils.DrawLargeText(screen, nextWaveText, 350, 10, 2.0)
	}

	// Coins info
	coinsText := fmt.Sprintf("Coins: %d", h.Coins)
	utils.DrawLargeText(screen, coinsText, 350, 45, 2.0)

	// Lives info
	livesText := fmt.Sprintf("Lives: %d", h.Lives)
	utils.DrawLargeText(screen, livesText, 350, 80, 2.0)

	// Draw Next Wave button
	h.drawButton(screen)

	// Draw Shop button
	h.drawShopButton(screen)
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
		// Ready for next wave - blue button
		buttonColor = color.RGBA{0, 120, 255, 220}
		if h.CurrentWave == 0 {
			buttonText = "Start Wave"
		} else {
			buttonText = "Next Wave"
		}
	}

	// Draw button background
	vector.FillRect(screen, h.buttonX, h.buttonY, h.buttonWidth, h.buttonHeight, buttonColor, false)

	// Draw button border
	vector.StrokeRect(screen, h.buttonX, h.buttonY, h.buttonWidth, h.buttonHeight, 3, color.RGBA{255, 255, 255, 255}, false)

	// Draw button text (centered)
	utils.DrawLargeText(screen, buttonText, float64(h.buttonX)+9, float64(h.buttonY)+12, 2.2)
}

// IsButtonClicked checks if the button was clicked at the given coordinates
func (h *HUD) IsButtonClicked(x, y int) bool {
	fx, fy := float32(x), float32(y)
	return fx >= h.buttonX && fx <= h.buttonX+h.buttonWidth &&
		fy >= h.buttonY && fy <= h.buttonY+h.buttonHeight
}

// drawShopButton draws the shop button
func (h *HUD) drawShopButton(screen *ebiten.Image) {
	buttonColor := color.RGBA{255, 165, 0, 220} // Orange color

	// Draw button background
	vector.FillRect(screen, h.shopButtonX, h.shopButtonY, h.shopButtonWidth, h.shopButtonHeight, buttonColor, false)

	// Draw button border
	vector.StrokeRect(screen, h.shopButtonX, h.shopButtonY, h.shopButtonWidth, h.shopButtonHeight, 2, color.RGBA{255, 255, 255, 255}, false)

	// Draw button text
	utils.DrawLargeText(screen, "SHOP", float64((h.shopButtonWidth/2)+h.shopButtonX-15), float64(h.shopButtonY), 1.5)
}

// IsShopButtonClicked checks if the shop button was clicked
func (h *HUD) IsShopButtonClicked(x, y int) bool {
	fx, fy := float32(x), float32(y)
	return fx >= h.shopButtonX && fx <= h.shopButtonX+h.shopButtonWidth &&
		fy >= h.shopButtonY && fy <= h.shopButtonY+h.shopButtonHeight
}
