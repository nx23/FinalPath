# 🎮 Final Path

A classic tower defense game built with Go and Ebiten game engine.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Ebiten](https://img.shields.io/badge/Ebiten-v2-FF6B6B?style=flat)
![License](https://img.shields.io/badge/license-MIT-green)

## 📝 About

**Final Path** is a tower defense game where you strategically place towers to defend against waves of enemies. Earn coins by defeating enemies and use them to purchase tower upgrades and unlock more tower slots.

### Features

- 🗺️ **Wave-based Gameplay**: Face increasingly difficult waves of enemies
- 🏰 **Strategic Tower Placement**: Place towers in optimal positions to defend your path
- 💰 **Economy System**: Earn coins by defeating enemies
- 🛒 **Upgrade Shop**: Purchase damage boosts, fire rate improvements, and additional tower slots
- ❤️ **Lives System**: Lose lives when enemies reach the end of the path
- 🎯 **Smart Targeting**: Towers automatically target enemies within range
- 📊 **HUD Dashboard**: Track your coins, lives, wave number, and tower count

## 📁 Project Structure

```
FinalPath/
├── main.go                      # Entry point
├── internal/
│   ├── config/
│   │   └── constants.go         # Game constants and configuration
│   ├── entity/
│   │   ├── enemy.go             # Enemy logic and behavior
│   │   ├── tower.go             # Tower logic and targeting
│   │   └── projectile.go        # Projectile physics
│   ├── game/
│   │   └── game.go              # Core game loop and state
│   ├── gamemap/
│   │   └── map.go               # Map and path system
│   ├── gameover/
│   │   └── gameover.go          # Game over screen
│   ├── hud/
│   │   └── hud.go               # Heads-up display
│   ├── instructions/
│   │   └── instructions.go      # Tutorial screen
│   ├── renderer/
│   │   └── renderer.go          # Rendering functions
│   ├── shop/
│   │   └── shop.go              # Shop system
│   └── utils/
│       └── utils.go             # Utility functions
├── go.mod                       # Go dependencies
└── README.md                    # This file
```

## 🎯 How to Play

### Objective
Prevent enemies from reaching the end of the path by strategically placing defensive towers.

### Controls
- **Left Click**: Place tower (15 coins) or interact with shop/buttons
- **Right Click**: Remove tower (refunds 10 coins)
- **Mouse**: Navigate menus and UI

### Game Mechanics
- **Starting Resources**: 10 lives, 50 coins
- **Tower Placement**: Place towers on green buildable areas (costs 15 coins per tower)
- **Tower Removal**: Right-click removes towers and refunds 10 coins
- **Earning Coins**: +5 coins per enemy defeated
- **Wave System**: Each wave spawns more enemies than the previous
- **Lives**: Lose 1 life per enemy that reaches the end

### Shop Items
1. **Tower Slot** (100 coins) - Unlock an additional tower slot
2. **Damage Upgrade** (25 coins) - Increase all towers' damage by +5
3. **Fire Rate Upgrade** (20 coins) - Increase all towers' fire rate by 10%

## 🔧 Setup and Installation

### Prerequisites
- Go 1.21 or higher
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/nx23/FinalPath.git
cd FinalPath

# Download dependencies
go mod download

# Build the game
go build -o finalpath

# Run the game
./finalpath
```

### Development Mode

```bash
# Run directly without building
go run .

# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air
```

## 🏗️ Architecture

The project follows a clean, modular architecture with clear separation of concerns:

- **Entity Layer**: Game objects (enemies, towers, projectiles) with their own behavior
- **Game Layer**: Core game loop, state management, and coordination
- **UI Layer**: HUD, shop, instructions, and game over screens
- **Rendering Layer**: Centralized drawing functions for all visual elements
- **Map Layer**: Path definitions and collision detection
- **Config Layer**: Constants and configuration values

### Design Principles
- **Single Responsibility**: Each module has a focused purpose
- **Encapsulation**: Internal logic hidden behind clean interfaces
- **Separation of Concerns**: Game logic separated from rendering and UI
- **No Circular Dependencies**: Clean dependency graph

## 🎮 Game Stats

- **Window Size**: 800x720 pixels
- **HUD Height**: 120 pixels
- **Economy**:
  - Tower Cost: 15 coins
  - Tower Refund: 10 coins
  - Enemy Reward: 5 coins
  - Starting Coins: 50
- **Tower Stats**:
  - Base Damage: 10
  - Base Fire Rate: 1.0 shots/second
  - Range: 100 pixels
- **Enemy Stats**:
  - Base Health: 10 HP (scales with wave: 10 + (1 + (wave-1)*2) + (20 * difficulty))
  - Base Speed: 2.0 (scales with wave: 2 * (1 + (wave-1)*0.1))
  - Size: 25x25 pixels
- **Wave Scaling**: Base 3 enemies + 2 per wave number
- **Shop Prices**:
  - Tower Slot: 100 coins
  - Damage +5: 25 coins
  - Fire Rate +10%: 20 coins

## �📝 Code Conventions

- **Exported** functions and types start with uppercase (e.g., `NewEnemy`)
- **Private** functions and fields start with lowercase (internal use only)
- Packages in `internal/` cannot be imported by external code
- Methods follow Go patterns (e.g., `enemy.IsAlive()`)
- All X/Y positions in entities represent the **center** of the entity
- Constants defined in `internal/config/constants.go`
- Render functions separated into `internal/renderer/`

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- Built with [Ebiten](https://ebiten.org/) - A dead simple 2D game library for Go
- Inspired by classic tower defense games

## 📞 Contact

- GitHub: [@nx23](https://github.com/nx23)
- Project Link: [https://github.com/nx23/FinalPath](https://github.com/nx23/FinalPath)

---

**Made with ❤️ and Go**
