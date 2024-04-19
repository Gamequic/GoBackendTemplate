package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/calleros/go-jwt/initializers"
	"github.com/calleros/go-jwt/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	// Find the profile with id 2
	var profile models.Profile
	result = initializers.DB.First(&profile, 2)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find profile with id 2",
		})
		return
	}

	// Associate the profile with the user
	initializers.DB.Model(&user).Association("Profiles").Append(&profile)

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	// Get the login and pass from request body
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

	// Look up requested user
	var user models.User
	initializers.DB.Preload("Profiles").First(&user, "login = ?", body.Login)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid login or password",
		})
		return
	}

	// Compare sent password with saved user password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid login or password",
		})
		return
	}

	// Check if the user is active
	if !user.Active {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User is not active",
		})
		return
	}

	// Check if the user has access (profile id 1 or 2)
	var hasAccess bool
	for _, profile := range user.Profiles {
		if profile.ID == 1 || profile.ID == 2 {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User without access",
		})
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Send back the token
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

// validate sencillo solo regresa los datos del usuario sin los perfiles
//func Validate(c *gin.Context) {
//	user, _ := c.Get("user")

//	c.JSON(http.StatusOK, gin.H{
//		"message": user,
//	})
//}

// validate con los datos del usuario y perfiles
func Validate(c *gin.Context) {
	// Obtener el usuario del contexto
	user, _ := c.Get("user")

	// Verificar si el usuario tiene un ID válido
	if user != nil {
		// Convertir el usuario a un modelo.User
		userModel, ok := user.(models.User)
		if !ok {
			// Manejar el caso en el que el usuario no sea de tipo models.User
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user type",
			})
			return
		}

		// Cargar los perfiles del usuario desde la base de datos
		var userProfiles []models.Profile
		initializers.DB.Model(&userModel).Association("Profiles").Find(&userProfiles)

		// Crear una estructura para almacenar la respuesta JSON
		response := gin.H{
			"ID":       userModel.ID,
			"Login":    userModel.Login,
			"Name":     userModel.Name,
			"Email":    userModel.Email,
			"Active":   userModel.Active,
			"Profiles": userProfiles,
		}

		// Responder con la información del usuario y sus perfiles
		c.JSON(http.StatusOK, response)
	} else {
		// Manejar el caso en el que no se pueda obtener el usuario del contexto
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found in context",
		})
	}
}

func GetUsersList(c *gin.Context) {
	// Obtener todos los usuarios omitiendo el usuario con ID igual a 1
	var users []models.User
	initializers.DB.Where("id != ?", 1).Find(&users)

	// Responder con la lista de usuarios
	c.JSON(http.StatusOK, users)
}
