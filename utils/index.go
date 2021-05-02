package utils

import (
	"app-auth/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/markbates/goth"
	"github.com/mongodb/mongo-go-driver/bson"
	"golang.org/x/crypto/bcrypt"
)


func UserBsonM(user types.UserData) bson.M {
	return bson.M{
		"id": user.Id,
		"username": user.Username,
		"age": user.Age,
		"firstname": user.Firstname,
		"lastname": user.Lastname,
		"email": user.Email,
		"verified": user.Verified,
		"password": user.Password,
		"provider": user.Provider,
		"picture": user.Picture,
	}
}

func CreateHttpCookie(value string, expires time.Time, secure bool) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "user_session"
	cookie.Expires = expires
	cookie.Value = value
	cookie.Path = ""
	cookie.HttpOnly = secure
	cookie.Secure = secure

	return cookie
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14); if err != nil {
		fmt.Println(err)
	}

	return string(bytes)
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateUserFromGoogleAuth(data goth.User) types.UserData {
	return types.UserData{
		Id: data.UserID,
		Password: "",
		Email: data.Email,
		Username: data.Name,
		Firstname: data.FirstName,
		Lastname: data.LastName,
		Picture: data.AvatarURL,
		Age: 0,
		Verified: true,
		Active: true,
		Provider: data.Provider,
		LastLoggedInScope: "",
	}
}

func Stringify(data interface{}) string {
	res, _ := json.Marshal(data)
	return string(res)
}

func ReadFile(path string) string {
	data, err := ioutil.ReadFile(path); if err != nil {
		return ""
	}

	return string(data)
}

func TransformClaimsToOrganisationMember(claims *types.UserDataClaims) types.UserData {
	return types.UserData{
		Id: claims.Id,
		Email: claims.Email,
		Lastname: claims.Lastname,
		Firstname: claims.Firstname,
		Username: claims.Username,
		Age: claims.Age,
		Verified: claims.Verified,
		Provider: claims.Provider,
		Active: claims.Active,
	}
}
