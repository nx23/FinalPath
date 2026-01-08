package config

// Entity sizes in pixels
const (
	EnemySize      float32 = 25
	TowerSize      float32 = 25
	PathWidth      float32 = 50
	ProjectileSize float32 = 5
)

// HUD configuration
const (
	HUDHeight   float32 = 120
	HUDFontSize float32 = 2.5
	MapOffsetY  float32 = 120
)

// Window holds the game window configuration
type Window struct {
	Width  int
	Height int
	Title  string
}

// Config is the default window configuration
var Config = Window{
	Width:  800,
	Height: 720,
	Title:  "Final Path v1.0",
}
