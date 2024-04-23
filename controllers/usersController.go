package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/calleros/sich/initializers"
	"github.com/calleros/sich/middleware" // Importar el paquete middleware
	"github.com/calleros/sich/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *gin.Context) {
	// Get te login/password off req body
	var body struct {
		Login    string
		Password string
		Name     string
		Email    string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the user
	user := models.User{Login: body.Login, Password: string(hash), Name: body.Name, Email: body.Email}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Find the profile with id 3
	var profile models.Profile
	result = initializers.DB.First(&profile, 3)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find profile with id 3",
		})
		return
	}

	// Associate the profile with the user
	initializers.DB.Model(&user).Association("Profiles").Append(&profile)

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	// Obtener el login y la contraseña del cuerpo de la solicitud
	var body struct {
		Login    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Buscar el usuario solicitado
	var user models.User
	err := initializers.DB.Preload("Profiles").First(&user, "login = ?", body.Login).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid login or password",
		})
		return
	}

	// Imprimir los perfiles obtenidos para depuración
	fmt.Println("Perfiles del usuario:", user.Profiles)

	// Comparar la contraseña enviada con el hash de la contraseña guardada del usuario
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid login or password",
		})
		return
	}

	// Construir un slice de perfiles del usuario con sus permisos
	var profiles []gin.H
	for _, profile := range user.Profiles {
		profiles = append(profiles, gin.H{
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

	// Generar los reclamos para el token JWT
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"profiles": profiles,
	}

	// Generar un token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar y obtener el token codificado completo como una cadena utilizando el secreto
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Enviar el token
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	// Obtener la cookie de la solicitud
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to get token from cookie",
		})
		return
	}

	// Parsear el token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Devolver la clave secreta almacenada en la variable de entorno "SECRET"
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to parse token",
		})
		return
	}

	// Verificar si el token es válido
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verificar si existen reclamos de perfiles en el token
		if profiles, ok := claims["profiles"]; ok {
			// Responder con los perfiles del usuario
			c.JSON(http.StatusOK, gin.H{"profiles": profiles})
			return
		}
		// Si no hay reclamos de perfiles en el token, responder con un error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No profiles found in token"})
		return
	}

	// Si el token no es válido, responder con un error
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
}

func Logout(c *gin.Context) {
	// Eliminar la cookie de autorización
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	// Responder con un mensaje de éxito
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func CreateUser(c *gin.Context) {
	// Verificar los permisos del usuario para crear un usuario
	err := middleware.VerifyAccess(c, []int{1, 2, 5, 6, 7}, map[string]bool{"PrivAccess": true, "PrivInsert": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener los datos del nuevo usuario del cuerpo de la solicitud
	var newUser struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Profiles []uint `json:"profiles" binding:"required"`
	}

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}

	// Crear el nuevo usuario
	user := models.User{
		Login:    newUser.Login,
		Password: string(hashedPassword),
		Name:     newUser.Name,
		Email:    newUser.Email,
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	// Asociar los perfiles con el usuario en la tabla user_profiles
	for _, profileID := range newUser.Profiles {
		var profile models.Profile
		result := initializers.DB.First(&profile, profileID)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find profile"})
			return
		}
		initializers.DB.Model(&user).Association("Profiles").Append(&profile)
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func UpdateUser(c *gin.Context) {
	// Verificar los permisos del usuario para actualizar un usuario
	err := middleware.VerifyAccess(c, []int{1, 2, 6, 7}, map[string]bool{"PrivAccess": true, "PrivUpdate": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del usuario a actualizar del parámetro de la ruta
	userID := c.Param("id")

	// Verificar si el ID del usuario es válido
	var existingUser models.User
	result := initializers.DB.First(&existingUser, userID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Obtener los datos actualizados del usuario del cuerpo de la solicitud
	var updatedUser struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Profiles []uint `json:"profiles"`
	}

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Actualizar los campos del usuario si se proporcionan
	if updatedUser.Login != "" {
		existingUser.Login = updatedUser.Login
	}
	if updatedUser.Password != "" {
		// Hash de la nueva contraseña
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
			return
		}
		existingUser.Password = string(hashedPassword)
	}
	if updatedUser.Name != "" {
		existingUser.Name = updatedUser.Name
	}
	if updatedUser.Email != "" {
		existingUser.Email = updatedUser.Email
	}

	// Actualizar los perfiles asociados con el usuario
	if len(updatedUser.Profiles) > 0 {
		// Limpiar los perfiles asociados existentes
		initializers.DB.Model(&existingUser).Association("Profiles").Clear()

		// Asociar los nuevos perfiles con el usuario
		for _, profileID := range updatedUser.Profiles {
			var profile models.Profile
			result := initializers.DB.First(&profile, profileID)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find profile"})
				return
			}
			initializers.DB.Model(&existingUser).Association("Profiles").Append(&profile)
		}
	}

	// Guardar los cambios en la base de datos
	result = initializers.DB.Save(&existingUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func GetUsers(c *gin.Context) {
	// Verificar los permisos del usuario para acceder a la lista de usuarios
	err := middleware.VerifyAccess(c, []int{1, 2, 4, 5, 6, 7}, map[string]bool{"PrivAccess": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener todos los usuarios de la base de datos, excepto el usuario con ID 1
	var users []models.User
	result := initializers.DB.Where("id != ?", 1).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get users"})
		return
	}

	// Crear un mapa para incluir el título "Users" y la lista de usuarios
	usersJSON := gin.H{
		"Users": users,
	}

	c.JSON(http.StatusOK, usersJSON)
}

func GetUserById(c *gin.Context) {
	// Verificar los permisos del usuario para acceder a un usuario específico
	err := middleware.VerifyAccess(c, []int{1, 2, 6, 7}, map[string]bool{"PrivAccess": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del usuario de la URL
	userID := c.Param("id")

	// Buscar el usuario en la base de datos por su ID
	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respuesta exitosa con el usuario encontrado
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	// Verificar los permisos del usuario para eliminar usuarios
	err := middleware.VerifyAccess(c, []int{1, 2, 7}, map[string]bool{"PrivAccess": true, "PrivDelete": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del usuario de los parámetros de la ruta
	userID := c.Param("id")

	// Eliminar los perfiles asociados al usuario en la tabla user_profiles
	if err := initializers.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Preload("Profiles").First(&user, userID).Error; err != nil {
			return err
		}
		// Eliminar los perfiles asociados al usuario
		if err := tx.Model(&user).Association("Profiles").Clear(); err != nil {
			return err
		}
		// Eliminar el usuario de la tabla users
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func DeleteUserCompletely(c *gin.Context) {
	// Verificar los permisos del usuario para eliminar usuarios
	err := middleware.VerifyAccess(c, []int{1}, map[string]bool{"PrivAccess": true, "PrivDelete": true})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del usuario de los parámetros de la ruta
	userID := c.Param("id")

	// Eliminar los perfiles asociados al usuario en la tabla user_profiles
	if err := initializers.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Preload("Profiles").First(&user, userID).Error; err != nil {
			return err
		}
		// Eliminar los perfiles asociados al usuario
		if err := tx.Model(&user).Association("Profiles").Clear(); err != nil {
			return err
		}
		// Eliminar completamente el usuario de la tabla users
		if err := tx.Unscoped().Delete(&user).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "User deleted completely"})
}
