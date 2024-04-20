package controllers

import (
	"net/http"

	"github.com/calleros/sich/initializers"
	"github.com/calleros/sich/middleware" // Importar el paquete middleware
	"github.com/calleros/sich/models"
	"github.com/gin-gonic/gin"
)

func CreateProfile(c *gin.Context) {
	// Verificar acceso utilizando la función VerifyAccess del paquete middleware
	if err := middleware.VerifyAccess(c, []int{1, 2, 8}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener los datos del cuerpo de la solicitud
	var profileData models.Profile
	if err := c.BindJSON(&profileData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Crear el perfil
	result := initializers.DB.Create(&profileData)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create profile"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"profile": profileData})
}

func GetProfiles(c *gin.Context) {
	// Verificar acceso utilizando el token JWT almacenado en la cookie
	if err := middleware.VerifyAccess(c, []int{1, 2, 9}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Definir una estructura para la respuesta JSON
	type ProfileResponse struct {
		ID          uint   `json:"ID"`
		Name        string `json:"Name"`
		Description string `json:"Description"`
	}

	// Obtener todos los perfiles de la base de datos, excluyendo el perfil 1
	var profiles []models.Profile
	if err := initializers.DB.Where("id != ?", 1).Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profiles"})
		return
	}

	// Crear una slice de ProfileResponse y asignar los valores de los perfiles
	var profilesResponse []ProfileResponse
	for _, profile := range profiles {
		profilesResponse = append(profilesResponse, ProfileResponse{
			ID:          profile.ID,
			Name:        profile.Name,
			Description: profile.Description,
		})
	}

	// Respuesta exitosa con los perfiles
	c.JSON(http.StatusOK, gin.H{"profiles": profilesResponse})
}

func GetProfileById(c *gin.Context) {
	// Verificar acceso utilizando el token JWT almacenado en la cookie
	if err := middleware.VerifyAccess(c, []int{1, 2, 10}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener el ID del perfil de los parámetros de la ruta
	profileID := c.Param("id")

	// Definir una estructura para la respuesta JSON
	type ProfileResponse struct {
		ID          uint   `json:"ID"`
		Name        string `json:"Name"`
		Description string `json:"Description"`
	}

	// Buscar el perfil por ID en la base de datos
	var profile models.Profile
	if err := initializers.DB.Select("id, name, description").First(&profile, profileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Crear una instancia de ProfileResponse y asignar los valores del perfil
	profileResponse := ProfileResponse{
		ID:          profile.ID,
		Name:        profile.Name,
		Description: profile.Description,
	}

	// Respuesta exitosa con el perfil
	c.JSON(http.StatusOK, gin.H{"profile": profileResponse})
}

func UpdateProfile(c *gin.Context) {
	// Verificar acceso utilizando el token JWT almacenado en la cookie
	if err := middleware.VerifyAccess(c, []int{1, 2, 10}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener el ID del perfil de los parámetros de la ruta o del cuerpo de la solicitud
	profileID := c.Param("id")

	// Buscar el perfil en la base de datos
	var profile models.Profile
	if err := initializers.DB.First(&profile, profileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Bind los datos del cuerpo de la solicitud al perfil existente
	if err := c.Bind(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process request body"})
		return
	}

	// Guardar los cambios en la base de datos
	if err := initializers.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetUsersByProfileId(c *gin.Context) {
	// Verificar acceso utilizando la función VerifyAccess del paquete middleware
	if err := middleware.VerifyAccess(c, []int{1, 2, 9}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener el ID del perfil de los parámetros de la ruta
	profileID := c.Param("id")

	// Buscar el perfil en la base de datos
	var profile models.Profile
	if err := initializers.DB.Preload("Users").First(&profile, profileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Definir una estructura para la respuesta JSON
	type UserResponse struct {
		ID    uint   `json:"ID"`
		Login string `json:"Login"`
		Name  string `json:"Name"`
		Email string `json:"Email"`
	}

	// Obtener los usuarios asociados al perfil y crear una lista de respuesta
	var usersResponse []UserResponse
	for _, user := range profile.Users {
		usersResponse = append(usersResponse, UserResponse{
			ID:    user.ID,
			Login: user.Login,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	// Respuesta exitosa con los usuarios asociados al perfil
	c.JSON(http.StatusOK, gin.H{"users": usersResponse})
}
