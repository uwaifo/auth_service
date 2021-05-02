package types

// user edit info payload
type UserEditPUTPayload struct {
	FirstName string `json:"firstName" form:"firstName" query:"firstName"`
	LastName  string `json:"lastName" form:"lastName" query:"lastName"`
	UserName  string `json:"userName" form:"userName" query:"userName"`
	Email     string `json:"email" form:"email" query:"email"`
	Picture   string `json:"picture" form:"picture" query:"picture"`
}
type UserUpdatePasswordPUTPayload struct {
	Password    string `json:"password" form:"password" query:"password"`
	OldPassword string `json:"oldPassword" form:"oldPassword" query:"oldPassword"`
}

// organisation actions payloads...
type OrganisationPOSTPayload struct {
	Name     string `json:"name" form:"name" query:"name"`
	ImageUrl string `json:"imageUrl" form:"imageUrl" query:"imageUrl"`
}
type OrganisationPUTPayload struct {
	Name     string `json:"name" form:"name" query:"name"`
	ImageUrl string `json:"imageUrl" form:"imageUrl" query:"imageUrl"`
}

// team actions payloads
type TeamPOSTPayload struct {
	Name     string `json:"name" form:"name" query:"name"`
	ImageUrl string `json:"imageUrl" form:"imageUrl" query:"imageUrl"`
}
type TeamPUTPayload struct {
	Name     string `json:"name" form:"name" query:"name"`
	ImageUrl string `json:"imageUrl" form:"imageUrl" query:"imageUrl"`
}
type TeamMemberInvitePOSTPayload struct {
	Email          string `json:"email" form:"email" query:"email"`
	SignupUrl      string `json:"signupUrl" form:"signupUrl" query:"signupUrl"`
	AppRedirectUrl string `json:"appRedirectUrl" form:"appRedirectUrl" query:"appRedirectUrl"`
}

// user actions payload...
type UserLoginPOSTPayload struct {
	App      string `json:"app" form:"app" query:"app"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}

// user actions payload...
type UserLoginSignupPOSTPayload struct {
	App      string `json:"app" form:"app" query:"app"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}
type UserSignUpPOSTPayload struct {
	App            string `json:"app" form:"app" query:"app"`
	Email          string `json:"email" form:"email" query:"email"`
	FirstName      string `json:"firstName" form:"firstName" query:"firstName"`
	LastName       string `json:"lastName" form:"lastName" query:"lastName"`
	Username       string `json:"username" form:"username" query:"username"`
	Password       string `json:"password" form:"password" query:"password"`
	Picture        string `json:"picture" form:"picture" query:"picture"`
	AppRedirectUrl string `json:"appRedirectUrl" form:"appRedirectUrl" query:"appRedirectUrl"`
}

// user actions payload...
type UserPasswordResetRequestPOSTPayload struct {
	App         string `json:"app" form:"app" query:"app"`
	RedirectUrl string `json:"redirectUrl" form:"redirectUrl" query:"redirectUrl"`
	Email       string `json:"email" form:"email" query:"email"`
}

// user actions payload...
type UserEmailChangeRequestPOSTPayload struct {
	App string `json:"app" form:"app" query:"app"`
	// you'll need this redirect url for redirection when the user clicks on the email change link
	RedirectUrl string `json:"redirectUrl" form:"redirectUrl" query:"redirectUrl"`
	// the new email the user needs to change
	Email string `json:"email" form:"email" query:"email"`
	// the user password
	Password string `json:"password" form:"password" query:"password"`
}

// user actions payload...
type UserPasswordResetChangePOSTPayload struct {
	App      string `json:"app" form:"app" query:"app"`
	Email    string `json:"email" form:"email" query:"email"`
	ResetId  string `json:"resetId" form:"resetId" query:"resetId"`
	Password string `json:"password" form:"password" query:"password"`
}

//
type UserEmailChangeRequestPostPayload struct {
	App         string `json:"app" form:"app" query:"app"`
	RedirectUrl string `json:"redirectUrl" form:"redirectUrl" query:"redirectUrl"`
	OldEmail    string `json:"oldemail" form:"oldemail" query:"oldemail"`
	NewEmail    string `json:"newemail" form:"newemail" query:"newemail"`
	Password    string `json:"password" form:"password" query:"password"`
}

// member scope actions payloads
type MemberScopePUTPayload struct {
	Scope string `json:"scope" form:"scope" query:"scope"`
}
type MemberScopeDELETEPayload struct {
	Scope string `json:"scope" form:"scope" query:"scope"`
}

// IAM Actions
type IAMScopePOSTPayload struct {
	Scope       string   `json:"scope" form:"scope" query:"scope"`
	Permissions []string `json:"permissions" form:"permissions" query:"permissions"`
}
type IAMScopePUTPayload struct {
	Scope string `json:"scope" form:"scope" query:"scope"`
}
type IAMScopeDELETEPermissionPayload struct {
	Permission string `json:"permission" form:"permission" query:"permission"`
}
type IAMScopePOSTPermissionPayload struct {
	Permission string `json:"permission" form:"permission" query:"permission"`
}
type IAMScopePOSTPermissionsPayload struct {
	Permissions []string `json:"permissions" form:"permissions" query:"permissions"`
}

// token payloads
type RefreshToken struct {
	Token string `json:"token" form:"token" query:"token"`
}
