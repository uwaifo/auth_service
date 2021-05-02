package types

import (
	"time"
)

type Team struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Id               string    `json:"id"`
	ImageUrl         string    `json:"image_url"`
	OrganisationName string    `json:"organisation_name"`
	OrganisationId   string    `json:"organisation_id"`
	Members          []string  `json:"members"` // this would correspond to the user_ids cos we have users on the same system
	App              string    `json:"app"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdateAt         time.Time `json:"update_at"`
	Deleted          bool      `json:"deleted"`
	DeletedBy        string    `json:"deleted_by"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type TeamAlreadyExists struct {
	Name    string `json:"name"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type TeamNotFound struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type TeamOrganisationNotFound struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type TeamOperation struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	State   string `json:"state"`
	Error   bool   `json:"error"`
}

type TeamsCount struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Count   int    `json:"count"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}
