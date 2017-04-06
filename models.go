package sky

import "time"

// Model model
type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
