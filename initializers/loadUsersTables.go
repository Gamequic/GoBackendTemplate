package initializers

import (
	"github.com/calleros/sich/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoadUsersTables initializes the database tables and creates the first user with the "root" profile
func LoadUsersTables(db *gorm.DB) error {
	// Define los perfiles iniciales que deseas crear
	initialProfiles := []models.Profile{
		{Description: "root", PrivAccess: true, PrivInsert: true, PrivUpdate: true, PrivDelete: true},
		{Description: "Admin", PrivAccess: true, PrivInsert: true, PrivUpdate: true, PrivDelete: true},
		{Description: "Acceso", PrivAccess: true},
		{Description: "Usuarios-Leer", PrivAccess: true, PrivInsert: true},
		{Description: "Usuarios-Leer-Crear", PrivAccess: true, PrivInsert: true},
		{Description: "Usuarios-Leer-Crear-Actualizar", PrivAccess: true, PrivInsert: true, PrivUpdate: true},
		{Description: "Usuarios-Leer-Crear-Actualizar-Borrar", PrivAccess: true, PrivInsert: true, PrivUpdate: true, PrivDelete: true},
	}

	// Crea los perfiles si no existen
	for _, profile := range initialProfiles {
		result := db.Where("description = ?", profile.Description).First(&profile)
		if result.RowsAffected == 0 {
			db.Create(&profile)
		}
	}

	// Verifica si hay usuarios en la base de datos
	var users []models.User
	db.Find(&users)
	if len(users) == 0 {
		// Si no hay usuarios, crea el primer usuario y asígnale el primer perfil
		password := "1234" // Cambia esto por una contraseña segura
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Obtén el primer perfil
		var firstProfile models.Profile
		db.First(&firstProfile)

		// Crea el usuario y asígnale el primer perfil
		user := models.User{
			Login:    "root",
			Password: string(hashedPassword),
			Name:     "Marco Antonio Calleros Lozano",
			Email:    "marco.calleros@gmail.com",
			Profiles: []models.Profile{firstProfile}, // Asigna el primer perfil al usuario
		}
		db.Create(&user)
	}

	return nil
}
