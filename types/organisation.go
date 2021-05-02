package types

import (
	"time"
)

type Organisation struct {
	Id          string    `json:"id"`
	App         string    `json:"app"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Teams       []string  `json:"teams"`
	ImageUrl    string    `json:"image_url"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
	Deleted     bool      `json:"deleted"`
	DeletedBy   string    `json:"deleted_by"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type OrganisationAlreadyExists struct {
	Name    string `json:"name"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type OrganisationNotFound struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type OrganisationOperation struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	State   string `json:"state"`
	Error   bool   `json:"error"`
}

type DeletedOrganisationsCount struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Count   int    `json:"count"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}
