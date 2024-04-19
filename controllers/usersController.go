package controllers

import (
	"net/http"
	"os"
	"strconv"
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

	// Check if the user is active (elimine el campo active ya que al eliminar al usuario
	// no se borra de la bd, solo le agrega fecha al campo deleted_at y eso lo inactiva)
	//if !user.Active {
	//	c.JSON(http.StatusUnauthorized, gin.H{
	//		"error": "User is not active",
	//	})
	//	return
	//}

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

func Validate(c *gin.Context) {
	//	user, _ := c.Get("user") //con esta linea regreso el usuario en un JSON

	c.JSON(http.StatusOK, gin.H{
		//		"message": user,
		"message": "I'm logged in",
	})
}

func CreateUser(c *gin.Context) {
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

			// Verificar si el usuario tiene perfil 1 o 4
			hasAccess := false
			for _, profile := range userWithProfiles.Profiles {
				if profile.ID == 1 || profile.ID == 4 {
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

	// Get te login/password off req body
	var body struct {
		Login    string
		Password string
		Name     string
		Email    string
		Profiles []int // Lista de IDs de perfiles
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

	// Asociar los perfiles con el usuario
	for _, profileID := range body.Profiles {
		var profile models.Profile
		if err := initializers.DB.First(&profile, profileID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to find profile with id " + strconv.Itoa(profileID),
			})
			return
		}
		user.Profiles = append(user.Profiles, profile)
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func GetUsers(c *gin.Context) {
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

	// Si el usuario tiene acceso, continuar obteniendo la lista de usuarios excluyendo al usuario con ID 1
	var users []models.User
	initializers.DB.Where("id != ?", 1).Find(&users)

	// Responder con la lista de usuarios
	c.JSON(http.StatusOK, users)
}

func GetUserById(c *gin.Context) {
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

			// Verificar si el usuario tiene perfil 1 o 5
			hasAccess := false
			for _, profile := range userWithProfiles.Profiles {
				if profile.ID == 1 || profile.ID == 5 {
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

	// Respond
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
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

			// Verificar si el usuario tiene perfil 1 o 5
			hasAccess := false
			for _, profile := range userWithProfiles.Profiles {
				if profile.ID == 1 || profile.ID == 5 {
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
	if err := initializers.DB.Preload("Profiles").First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parsear el cuerpo de la solicitud
	var updateUser struct {
		Login    string
		Password string
		Name     string
		Email    string
		Profiles []uint // Lista de IDs de perfiles
	}

	if c.Bind(&updateUser) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Actualizar los campos del usuario
	if updateUser.Login != "" {
		user.Login = updateUser.Login
	}
	if updateUser.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hash)
	}
	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	// Eliminar los perfiles existentes del usuario
	if err := initializers.DB.Model(&user).Association("Profiles").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear user profiles"})
		return
	}

	// Agregar los nuevos perfiles proporcionados en la solicitud
	var profiles []models.Profile
	if err := initializers.DB.Find(&profiles, updateUser.Profiles).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find profiles"})
		return
	}
	user.Profiles = profiles

	// Guardar los cambios en la base de datos
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
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
	if err := initializers.DB.Preload("Profiles").First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Eliminar los perfiles asociados al usuario
	if err := initializers.DB.Model(&user).Association("Profiles").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear user profiles"})
		return
	}

	// Eliminar el usuario
	if err := initializers.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func DeleteUserCompletely(c *gin.Context) {
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

			// Verificar si el usuario tiene perfil 1
			hasAccess := false
			for _, profile := range userWithProfiles.Profiles {
				if profile.ID == 1 {
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
	if err := initializers.DB.Preload("Profiles").First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Eliminar los perfiles asociados al usuario
	if err := initializers.DB.Model(&user).Association("Profiles").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear user profiles"})
		return
	}

	// Eliminar completamente el usuario
	if err := initializers.DB.Unscoped().Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted completely"})
}

func CreateProfile(c *gin.Context) {
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

	// Parsear el cuerpo de la solicitud
	var newProfile models.Profile
	if err := c.BindJSON(&newProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Crear el nuevo perfil
	result := initializers.DB.Create(&newProfile)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusOK, newProfile)
}

func GetProfiles(c *gin.Context) {
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

	// Obtener todos los perfiles
	var profiles []models.Profile
	initializers.DB.Find(&profiles)

	c.JSON(http.StatusOK, profiles)
}

func GetProfileById(c *gin.Context) {
	// Obtener el ID del perfil de los parámetros de la ruta
	profileId := c.Param("id")

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

	// Buscar el perfil por su ID
	var profile models.Profile
	if err := initializers.DB.First(&profile, profileId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func UpdateProfile(c *gin.Context) {
	// Obtener el ID del perfil de los parámetros de la ruta
	profileId := c.Param("id")

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

	// Parsear el cuerpo de la solicitud
	var updatedProfile models.Profile
	if err := c.BindJSON(&updatedProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Actualizar el perfil
	if err := initializers.DB.Model(&models.Profile{}).Where("id = ?", profileId).Updates(&updatedProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, updatedProfile)
}

func DeleteProfile(c *gin.Context) {
	// Obtener el ID del perfil de los parámetros de la ruta
	profileId := c.Param("id")

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

	// Eliminar el perfil
	if err := initializers.DB.Delete(&models.Profile{}, profileId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}

func DeleteProfileCompletely(c *gin.Context) {
	// Obtener el ID del perfil de los parámetros de la ruta
	profileId := c.Param("id")

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

	// Eliminar el perfil completamente
	if err := initializers.DB.Unscoped().Delete(&models.Profile{}, profileId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile completely"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile completely deleted successfully"})
}
