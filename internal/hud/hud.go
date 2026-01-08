package hud

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// HUD represents the game's heads-up display
type HUD struct {
	TowersBuilt     int
	TowersLimit     int
	EnemiesDefeated int
}

// Draw renders the HUD on the screen
func (h *HUD) Draw(screen *ebiten.Image) {
	// Draw semi-transparent background panel
	panelHeight := float32(80)
	screenWidth := float32(screen.Bounds().Dx())
	vector.DrawFilledRect(screen, 0, 0, screenWidth, panelHeight, color.RGBA{0, 0, 0, 180}, false)

	// Tower info
	towerText := fmt.Sprintf("Towers: %d/%d", h.TowersBuilt, h.TowersLimit)
	ebitenutil.DebugPrintAt(screen, towerText, 15, 15)

	// Enemy info
	enemyText := fmt.Sprintf("Enemies Defeated: %d", h.EnemiesDefeated)
	ebitenutil.DebugPrintAt(screen, enemyText, 15, 40)
}
