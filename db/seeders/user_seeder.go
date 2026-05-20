// internal/database/seeders/user_seeder.go
package seeders

import (
	"errors"

	"github.com/mustafazeren/go-ecommerce-course/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func (s UserSeeder) Seed(db *gorm.DB) error {
	var existingUser models.User

	err := db.Where("email = ?", "admin@zdev.com").
		First(&existingUser).Error

	if err == nil {
		// user zaten var
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("supersecret"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	admin := models.User{
		Email:     "admin@zdev.com",
		Password:  string(hashedPassword),
		FirstName: "System",
		LastName:  "Admin",
		Role:      models.UserRoleAdmin,
		IsActive:  true,
	}

	return db.Create(&admin).Error
}
