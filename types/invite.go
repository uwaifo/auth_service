package types

type ConfirmationData struct {
	Email string `json:"email"`
	Confirmid string `json:"confirmid"`
	Userid string `json:"userid"`
}

type ConfirmObject struct{
	ConfirmationId string `json:"confirmation_id"`
	TeamId string `json:"team_id"`
	OrganisationId string `json:"organisation_id"`
}

