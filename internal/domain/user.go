package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is an account stored in the database.
type User struct {
	ID          string `gorm:"type:uuid;primaryKey;column:id;" json:"id"`
	Username    string `gorm:"unique;not null" json:"username"`
	Password    string `gorm:"not null" json:"password"`
	DisplayName string `json:"display_name"`
	Email       string `gorm:"unique;not null" json:"email"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// BeforeCreate is a GORM hook that assigns a new UUID as the ID before the user is inserted.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	return nil
}
