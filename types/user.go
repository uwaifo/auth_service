package types

import "github.com/dgrijalva/jwt-go"

type UserData struct {
	Email             string `json:"email"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Firstname         string `json:"firstname"`
	Lastname          string `json:"lastname"`
	Picture           string `json:"picture"`
	Id                string `json:"id"`
	Age               int    `json:"age"`
	Verified          bool   `json:"verified"`
	Provider          string `json:"provider"`
	Active            bool   `json:"active"`
	LastLoggedInScope string `json:"lastloggedinscope"`
}

type UserDataClaims struct {
	Id                string `json:"id"`
	Username          string `json:"username"`
	Age               int    `json:"age"`
	Firstname         string `json:"firstname"`
	Lastname          string `json:"lastname"`
	Email             string `json:"email"`
	Verified          bool   `json:"verified"`
	Provider          string `json:"provider"`
	Active            bool   `json:"active"`
	App               string `json:"app"`
	Expires           string `json:"expires"`
	LastLoggedInScope string `json:"team"`
	Scope 			  *UserMemberScope `json:"scope"`
	jwt.StandardClaims
}

type UserStats struct {
	Organisations   int    `json:"organisations"`
	Teams           int    `json:"teams"`
	ReceivedInvites int    `json:"received_invites"`
	SentInvites     int    `json:"sent_invites"`
	Username        string `json:"username"`
	Verified        bool   `json:"verified"`
	Active          bool   `json:"active"`
}
