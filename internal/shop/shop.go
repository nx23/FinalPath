package shop

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Shop struct {
	Open   bool
	X      float32
	Y      float32
	Width  float32
	Height float32
	Items  []ShopItem
}

type ShopItem struct {
	ID          int
	Name        string
	Cost        int
	Description string
	Y           float32 // Y position relative to shop
}

func NewShop() *Shop {
	return &Shop{
		Open:   false,
		X:      200,
		Y:      200,
		Width:  400,
		Height: 320,
		Items: []ShopItem{
			{ID: 1, Name: "Buy Tower Slot", Cost: 30, Description: "Add +1 tower slot", Y: 100},
			{ID: 2, Name: "Tower Damage +10", Cost: 50, Description: "Increase all tower damage by +10", Y: 160},
			{ID: 4, Name: "Fire Rate +10%", Cost: 45, Description: "Increase all tower fire rate by 10%", Y: 220},
		},
	}
}

func (s *Shop) Draw(screen *ebiten.Image, coins int, drawTextFunc func(*ebiten.Image, string, float64, float64, float64)) {
	if !s.Open {
		return
	}

	// Semi-transparent overlay
	vector.FillRect(screen, 0, 0, float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy()),
		color.RGBA{0, 0, 0, 180}, false)

	// Shop panel
	vector.FillRect(screen, s.X, s.Y, s.Width, s.Height, color.RGBA{40, 40, 40, 255}, false)
	vector.StrokeRect(screen, s.X, s.Y, s.Width, s.Height, 3, color.RGBA{255, 165, 0, 255}, false)

	// Title
	drawTextFunc(screen, "SHOP", 360, 215, 3.0)

	// Coins display
	coinsText := fmt.Sprintf("Coins: %d", coins)
	drawTextFunc(screen, coinsText, 340, 260, 2.0)

	// Draw shop items
	for _, item := range s.Items {
		s.drawItem(screen, item, coins, drawTextFunc)
	}

	// Close instruction
	drawTextFunc(screen, "Right-click to close", 320, 600, 1.5)
}

// drawItem renders a single shop item
func (s *Shop) drawItem(screen *ebiten.Image, item ShopItem, coins int, drawTextFunc func(*ebiten.Image, string, float64, float64, float64)) {
	itemX := s.X + 20
	itemY := s.Y + item.Y
	itemWidth := float32(360)
	itemHeight := float32(50)

	canAfford := coins >= item.Cost

	// Background color based on affordability
	var bgColor color.RGBA
	if canAfford {
		bgColor = color.RGBA{0, 100, 0, 200} // Green if affordable
	} else {
		bgColor = color.RGBA{100, 0, 0, 200} // Red if not
	}

	vector.FillRect(screen, itemX, itemY, itemWidth, itemHeight, bgColor, false)
	vector.StrokeRect(screen, itemX, itemY, itemWidth, itemHeight, 2, color.RGBA{255, 255, 255, 255}, false)

	// Item text
	itemText := fmt.Sprintf("%s - %d coins", item.Name, item.Cost)
	drawTextFunc(screen, itemText, float64(itemX+10), float64(itemY+15), 1.8)
}

// HandleClick processes click events on shop items
// Returns the item ID if a purchase was made, or 0 if no purchase
func (s *Shop) HandleClick(mx, my int, coins int) (itemID int, purchased bool) {
	if !s.Open {
		return 0, false
	}

	for _, item := range s.Items {
		itemX := int(s.X + 20)
		itemY := int(s.Y + item.Y)
		itemWidth := 360
		itemHeight := 50

		// Check if click is within item bounds
		if mx >= itemX && mx <= itemX+itemWidth && my >= itemY && my <= itemY+itemHeight {
			if coins >= item.Cost {
				return item.ID, true
			}
			return item.ID, false // Clicked but can't afford
		}
	}

	return 0, false
}

// Toggle opens or closes the shop
func (s *Shop) Toggle() {
	s.Open = !s.Open
}

// Close closes the shop
func (s *Shop) Close() {
	s.Open = false
}

func (s *Shop) PurchaseItem(itemID, coins, towerLimit, towerDamageBoost int, towerFireRateBoost float32) (newCoins, newTowerLimit, newDamageBoost int, newFireRateBoost float32, success bool) {
	// Find the item
	var item *ShopItem
	for i := range s.Items {
		if s.Items[i].ID == itemID {
			item = &s.Items[i]
			break
		}
	}

	if item == nil || coins < item.Cost {
		return coins, towerLimit, towerDamageBoost, towerFireRateBoost, false
	}

	// Deduct cost
	newCoins = coins - item.Cost
	newTowerLimit = towerLimit
	newDamageBoost = towerDamageBoost
	newFireRateBoost = towerFireRateBoost

	// Apply effect based on item ID
	switch itemID {
	case 1: // Buy Tower Slot
		newTowerLimit++
	case 2: // Tower Damage Upgrade
		newDamageBoost += 10
	case 4: // Fire Rate Upgrade
		newFireRateBoost += 0.1
	default:
		return coins, towerLimit, towerDamageBoost, towerFireRateBoost, false
	}

	return newCoins, newTowerLimit, newDamageBoost, newFireRateBoost, true
}
