package initializers

import (
	"github.com/calleros/go-jwt/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoadDbTables initializes the database tables and creates the first user with the "root" profile
func LoadDbTables(db *gorm.DB) error {
	// Check if the "root" profile already exists
	var rootProfile models.Profile
	result := db.Where("name = ?", "root").First(&rootProfile)
	if result.RowsAffected == 0 {
		// If "root" profile doesn't exist, create it
		rootProfile = models.Profile{Name: "Root", Description: "Usuario root"}
		db.Create(&rootProfile)
	}

	// Check if the "admin" profile already exists
	var adminProfile models.Profile
	result = db.Where("name = ?", "admin").First(&adminProfile)
	if result.RowsAffected == 0 {
		// If "admin" profile doesn't exist, create it
		adminProfile = models.Profile{Name: "Admin", Description: "Administrador del sistema"}
		db.Create(&adminProfile)
	}

	// Check if the "access" profile already exists
	var accessProfile models.Profile
	result = db.Where("name = ?", "access").First(&accessProfile)
	if result.RowsAffected == 0 {
		// If "access" profile doesn't exist, create it
		accessProfile = models.Profile{Name: "Acceso", Description: "Acceso al sistema"}
		db.Create(&accessProfile)
	}

	// Check if the "createUser" profile already exists
	var createUserProfile models.Profile
	result = db.Where("name = ?", "createUser").First(&createUserProfile)
	if result.RowsAffected == 0 {
		// If "createUser" profile doesn't exist, create it
		createUserProfile = models.Profile{Name: "Usuarios-Registro", Description: "Acceso a registro de nuevos usuarios"}
		db.Create(&createUserProfile)
	}

	// Check if the "userlist" profile already exists
	var userlistProfile models.Profile
	result = db.Where("name = ?", "userlist").First(&userlistProfile)
	if result.RowsAffected == 0 {
		// If "userlist" profile doesn't exist, create it
		userlistProfile = models.Profile{Name: "Usuarios-Consulta", Description: "Acceso a consulta de usuarios"}
		db.Create(&userlistProfile)
	}

	// Check if the "editUser" profile already exists
	var editUserProfile models.Profile
	result = db.Where("name = ?", "editUser").First(&editUserProfile)
	if result.RowsAffected == 0 {
		// If "editUser" profile doesn't exist, create it
		editUserProfile = models.Profile{Name: "Usuarios-Editar", Description: "Acceso a edición de usuarios"}
		db.Create(&editUserProfile)
	}

	// Check if the "deleteUser" profile already exists
	var deleteUserProfile models.Profile
	result = db.Where("name = ?", "deleteUser").First(&deleteUserProfile)
	if result.RowsAffected == 0 {
		// If "deleteUser" profile doesn't exist, create it
		deleteUserProfile = models.Profile{Name: "Usuarios-Eliminar", Description: "Acceso a eliminar usuarios"}
		db.Create(&deleteUserProfile)
	}

	// Check if the "createProfile" profile already exists
	var createProfileProfile models.Profile
	result = db.Where("name = ?", "createProfile").First(&createProfileProfile)
	if result.RowsAffected == 0 {
		// If "createProfile" profile doesn't exist, create it
		createProfileProfile = models.Profile{Name: "Perfiles-Registro", Description: "Acceso a registro de nuevos perfiles"}
		db.Create(&createProfileProfile)
	}

	// Check if the "profilelist" profile already exists
	var profilelistProfile models.Profile
	result = db.Where("name = ?", "profilelist").First(&profilelistProfile)
	if result.RowsAffected == 0 {
		// If "profilelist" profile doesn't exist, create it
		profilelistProfile = models.Profile{Name: "Perfiles-Consulta", Description: "Acceso a consulta de perfiles"}
		db.Create(&profilelistProfile)
	}

	// Check if the "editProfile" profile already exists
	var editProfileProfile models.Profile
	result = db.Where("name = ?", "editProfile").First(&editProfileProfile)
	if result.RowsAffected == 0 {
		// If "editProfile" profile doesn't exist, create it
		editProfileProfile = models.Profile{Name: "Perfiles-Editar", Description: "Acceso a edición de perfiles"}
		db.Create(&editProfileProfile)
	}

	// Check if there are any users in the database
	var users []models.User
	db.Find(&users)
	if len(users) == 0 {
		// If no users exist, create the first user and assign the peofiles to it
		password := "1234" // You may want to change this to a secure password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Login:    "root",
			Password: string(hashedPassword),
			Name:     "Marco Antonio Calleros Lozano",
			Email:    "marco.calleros@gmail.com",
			Profiles: []models.Profile{rootProfile}, // Pass the profile to the user
		}
		db.Create(&user)
	}

	return nil
}
