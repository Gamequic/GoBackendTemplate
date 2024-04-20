package middleware

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func VerifyAccess(c *gin.Context, requiredProfiles []int) error {
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
		// Verificar si el usuario tiene los perfiles necesarios
		if profiles, ok := claims["profiles"]; ok {
			// Convertir los perfiles a un slice de enteros
			var profilesClaim []interface{}
			profilesClaim, ok := profiles.([]interface{})
			if !ok {
				return fmt.Errorf("invalid profiles claim")
			}
			var profilesIDs []int
			for _, p := range profilesClaim {
				profile, ok := p.(float64)
				if !ok {
					return fmt.Errorf("invalid profile format")
				}
				profilesIDs = append(profilesIDs, int(profile))
			}

			// Verificar si el usuario tiene los perfiles necesarios
			hasRequiredProfile := false
			for _, requiredProfile := range requiredProfiles {
				for _, profileID := range profilesIDs {
					if profileID == requiredProfile {
						hasRequiredProfile = true
						break
					}
				}
			}

			if !hasRequiredProfile {
				return fmt.Errorf("unauthorized")
			}
		} else {
			return fmt.Errorf("no profiles found in token")
		}
	} else {
		return fmt.Errorf("invalid token")
	}

	return nil
}
