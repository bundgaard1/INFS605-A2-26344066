package domain

import "time"

type Notification struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"index;not null" json:"user_id"` // Hvem notifikationen tilhører
	Title     string    `gorm:"not null" json:"title"`         // f.eks. "Ny karakter modtaget"
	Message   string    `gorm:"not null" json:"message"`       // f.eks. "Du har fået 12 i INFS-605"
	LinkURL   string    `json:"link_url,omitempty"`            // f.eks. "/assignments/asg_9921" (så brugeren kan klikke)
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
