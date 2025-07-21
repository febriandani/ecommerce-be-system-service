package utils

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// Middleware JWT function
func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(secret),
	})
}

// func ExtractClaims(c *fiber.Ctx) (proto.CredentialData, error) {
// 	authHeader := c.Get("Authorization")
// 	if authHeader == "" {
// 		return proto.CredentialData{}, fmt.Errorf("Authorization header is missing")
// 	}

// 	const bearerPrefix = "Bearer "
// 	if !strings.HasPrefix(authHeader, bearerPrefix) {
// 		return proto.CredentialData{}, fmt.Errorf("Invalid token format")
// 	}

// 	// Extract the token
// 	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

// 	// Parse the JWT token
// 	claims := jtoken.MapClaims{}
// 	token, err := jtoken.ParseWithClaims(tokenString, claims, func(token *jtoken.Token) (interface{}, error) {
// 		// Replace with your secret key used to sign the token
// 		return []byte("session"), nil
// 	})
// 	if err != nil || !token.Valid {
// 		return proto.CredentialData{}, fmt.Errorf("Invalid token")
// 	}

// 	fmt.Println("DEBUG : ", claims["user_id"])
// 	// Extract claims
// 	// Convert "user_id" from float64 to int64
// 	userId, ok := claims["user_id"].(float64)
// 	if !ok {
// 		return proto.CredentialData{}, fmt.Errorf("Invalid user_id format")
// 	}

// 	credential := proto.CredentialData{
// 		Id:       int64(userId), // Convert float64 to int64
// 		Email:    claims["email"].(string),
// 		Fullname: claims["fullname"].(string),
// 	}

// 	return credential, nil
// }
