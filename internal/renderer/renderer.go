package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nx23/final-path/internal/config"
	"github.com/nx23/final-path/internal/entity"
	"github.com/nx23/final-path/internal/gamemap"
	"github.com/nx23/final-path/internal/utils"
)

func DrawBuildableAreas(screen *ebiten.Image, gameMap gamemap.Map) {
	const gridSize float32 = 40
	screenWidth := float32(screen.Bounds().Dx())
	screenHeight := float32(screen.Bounds().Dy())

	// Draw grid of buildable areas
	for x := float32(0); x < screenWidth; x += gridSize {
		for y := float32(0); y < screenHeight; y += gridSize {
			centerX := x + gridSize/2
			centerY := y + gridSize/2

			if !gamemap.IsPositionOnPath(centerX, centerY, gameMap) {
				vector.FillRect(screen, x, y, gridSize, gridSize, color.RGBA{0, 100, 0, 30}, false)
				vector.StrokeRect(screen, x, y, gridSize, gridSize, 1, color.RGBA{0, 150, 0, 50}, false)
			}
		}
	}
}

func DrawEnemies(screen *ebiten.Image, enemies []*entity.Enemy) {
	for _, enemy := range enemies {
		if enemy.IsAlive() {
			topLeftX, topLeftY := utils.CenteredPosition{X: enemy.PositionX, Y: enemy.PositionY, Size: config.EnemySize}.TopLeft()
			vector.FillRect(screen, topLeftX, topLeftY, config.EnemySize, config.EnemySize, color.RGBA{255, 0, 0, 255}, false)
		}
	}
}

func DrawTowers(screen *ebiten.Image, towers []entity.Tower) {
	for _, tower := range towers {
		// Draw range circle centered on tower
		vector.StrokeCircle(screen, tower.PositionX, tower.PositionY, tower.Range, 2, color.RGBA{0, 0, 255, 20}, false)
		topLeftX, topLeftY := utils.CenteredPosition{X: tower.PositionX, Y: tower.PositionY, Size: config.TowerSize}.TopLeft()
		vector.FillRect(screen, topLeftX, topLeftY, config.TowerSize, config.TowerSize, color.RGBA{0, 255, 255, 255}, false)
	}
}

func DrawProjectiles(screen *ebiten.Image, projectiles []entity.Projectile) {
	for _, projectile := range projectiles {
		vector.FillCircle(screen, projectile.PositionX, projectile.PositionY, config.ProjectileSize, color.RGBA{255, 255, 0, 255}, false)
	}
}
