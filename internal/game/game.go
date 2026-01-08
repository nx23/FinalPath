package game

import (
	"fmt"
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
	maps            []gamemap.Map
	enemy           *entity.Enemy
	towers          []entity.Tower
	projectiles     []entity.Projectile
	towerLimit      int
	enemiesDefeated int
	enemyJustDied   bool
	mousePressed    bool
	tick            int // Frame counter (60 fps)
	errorMessage    string
	errorTimer      int
	hud             *hud.HUD
}

// NewGame initializes a new game with the default map
func NewGame() *Game {
	gameMap := gamemap.DefaultMap()
	towerLimit := 3
	return &Game{
		maps:            []gamemap.Map{gameMap},
		enemy:           entity.NewEnemy(gameMap),
		towerLimit:      towerLimit,
		enemiesDefeated: 0,
		hud:             hud.NewHUD(towerLimit),
	}
}

// Update is called every frame (60x per second) to update the game state
func (g *Game) Update() error {
	g.tick++

	// Update enemy movement only if wave is active
	if g.hud.WaveActive {
		if g.enemy.IsAlive() {
			g.enemy.FollowPath(g.maps[0])
			g.enemyJustDied = false
		} else if !g.enemyJustDied {
			// Enemy just died, respawn it
			g.enemy = entity.NewEnemy(g.maps[0])
			g.enemiesDefeated++
			g.hud.EnemiesDefeated = g.enemiesDefeated
			g.enemyJustDied = true
			fmt.Printf("Enemy defeated! Total: %d\n", g.enemiesDefeated)
		}

		// Check for tower attacks
		for i := range g.towers {
			tower := &g.towers[i]
			if tower.IsEnemyInRange(g.enemy) && g.enemy.IsAlive() && tower.CanFire(g.tick) {
				g.projectiles = append(g.projectiles, tower.Attack(g.enemy))
				tower.LastFireTime = g.tick
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
	g.enemy = entity.NewEnemy(g.maps[0])
	g.enemyJustDied = false
	fmt.Printf("Wave %d started!\n", g.hud.CurrentWave)
}

// placeTower attempts to place a tower at the given position
func (g *Game) placeTower(x, y float32) {
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

	// Draw enemy
	if g.enemy.IsAlive() {
		topLeftX, topLeftY := utils.CenteredPosition{X: g.enemy.PositionX, Y: g.enemy.PositionY, Size: config.EnemySize}.TopLeft()
		vector.FillRect(screen, topLeftX, topLeftY, config.EnemySize, config.EnemySize, color.RGBA{255, 0, 0, 255}, false)
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

	// Draw error message (below HUD)
	if g.errorMessage != "" {
		ebitenutil.DebugPrintAt(screen, g.errorMessage, 20, int(config.HUDHeight)+10)
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
