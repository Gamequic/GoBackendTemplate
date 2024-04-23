package controllers

import (
	"net/http"

	"github.com/calleros/sich/initializers"
	"github.com/calleros/sich/middleware"
	"github.com/calleros/sich/models"
	"github.com/gin-gonic/gin"
)

func CreateProfile(c *gin.Context) {
	// Verificar los permisos del usuario para crear perfiles
	err := middleware.VerifyAccess(c, []int{1}, map[string]bool{"PrivInsert": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener los datos del nuevo perfil del cuerpo de la solicitud
	var newProfile models.Profile
	if err := c.ShouldBindJSON(&newProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Crear el nuevo perfil
	result := initializers.DB.Create(&newProfile)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create profile"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Profile created successfully"})
}

func UpdateProfile(c *gin.Context) {
	// Verificar los permisos del usuario para actualizar perfiles
	err := middleware.VerifyAccess(c, []int{1}, map[string]bool{"PrivUpdate": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del perfil de los parámetros de la ruta
	profileID := c.Param("id")

	// Obtener el perfil existente de la base de datos
	var existingProfile models.Profile
	result := initializers.DB.First(&existingProfile, profileID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find profile"})
		return
	}

	// Decodificar los datos del cuerpo de la solicitud en el perfil existente
	if err := c.ShouldBindJSON(&existingProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Actualizar el perfil en la base de datos
	result = initializers.DB.Save(&existingProfile)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update profile"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetProfiles(c *gin.Context) {
	// Verificar los permisos del usuario para obtener perfiles
	err := middleware.VerifyAccess(c, []int{1}, map[string]bool{"PrivAccess": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener todos los perfiles excepto el perfil con ID 1
	var profiles []models.Profile
	result := initializers.DB.Not("id = ?", 1).Find(&profiles)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch profiles"})
		return
	}

	// Crear un slice para almacenar los perfiles en formato JSON
	var profilesJSON []gin.H
	for _, profile := range profiles {
		profilesJSON = append(profilesJSON, gin.H{
			"ID":          profile.ID,
			"Description": profile.Description,
			"PrivAccess":  profile.PrivAccess,
			"PrivExport":  profile.PrivExport,
			"PrivPrint":   profile.PrivPrint,
			"PrivInsert":  profile.PrivInsert,
			"PrivUpdate":  profile.PrivUpdate,
			"PrivDelete":  profile.PrivDelete,
		})
	}

	// Crear un mapa con la respuesta JSON
	response := gin.H{"profiles": profilesJSON}

	// Agregar el título "Profiles" al mapa de respuesta
	response["title"] = "Profiles"

	// Respuesta exitosa
	c.JSON(http.StatusOK, response)
}

func GetProfileById(c *gin.Context) {
	// Verificar los permisos del usuario para obtener perfiles
	err := middleware.VerifyAccess(c, []int{1}, map[string]bool{"PrivAccess": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del perfil de los parámetros de la ruta
	profileID := c.Param("id")

	// Buscar el perfil en la base de datos por su ID
	var profile models.Profile
	result := initializers.DB.First(&profile, profileID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find profile"})
		return
	}

	// Crear un mapa con la información del perfil en formato JSON
	profileJSON := gin.H{
		"ID":          profile.ID,
		"Description": profile.Description,
		"PrivAccess":  profile.PrivAccess,
		"PrivExport":  profile.PrivExport,
		"PrivPrint":   profile.PrivPrint,
		"PrivInsert":  profile.PrivInsert,
		"PrivUpdate":  profile.PrivUpdate,
		"PrivDelete":  profile.PrivDelete,
	}

	// Crear un mapa con la respuesta JSON
	response := gin.H{"profile": profileJSON}

	// Agregar el título "Profile" al mapa de respuesta
	response["title"] = "Profile"

	// Respuesta exitosa
	c.JSON(http.StatusOK, response)
}

func GetUsersByProfileId(c *gin.Context) {
	// Verificar los permisos del usuario para obtener usuarios por perfil
	err := middleware.VerifyAccess(c, []int{1}, map[string]bool{"PrivAccess": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del perfil de los parámetros de la ruta
	profileID := c.Param("id")

	// Buscar usuarios asociados al perfil en la tabla de usuarios
	var users []models.User
	result := initializers.DB.Joins("JOIN sec_users_profiles ON sec_users_profiles.user_id = sec_users.id").Where("sec_users_profiles.profile_id = ?", profileID).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find users"})
		return
	}

	// Crear un slice de mapas con la información de los usuarios en formato JSON
	var usersJSON []gin.H
	for _, user := range users {
		userJSON := gin.H{
			"ID":    user.ID,
			"Login": user.Login,
			"Name":  user.Name,
			"Email": user.Email,
		}
		usersJSON = append(usersJSON, userJSON)
	}

	// Crear un mapa con la respuesta JSON
	response := gin.H{"users": usersJSON}

	// Agregar el título "Users" al mapa de respuesta
	response["title"] = "Users"

	// Respuesta exitosa
	c.JSON(http.StatusOK, response)
}
