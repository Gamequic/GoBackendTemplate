package middleware

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func VerifyAccess(c *gin.Context, profileIDs []int, requiredPrivileges map[string]bool) error {
	// Verificar si se proporciona un token JWT en la cookie
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		return fmt.Errorf("failed to get token from cookie: %v", err)
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
		return fmt.Errorf("failed to parse token: %v", err)
	}

	// Verificar si el token es válido
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verificar si los perfiles y permisos son válidos
		if profiles, ok := claims["profiles"].([]interface{}); ok {
			for _, p := range profiles {
				profile, ok := p.(map[string]interface{})
				if !ok {
					return fmt.Errorf("invalid profile format")
				}

				// Verificar si el perfil está en los perfiles permitidos
				profileID := int(profile["ID"].(float64))
				if contains(profileIDs, profileID) {
					// Verificar si el perfil tiene los permisos necesarios
					for privilege, required := range requiredPrivileges {
						if profile[privilege] != required {
							return fmt.Errorf("unauthorized")
						}
					}
					// Si se encontró un perfil válido, salir del bucle
					return nil
				}
			}
		} else {
			return fmt.Errorf("no profiles found in token")
		}
	} else {
		return fmt.Errorf("invalid token")
	}

	return fmt.Errorf("unauthorized")
}

func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
