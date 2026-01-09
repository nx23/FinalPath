package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nx23/final-path/internal/config"
	"github.com/nx23/final-path/internal/entity"
	"github.com/nx23/final-path/internal/gamemap"
	"github.com/nx23/final-path/internal/gameover"
	"github.com/nx23/final-path/internal/hud"
	"github.com/nx23/final-path/internal/instructions"
	"github.com/nx23/final-path/internal/renderer"
	"github.com/nx23/final-path/internal/shop"
	"github.com/nx23/final-path/internal/utils"
)

type Game struct {
	maps                   []gamemap.Map
	enemies                []*entity.Enemy
	towers                 []entity.Tower
	projectiles            []entity.Projectile
	towerLimit             int
	towerCost              int
	towerRefund            int
	enemiesDefeated        int
	difficultyModifier     int
	mousePressed           bool
	mouseRightPressed      bool
	tick                   int
	errorMessage           string
	errorTimer             int
	hud                    *hud.HUD
	enemiesPerWave         int
	enemiesSpawnedInWave   int
	lastSpawnTick          int
	spawnInterval          int
	lives                  int
	coins                  int
	towerDamageBoost       int
	towerDamageBoostCost   int
	towerFireRateBoost     float32
	towerFireRateBoostCost int
	shop                   *shop.Shop
	gameOverScreen         *gameover.GameOver
	instructionsScreen     *instructions.Instructions
}

// NewGame initializes a new game with the default map
func NewGame() *Game {
	gameMap := gamemap.DefaultMap()

	g := &Game{
		maps:               []gamemap.Map{gameMap},
		enemies:            []*entity.Enemy{},
		towerLimit:         config.GameConstants.TowerLimit,
		towerCost:          config.GameConstants.InitialTowerCost,
		towerRefund:        config.GameConstants.InitialTowerRefund,
		enemiesDefeated:    config.GameConstants.EnemiesDefeated,
		difficultyModifier: config.GameConstants.DifficultyModifier,
		hud:                hud.NewHUD(config.GameConstants.TowerLimit, config.GameConstants.InitialTowerCost, config.GameConstants.InitialTowerRefund, config.GameConstants.InitialLives, config.GameConstants.InitialCoins),
		spawnInterval:      config.GameConstants.SpawnInterval,
		lives:              config.GameConstants.InitialLives,
		coins:              config.GameConstants.InitialCoins,
		towerDamageBoost:   config.GameConstants.TowerDamageBoost,
		towerFireRateBoost: config.GameConstants.TowerFireRateBoost,
		shop:               shop.NewShop(),
		gameOverScreen:     gameover.NewGameOver(),
		instructionsScreen: instructions.NewInstructions(),
	}

	return g
}

func (g *Game) Update() error {
	g.tick++

	if g.instructionsScreen.Active {
		if g.instructionsScreen.Update() {
			// Instructions were just closed, consume the click to prevent tower placement
			g.mousePressed = true
		}
		return nil
	}

	// Handle game over state
	if g.gameOverScreen.Active {
		if g.gameOverScreen.Update() {
			g.restartGame()
		}
		return nil
	}

	// Update enemy movement only if wave is active
	if g.hud.WaveActive {
		if g.enemiesSpawnedInWave < g.enemiesPerWave {
			if g.tick-g.lastSpawnTick >= g.spawnInterval || g.enemiesSpawnedInWave == 0 {
				g.enemies = append(g.enemies, entity.NewEnemy(entity.NewEnemyParams{
					Map:   g.maps[0],
					Speed: 2 * (1 + float32(g.hud.CurrentWave-1)*0.1),
					Life:  10 + (1 + (g.hud.CurrentWave-1)*2) + (20 * g.difficultyModifier),
				}))
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
					g.lives--
					g.hud.Lives = g.lives
					fmt.Printf("Enemy escaped! Lives remaining: %d\n", g.lives)

					if g.lives <= 0 {
						g.gameOverScreen.Activate()
						g.hud.WaveActive = false
						fmt.Println("Game Over!")
					}
				} else {
					enemy.FollowPath(g.maps[0])
					aliveEnemies = append(aliveEnemies, enemy)
				}
			} else {
				g.enemiesDefeated++
				g.coins += 5
				g.hud.EnemiesDefeated = g.enemiesDefeated
				g.hud.Coins = g.coins
				g.hud.EnemiesKilledInWave++
				fmt.Printf("Enemy defeated! Total: %d, Coins: %d\n", g.enemiesDefeated, g.coins)
			}
		}
		g.enemies = aliveEnemies

		// Check if wave is complete (all enemies spawned and all dead)
		if g.enemiesSpawnedInWave >= g.enemiesPerWave && len(g.enemies) == 0 {
			g.hud.WaveActive = false
			g.hud.EnemiesKilledInWave = 0
			// Calculate enemies for next wave
			nextWaveEnemies := 3 + g.hud.CurrentWave*2 + (g.difficultyModifier-1)*2
			g.hud.EnemiesInWave = nextWaveEnemies
			fmt.Printf("Wave %d complete!\n", g.hud.CurrentWave)
		}

		// Check for tower attacks on all enemies
		for i := range g.towers {
			tower := &g.towers[i]
			// Apply global fire rate boost
			boostedFireRate := tower.FireRate * g.towerFireRateBoost
			ticksPerShot := int(60.0 / boostedFireRate)
			canFire := g.tick-tower.LastFireTime >= ticksPerShot

			if canFire {
				// Find closest enemy in range
				for _, enemy := range g.enemies {
					if tower.IsEnemyInRange(enemy) && enemy.IsAlive() {
						g.projectiles = append(g.projectiles, tower.Attack(enemy))
						g.towers[i].LastFireTime = g.tick
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

				if projectile.Target != nil && projectile.Target.IsAlive() {
					totalDamage := 10 + g.towerDamageBoost
					projectile.Target.TakeDamage(totalDamage)
					fmt.Printf("Enemy hit! Damage: %d, Life: %d\n", totalDamage, projectile.Target.Life)
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
	mouseRightPressedCurrent := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	if mousePressedCurrent && !g.mousePressed {
		mx, my := ebiten.CursorPosition()

		if g.hud.IsShopButtonClicked(mx, my) {
			g.shop.Toggle()
			fmt.Printf("Shop %s\n", map[bool]string{true: "opened", false: "closed"}[g.shop.Open])
		} else if g.shop.Open {
			// Handle shop item clicks
			g.handleShopClick(mx, my)
		} else if g.hud.IsButtonClicked(mx, my) && !g.hud.WaveActive {
			// Check if clicking the Next Wave button
			g.startNextWave()
		} else {
			// Try to place a tower
			g.placeTower(float32(mx), float32(my))
		}

	}

	// Handle right click (remove tower or close shop)
	if mouseRightPressedCurrent && !g.mouseRightPressed {
		if g.shop.Open {
			g.shop.Close()
		} else {
			mx, my := ebiten.CursorPosition()
			g.removeTower(float32(mx), float32(my))
		}
	}

	// Update mouse pressed states
	g.mousePressed = mousePressedCurrent
	g.mouseRightPressed = mouseRightPressedCurrent
}

// startNextWave starts the next wave of enemies
func (g *Game) startNextWave() {
	g.hud.CurrentWave++
	g.hud.WaveActive = true

	// Set number of enemies for this wave (increases with wave number)
	g.enemiesPerWave = 3 + (g.hud.CurrentWave-1)*2
	g.enemiesSpawnedInWave = 0
	g.lastSpawnTick = g.tick - g.spawnInterval - (g.hud.CurrentWave * 2)

	// Update difficulty modifier every 5 waves
	if g.hud.CurrentWave%5 == 0 {
		g.difficultyModifier++
		fmt.Printf("Difficulty increased! Modifier: %d\n", g.difficultyModifier)
	}

	// Update HUD with wave info
	g.hud.EnemiesInWave = g.enemiesPerWave
	g.hud.EnemiesKilledInWave = 0

	fmt.Printf("Wave %d started! (%d enemies)\n", g.hud.CurrentWave, g.enemiesPerWave)
}

func (g *Game) placeTower(x, y float32) {
	// Check if clicking in HUD area
	if y < config.HUDHeight {
		g.errorMessage = "Cannot place tower in HUD area!"
		g.errorTimer = 120
		return
	}

	// Check if there's already a tower at this position
	for _, tower := range g.towers {
		dx := x - tower.PositionX
		dy := y - tower.PositionY
		distance := dx*dx + dy*dy
		minDistance := config.TowerSize * config.TowerSize

		if distance < minDistance {
			g.errorMessage = "Cannot place tower on another tower!"
			g.errorTimer = 120
			return
		}
	}

	// Check if tower limit reached
	if len(g.towers) >= g.towerLimit {
		g.errorMessage = "Tower limit reached! Buy more slots in the shop."
		g.errorTimer = 120
		return
	}

	// Check if has enough coins
	if g.coins < g.hud.TowerCost {
		g.errorMessage = "Not enough coins to place tower!"
		g.errorTimer = 120
		return
	}

	// Validate tower placement (not on path)
	if !entity.CanPlaceTower(x, y, g.maps[0]) {
		g.errorMessage = "Cannot place tower on path!"
		g.errorTimer = 120
		return
	}

	// Deduct coins and place tower
	g.coins -= g.hud.TowerCost
	g.hud.Coins = g.coins
	fmt.Printf("Tower placed at (%.1f, %.1f)! Coins left: %d\n", x, y, g.coins)
	g.towers = append(g.towers, entity.NewTower(x, y))
	g.hud.TowersBuilt = len(g.towers)
}

func (g *Game) removeTower(x, y float32) {
	if y < config.HUDHeight {
		return
	}

	for i, tower := range g.towers {
		halfSize := config.TowerSize / 2
		if x >= tower.PositionX-halfSize && x <= tower.PositionX+halfSize &&
			y >= tower.PositionY-halfSize && y <= tower.PositionY+halfSize {

			// Refund half the tower cost
			g.coins += g.towerRefund
			g.hud.Coins = g.coins

			// Remove tower
			g.towers[i] = g.towers[len(g.towers)-1]
			g.towers = g.towers[:len(g.towers)-1]
			g.hud.TowersBuilt = len(g.towers)
			fmt.Printf("Tower removed! Remaining: %d\n", len(g.towers))
			return
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy()), color.Black, false)

	renderer.DrawBuildableAreas(screen, g.maps[0])

	g.maps[0].Draw(screen)

	renderer.DrawEnemies(screen, g.enemies)

	renderer.DrawTowers(screen, g.towers)

	renderer.DrawProjectiles(screen, g.projectiles)

	g.hud.Draw(screen)

	g.shop.Draw(screen, g.coins, utils.DrawLargeText)

	g.gameOverScreen.Draw(screen, g.enemiesDefeated, utils.DrawLargeText)

	g.instructionsScreen.Draw(screen, utils.DrawLargeText)

	// Draw error message (below HUD, larger text)
	if g.errorMessage != "" {
		utils.DrawLargeText(screen, g.errorMessage, 20, float64(config.HUDHeight)+10, 1.5)
	}
}

// Layout defines the game's logical screen size (required by ebiten.Game interface)
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.Config.Width, config.Config.Height
}

// drawShop draws the shop overlay
func (g *Game) restartGame() {
	fmt.Println("Restarting game...")
	g.enemies = []*entity.Enemy{}
	g.towers = []entity.Tower{}
	g.projectiles = []entity.Projectile{}
	g.enemiesDefeated = 0
	g.lives = config.GameConstants.InitialLives
	g.coins = config.GameConstants.InitialCoins
	g.towerLimit = config.GameConstants.TowerLimit
	g.towerDamageBoost = config.GameConstants.TowerDamageBoost
	g.towerFireRateBoost = config.GameConstants.TowerFireRateBoost
	g.difficultyModifier = 1
	g.spawnInterval = config.GameConstants.SpawnInterval
	g.tick = 0
	g.enemiesPerWave = 0
	g.enemiesSpawnedInWave = 0
	g.lastSpawnTick = 0
	g.errorMessage = ""
	g.errorTimer = 0

	// Reset shop and game over screen
	g.shop.Close()
	g.shop = shop.NewShop()
	g.gameOverScreen.Reset()
	g.instructionsScreen.Hide()

	// Reset HUD
	g.hud = hud.NewHUD(g.towerLimit, config.GameConstants.InitialTowerCost, config.GameConstants.InitialTowerRefund, config.GameConstants.InitialLives, config.GameConstants.InitialCoins)
}

func (g *Game) handleShopClick(mx, my int) {
	itemID, purchased := g.shop.HandleClick(mx, my, g.coins)

	if !purchased {
		return
	}

	// Process the purchase
	newCoins, newTowerLimit, newDamageBoost, newFireRateBoost, success := g.shop.PurchaseItem(
		itemID, g.coins, g.towerLimit, g.towerDamageBoost, g.towerFireRateBoost,
	)

	if success {
		g.coins = newCoins
		g.towerLimit = newTowerLimit
		g.towerDamageBoost = newDamageBoost
		g.towerFireRateBoost = newFireRateBoost
		g.towerCost = 5 * newTowerLimit
		g.towerRefund = 3 * newTowerLimit

		// Update HUD
		g.hud.Coins = g.coins
		g.hud.TowersLimit = g.towerLimit
		g.hud.TowerCost = g.towerCost
		g.hud.TowerRefund = g.towerRefund

		// Log purchase
		switch itemID {
		case 1:
			fmt.Printf("Bought tower slot! New limit: %d\n", g.towerLimit)
		case 2:
			fmt.Printf("Tower damage increased! Total bonus: +%d\n", g.towerDamageBoost)
		case 4:
			fmt.Printf("Tower fire rate increased! Multiplier: %.1fx\n", g.towerFireRateBoost)
		}
	}
}
