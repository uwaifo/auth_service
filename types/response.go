package types

import "time"

type InvalidRequest struct {
	Status  int
	Message string
}

type LoginResponse struct {
	Token   string    `json:"token"`
	UserId  string    `json:"user_id"`
	Refresh string    `json:"refresh"`
	Expires time.Time `json:"expires"`
	Message string    `json:"message"`
	Scope   string    `json:"scope"`
	Teams   []Team    `json:"teams"`
}

type SignupResponse struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Status  string    `json:"status"`
}

type LogoutResponse struct {
	Status  string
	Message string
}

type CalbackJSON struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type GenericResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Time    string `json:"time"`
}

type GenericResponseWithError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   error  `json:"error"`
	Time    string `json:"time"`
}

type EmailResetSuccess struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	Error         error  `json:"error"`
	Time          string `json:"time"`
	EmailChangeId string `json:"emailchangeid"`
}
