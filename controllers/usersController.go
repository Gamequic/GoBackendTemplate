package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/calleros/go-jwt/initializers"
	"github.com/calleros/go-jwt/middleware" // Importar el paquete middleware
	"github.com/calleros/go-jwt/models"
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
	initializers.DB.Preload("Profiles").First(&user, "login = ?", body.Login)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid login or password",
		})
		return
	}

	// Comparar la contraseña enviada con el hash de la contraseña guardada del usuario
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid login or password",
		})
		return
	}

	// Construir un slice de IDs de perfiles del usuario
	var profileIDs []uint
	for _, profile := range user.Profiles {
		profileIDs = append(profileIDs, profile.ID)
	}

	// Generar los reclamos para el token JWT
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"profiles": profileIDs,
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
	// Verificar acceso utilizando la función VerifyAccess del paquete middleware
	if err := middleware.VerifyAccess(c, []int{1, 2, 4}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener los datos del cuerpo de la solicitud
	var body struct {
		Login    string
		Password string
		Name     string
		Email    string
		Profiles []uint // Lista de IDs de perfiles asociados al usuario
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Hash de la contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Crear el usuario
	user := models.User{Login: body.Login, Password: string(hash), Name: body.Name, Email: body.Email}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Asociar los perfiles con el usuario en la tabla user_profiles
	for _, profileID := range body.Profiles {
		var profile models.Profile
		result := initializers.DB.First(&profile, profileID)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to find profile",
			})
			return
		}
		initializers.DB.Model(&user).Association("Profiles").Append(&profile)
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{})
}

func GetUsers(c *gin.Context) {
	// Verificar acceso utilizando la función VerifyAccess del paquete middleware
	if err := middleware.VerifyAccess(c, []int{1, 2, 5}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener todos los usuarios excluyendo al usuario con ID 1
	var users []models.User
	initializers.DB.Where("id != ?", 1).Find(&users)

	// Estructura para almacenar los datos del usuario en el orden deseado
	type UserResponse struct {
		ID    uint   `json:"ID"`
		Login string `json:"Login"`
		Name  string `json:"Name"`
		Email string `json:"Email"`
	}

	// Preparar los datos para la respuesta JSON
	var responseData []UserResponse
	for _, user := range users {
		userData := UserResponse{
			ID:    user.ID,
			Login: user.Login,
			Name:  user.Name,
			Email: user.Email,
		}
		responseData = append(responseData, userData)
	}

	// Crear un mapa para contener la matriz de usuarios
	responseMap := gin.H{"users": responseData}

	// Respuesta exitosa
	c.JSON(http.StatusOK, responseMap)
}

func GetUserById(c *gin.Context) {
	// Verificar acceso utilizando el token JWT almacenado en la cookie
	if err := middleware.VerifyAccess(c, []int{1, 3}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener el ID del usuario de los parámetros de la ruta
	userID := c.Param("id")

	// Definir una estructura para la respuesta JSON
	type UserProfile struct {
		ID          uint   `json:"ID"`
		Name        string `json:"Name"`
		Description string `json:"Description"`
	}

	type UserResponse struct {
		ID       uint          `json:"ID"`
		Login    string        `json:"Login"`
		Name     string        `json:"Name"`
		Email    string        `json:"Email"`
		Profiles []UserProfile `json:"Profiles"`
	}

	// Obtener el usuario de la base de datos
	var user models.User
	if err := initializers.DB.Preload("Profiles").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Crear una instancia de UserResponse y asignar los valores del usuario y los perfiles
	response := UserResponse{
		ID:    user.ID,
		Login: user.Login,
		Name:  user.Name,
		Email: user.Email,
		Profiles: func() []UserProfile {
			profiles := make([]UserProfile, len(user.Profiles))
			for i, profile := range user.Profiles {
				profiles[i] = UserProfile{
					ID:          profile.ID,
					Name:        profile.Name,
					Description: profile.Description,
				}
			}
			return profiles
		}(),
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"user": response})
}

func UpdateUser(c *gin.Context) {
	// Verificar acceso utilizando la función VerifyAccess del paquete middleware
	if err := middleware.VerifyAccess(c, []int{1, 2, 6}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener el ID del usuario de los parámetros de la ruta
	userID := c.Param("id")

	// Buscar el usuario por ID
	var user models.User
	err := initializers.DB.Preload("Profiles").First(&user, userID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Obtener los datos del cuerpo de la solicitud
	var body struct {
		Login    string
		Password string
		Name     string
		Email    string
		Profiles []uint // Lista de IDs de perfiles asociados al usuario
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Actualizar los campos del usuario
	if body.Login != "" {
		user.Login = body.Login
	}
	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		user.Password = string(hash)
	}
	if body.Name != "" {
		user.Name = body.Name
	}
	if body.Email != "" {
		user.Email = body.Email
	}

	// Borrar los registros de la tabla user_profiles para este usuario
	initializers.DB.Model(&user).Association("Profiles").Clear()

	// Asociar los perfiles con el usuario en la tabla user_profiles
	for _, profileID := range body.Profiles {
		var profile models.Profile
		result := initializers.DB.First(&profile, profileID)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to find profile",
			})
			return
		}
		initializers.DB.Model(&user).Association("Profiles").Append(&profile)
	}

	// Guardar los cambios en la base de datos
	result := initializers.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{})
}

func DeleteUser(c *gin.Context) {
	// Verificar acceso utilizando la función VerifyAccess del paquete middleware
	if err := middleware.VerifyAccess(c, []int{1, 2, 7}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
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
	// Verificar acceso utilizando la función VerifyAccess del paquete middleware
	if err := middleware.VerifyAccess(c, []int{1}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
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
