package game

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nx23/final-path/internal/config"
	"github.com/nx23/final-path/internal/entity"
	"github.com/nx23/final-path/internal/gamemap"
	"github.com/nx23/final-path/internal/hud"
	"github.com/nx23/final-path/internal/utils"
)

// Game holds all the game state
type Game struct {
	maps                 []gamemap.Map
	enemies              []*entity.Enemy
	towers               []entity.Tower
	projectiles          []entity.Projectile
	towerLimit           int
	enemiesDefeated      int
	mousePressed         bool
	tick                 int // Frame counter (60 fps)
	errorMessage         string
	errorTimer           int
	hud                  *hud.HUD
	enemiesPerWave       int // Number of enemies in current wave
	enemiesSpawnedInWave int // Number of enemies spawned in current wave
	lastSpawnTick        int // Last tick when an enemy was spawned
	spawnInterval        int // Ticks between enemy spawns (60 ticks = 1 second)
	lives                int // Player lives
	gameOver             bool
	restartButtonX       float32
	restartButtonY       float32
	restartButtonWidth   float32
	restartButtonHeight  float32
}

// NewGame initializes a new game with the default map
func NewGame() *Game {
	gameMap := gamemap.DefaultMap()
	towerLimit := 3
	initialLives := 10

	g := &Game{
		maps:                []gamemap.Map{gameMap},
		enemies:             []*entity.Enemy{},
		towerLimit:          towerLimit,
		enemiesDefeated:     0,
		hud:                 hud.NewHUD(towerLimit),
		spawnInterval:       60, // 1 second between spawns (60 fps * 1)
		lives:               initialLives,
		gameOver:            false,
		restartButtonX:      300,
		restartButtonY:      360,
		restartButtonWidth:  200,
		restartButtonHeight: 60,
	}

	// Sync lives with HUD
	g.hud.Lives = initialLives
	return g
}

// Update is called every frame (60x per second) to update the game state
func (g *Game) Update() error {
	g.tick++

	// Handle game over state
	if g.gameOver {
		return g.handleGameOverInput()
	}

	// Update enemy movement only if wave is active
	if g.hud.WaveActive {
		// Spawn new enemies at intervals
		if g.enemiesSpawnedInWave < g.enemiesPerWave {
			if g.tick-g.lastSpawnTick >= g.spawnInterval || g.enemiesSpawnedInWave == 0 {
				// Spawn a new enemy
				g.enemies = append(g.enemies, entity.NewEnemy(g.maps[0]))
				g.enemiesSpawnedInWave++
				g.lastSpawnTick = g.tick
				fmt.Printf("Enemy spawned! (%d/%d)\n", g.enemiesSpawnedInWave, g.enemiesPerWave)
			}
		}

		// Update all enemies
		var aliveEnemies []*entity.Enemy
		for _, enemy := range g.enemies {
			if enemy.IsAlive() {
				// Check if enemy reached the end of the path
				if enemy.CurrentPathIndex >= len(g.maps[0]) {
					// Enemy escaped! Lose a life
					g.lives--
					g.hud.Lives = g.lives
					fmt.Printf("Enemy escaped! Lives remaining: %d\n", g.lives)

					// Check for game over
					if g.lives <= 0 {
						g.gameOver = true
						g.hud.WaveActive = false
						fmt.Println("Game Over!")
					}
					// Don't add to aliveEnemies (despawn)
				} else {
					enemy.FollowPath(g.maps[0])
					aliveEnemies = append(aliveEnemies, enemy)
				}
			} else {
				// Enemy just died
				g.enemiesDefeated++
				g.hud.EnemiesDefeated = g.enemiesDefeated
				g.hud.EnemiesKilledInWave++
				fmt.Printf("Enemy defeated! Total: %d\n", g.enemiesDefeated)
			}
		}
		g.enemies = aliveEnemies

		// Check if wave is complete (all enemies spawned and all dead)
		if g.enemiesSpawnedInWave >= g.enemiesPerWave && len(g.enemies) == 0 {
			g.hud.WaveActive = false
			g.hud.EnemiesKilledInWave = 0
			// Calculate enemies for next wave
			nextWaveEnemies := 3 + g.hud.CurrentWave*2
			g.hud.EnemiesInWave = nextWaveEnemies
			fmt.Printf("Wave %d complete!\n", g.hud.CurrentWave)
		}

		// Check for tower attacks on all enemies
		for i := range g.towers {
			tower := &g.towers[i]
			if tower.CanFire(g.tick) {
				// Find closest enemy in range
				for _, enemy := range g.enemies {
					if tower.IsEnemyInRange(enemy) && enemy.IsAlive() {
						g.projectiles = append(g.projectiles, tower.Attack(enemy))
						tower.LastFireTime = g.tick
						break // Only attack one enemy per tower per fire cycle
					}
				}
			}
		}

		// Update projectiles
		var activeProjectiles []entity.Projectile
		for i := range g.projectiles {
			projectile := &g.projectiles[i]
			if projectile.Hit() {
				// Projectile hit the target
				if projectile.Target != nil && projectile.Target.IsAlive() {
					projectile.Target.TakeDamage(10)
					fmt.Printf("Enemy hit! Life: %d\n", projectile.Target.Life)
				}
			} else if projectile.Target != nil && projectile.Target.IsAlive() {
				// Projectile still moving
				activeProjectiles = append(activeProjectiles, *projectile)
			}
		}
		g.projectiles = activeProjectiles
	}

	// Handle mouse input
	g.handleMouseInput()

	// Decrement error message timer
	if g.errorTimer > 0 {
		g.errorTimer--
		if g.errorTimer == 0 {
			g.errorMessage = ""
		}
	}

	return nil
}

// handleMouseInput handles all mouse interactions
func (g *Game) handleMouseInput() {
	mousePressedCurrent := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if mousePressedCurrent && !g.mousePressed {
		mx, my := ebiten.CursorPosition()

		// Check if clicking the Next Wave button
		if g.hud.IsButtonClicked(mx, my) && !g.hud.WaveActive {
			g.startNextWave()
		} else if len(g.towers) < g.towerLimit {
			// Try to place a tower
			g.placeTower(float32(mx), float32(my))
		}
	}

	// Update mouse pressed state
	g.mousePressed = mousePressedCurrent
}

// startNextWave starts the next wave of enemies
func (g *Game) startNextWave() {
	g.hud.CurrentWave++
	g.hud.WaveActive = true

	// Set number of enemies for this wave (increases with wave number)
	g.enemiesPerWave = 3 + (g.hud.CurrentWave-1)*2
	g.enemiesSpawnedInWave = 0
	g.lastSpawnTick = g.tick - g.spawnInterval // Allow immediate spawn

	// Update HUD with wave info
	g.hud.EnemiesInWave = g.enemiesPerWave
	g.hud.EnemiesKilledInWave = 0

	fmt.Printf("Wave %d started! (%d enemies)\n", g.hud.CurrentWave, g.enemiesPerWave)
}

// placeTower attempts to place a tower at the given position
func (g *Game) placeTower(x, y float32) {
	// Check if clicking in HUD area
	if y < config.HUDHeight {
		g.errorMessage = "Cannot place tower in HUD area!"
		g.errorTimer = 120
		return
	}

	if entity.CanPlaceTower(x, y, g.maps[0]) {
		g.towers = append(g.towers, entity.NewTower(x, y))
		g.hud.TowersBuilt = len(g.towers)
	} else {
		g.errorMessage = "Cannot place tower on path!"
		g.errorTimer = 120
	}
}

// Draw renders everything on the screen every frame
func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	vector.FillRect(screen, 0, 0, float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy()), color.Black, false)

	// Draw buildable areas
	g.drawBuildableAreas(screen)

	// Draw map
	g.maps[0].Draw(screen)

	// Draw all enemies
	for _, enemy := range g.enemies {
		if enemy.IsAlive() {
			topLeftX, topLeftY := utils.CenteredPosition{X: enemy.PositionX, Y: enemy.PositionY, Size: config.EnemySize}.TopLeft()
			vector.FillRect(screen, topLeftX, topLeftY, config.EnemySize, config.EnemySize, color.RGBA{255, 0, 0, 255}, false)
		}
	}

	// Draw towers
	for _, tower := range g.towers {
		// Draw range circle centered on tower
		vector.StrokeCircle(screen, tower.PositionX, tower.PositionY, tower.Range, 2, color.RGBA{0, 0, 255, 20}, false)
		topLeftX, topLeftY := utils.CenteredPosition{X: tower.PositionX, Y: tower.PositionY, Size: config.TowerSize}.TopLeft()
		vector.FillRect(screen, topLeftX, topLeftY, config.TowerSize, config.TowerSize, color.RGBA{0, 255, 255, 255}, false)
	}

	// Draw projectiles
	for _, projectile := range g.projectiles {
		vector.FillCircle(screen, projectile.PositionX, projectile.PositionY, config.ProjectileSize, color.RGBA{255, 255, 0, 255}, false)
	}

	// Draw HUD
	g.hud.Draw(screen)

	// Draw game over screen if game is over
	if g.gameOver {
		g.drawGameOverScreen(screen)
	}

	// Draw error message (below HUD, larger text)
	if g.errorMessage != "" {
		g.drawLargeText(screen, g.errorMessage, 20, float64(config.HUDHeight)+10, 1.5)
	}
}

// drawBuildableAreas draws a green grid showing where towers can be placed
func (g *Game) drawBuildableAreas(screen *ebiten.Image) {
	const gridSize float32 = 40
	screenWidth := float32(screen.Bounds().Dx())
	screenHeight := float32(screen.Bounds().Dy())

	// Draw grid of buildable areas
	for x := float32(0); x < screenWidth; x += gridSize {
		for y := float32(0); y < screenHeight; y += gridSize {
			centerX := x + gridSize/2
			centerY := y + gridSize/2

			if !gamemap.IsPositionOnPath(centerX, centerY, g.maps[0]) {
				// Valid buildable area
				vector.FillRect(screen, x, y, gridSize, gridSize, color.RGBA{0, 100, 0, 30}, false)
				// Draw grid border
				vector.StrokeRect(screen, x, y, gridSize, gridSize, 1, color.RGBA{0, 150, 0, 50}, false)
			}
		}
	}
}

// Layout defines the game's logical screen size (required by ebiten.Game interface)
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.Config.Width, config.Config.Height
}

// handleGameOverInput handles mouse input during game over state
func (g *Game) handleGameOverInput() error {
	mousePressedCurrent := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if mousePressedCurrent && !g.mousePressed {
		mx, my := ebiten.CursorPosition()

		// Check if clicking the restart button
		if g.isRestartButtonClicked(mx, my) {
			g.restartGame()
		}
	}

	g.mousePressed = mousePressedCurrent
	return nil
}

// isRestartButtonClicked checks if the restart button was clicked
func (g *Game) isRestartButtonClicked(x, y int) bool {
	fx, fy := float32(x), float32(y)
	return fx >= g.restartButtonX && fx <= g.restartButtonX+g.restartButtonWidth &&
		fy >= g.restartButtonY && fy <= g.restartButtonY+g.restartButtonHeight
}

// restartGame resets the game to initial state
func (g *Game) restartGame() {
	fmt.Println("Restarting game...")
	g.enemies = []*entity.Enemy{}
	g.towers = []entity.Tower{}
	g.projectiles = []entity.Projectile{}
	g.enemiesDefeated = 0
	g.lives = 10
	g.gameOver = false
	g.tick = 0
	g.enemiesPerWave = 0
	g.enemiesSpawnedInWave = 0
	g.lastSpawnTick = 0
	g.errorMessage = ""
	g.errorTimer = 0

	// Reset HUD
	g.hud.TowersBuilt = 0
	g.hud.EnemiesDefeated = 0
	g.hud.CurrentWave = 0
	g.hud.WaveActive = false
	g.hud.EnemiesInWave = 3
	g.hud.EnemiesKilledInWave = 0
	g.hud.Lives = 10
}

// drawGameOverScreen draws the game over overlay with restart button
func (g *Game) drawGameOverScreen(screen *ebiten.Image) {
	// Semi-transparent dark overlay
	vector.FillRect(screen, 0, 0, float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy()),
		color.RGBA{0, 0, 0, 200}, false)

	// Game Over text (large)
	g.drawLargeText(screen, "GAME OVER!", 280, 250, 4.0)

	// Final score
	scoreText := fmt.Sprintf("Enemies Defeated: %d", g.enemiesDefeated)
	g.drawLargeText(screen, scoreText, 240, 310, 2.5)

	// Restart button
	buttonColor := color.RGBA{0, 200, 0, 255}
	vector.FillRect(screen, g.restartButtonX, g.restartButtonY,
		g.restartButtonWidth, g.restartButtonHeight, buttonColor, false)

	// Button border
	vector.StrokeRect(screen, g.restartButtonX, g.restartButtonY,
		g.restartButtonWidth, g.restartButtonHeight, 3, color.RGBA{255, 255, 255, 255}, false)

	// Button text (centered)
	g.drawLargeText(screen, "RESTART", float64(g.restartButtonX+50), float64(g.restartButtonY+15), 2.5)
}

// drawLargeText draws text with actual scaling for better readability
func (g *Game) drawLargeText(screen *ebiten.Image, text string, x, y, scale float64) {
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
