package auth

import (
	"app-auth/types"
	"app-auth/utils"
	"crypto/rsa"

	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateRSA256SignedToken(key string, user types.UserData, scope *types.UserMemberScope, expires time.Time, app string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":        user.Id,
		"username":  user.Username,
		"age":       user.Age,
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
		"email":     user.Email,
		"verified":  user.Verified,
		"provider":  user.Provider,
		"active":    user.Active,
		"expires":   expires,
		"app":       app,
		"team":      user.LastLoggedInScope,
		"scope":    scope,
		"StandardClaims": jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
			Issuer:    "Clipsynphony",
			IssuedAt:  time.Now().Unix(),
		},
	})

	secret, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))

	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

func CreateRSA256RefreshToken(key string, user types.UserData, expires time.Time, app string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":        user.Id,
		"username":  user.Username,
		"age":       user.Age,
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
		"email":     user.Email,
		"verified":  user.Verified,
		"provider":  user.Provider,
		"expires":   expires,
		"app":       app,
		"active":    user.Active,
		"team":      user.LastLoggedInScope,
		"StandardClaims": jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
			Issuer:    "Clipsynphony",
			IssuedAt:  time.Now().Unix(),
		},
	})

	secret, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))

	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

func ParseRSA256SignedToken(tokenString string, key string) *types.UserDataClaims {

	secret, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))

	token, err := jwt.ParseWithClaims(tokenString, &types.UserDataClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return secret, nil
	})

	claims, ok := token.Claims.(*types.UserDataClaims)
	if ok && token.Valid {
		fmt.Printf("%s %v", claims.Username, claims.StandardClaims.ExpiresAt)
	} else {
		fmt.Println(err)
	}

	return claims
}

func GetSecret() *rsa.PublicKey {

	key := utils.ReadFile("./jwtRS256.key.pub")

	secret, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	if err != nil {
		return nil
	}

	return secret
}
