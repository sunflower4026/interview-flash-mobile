package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"unique"`
	Pin         string
	Address     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(*gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

func (u *User) BeforeUpdate(*gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
