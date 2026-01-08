package config

// Entity sizes in pixels
const (
	EnemySize      float32 = 25
	TowerSize      float32 = 25
	PathWidth      float32 = 50
	ProjectileSize float32 = 5
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
	Height: 600,
	Title:  "Final Path v1.0",
}
