package initializers

import (
	"github.com/calleros/go-jwt/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoadDbTables initializes the database tables and creates the first user with the "admin" and "access" profile
func LoadDbTables(db *gorm.DB) error {
	// Check if the "admin" peofile already exists
	var adminProfile models.Profile
	result := db.Where("name = ?", "admin").First(&adminProfile)
	if result.RowsAffected == 0 {
		// If "admin" profile doesn't exist, create it
		adminProfile = models.Profile{Name: "admin", Description: "Administrador del sistema"}
		db.Create(&adminProfile)
	}

	// Check if the "access" profile already exists
	var accessProfile models.Profile
	result = db.Where("name = ?", "access").First(&accessProfile)
	if result.RowsAffected == 0 {
		// If "access" profile doesn't exist, create it
		accessProfile = models.Profile{Name: "acceso", Description: "Usuario con acceso al sistema"}
		db.Create(&accessProfile)
	}

	// Check if there are any users in the database
	var users []models.User
	db.Find(&users)
	if len(users) == 0 {
		// If no users exist, create the first user and assign the peofiles to it
		password := "92631043" // You may want to change this to a secure password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Login:    "admin",
			Password: string(hashedPassword),
			Name:     "Marco Antonio Calleros Lozano",
			Email:    "marco.calleros@gmail.com",
			Profiles: []models.Profile{adminProfile, accessProfile}, // Pass both profiles to the user
		}
		db.Create(&user)
	}

	return nil
}