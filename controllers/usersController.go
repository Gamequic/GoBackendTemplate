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
func Validate(c *gin.Context) {
	//	user, _ := c.Get("user") //con esta linea regreso el usuario en un JSON

	c.JSON(http.StatusOK, gin.H{
		//		"message": user,
		"message": "I'm logged in",
	})
}

func GetUsersList(c *gin.Context) {
	// Obtener el usuario del contexto
	user, _ := c.Get("user")

	// Verificar si el usuario tiene acceso permitido
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			// Cargar los perfiles del usuario
			var userWithProfiles models.User
			if err := initializers.DB.Preload("Profiles").First(&userWithProfiles, userModel.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// Verificar si el usuario tiene perfil 1 o 3
			hasAccess := false
			for _, profile := range userWithProfiles.Profiles {
				if profile.ID == 1 || profile.ID == 3 {
					hasAccess = true
					break
				}
			}
			if !hasAccess {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
				return
			}
		}
	} else {
		// Si el usuario no está en el contexto, devolver un error de autorización
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Si el usuario tiene acceso, continuar obteniendo la lista de usuarios
	var users []models.User
	initializers.DB.Find(&users)

	// Responder con la lista de usuarios
	c.JSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	// Obtener el id del usuario de los parámetros de la ruta
	userId := c.Param("id")

	// Obtener el usuario del contexto
	authUser, _ := c.Get("user")

	// Verificar si el usuario tiene acceso permitido
	if authUser != nil {
		if userModel, ok := authUser.(models.User); ok {
			// Cargar los perfiles del usuario
			var userWithProfiles models.User
			if err := initializers.DB.Preload("Profiles").First(&userWithProfiles, userModel.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// Verificar si el usuario tiene perfil 1 o 6
			hasAccess := false
			for _, profile := range userWithProfiles.Profiles {
				if profile.ID == 1 || profile.ID == 6 {
					hasAccess = true
					break
				}
			}
			if !hasAccess {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
				return
			}
		}
	} else {
		// Si el usuario no está en el contexto, devolver un error de autorización
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verificar si el id es válido
	var user models.User
	err := initializers.DB.Preload("Profiles").First(&user, userId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Eliminar el usuario
	if err := initializers.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
