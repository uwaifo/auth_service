package mail

import (
	"fmt"
	"os"
	"strings"
)

func SignupTemplate(confirmationId string, redirectId string, app string) string {
	authUrl := os.Getenv("ENDPOINT")
	if authUrl == "" {
		authUrl = "http://app-auth"
	}
	confirmationUrl := authUrl + "/confirm/" + confirmationId + "/redirect/" + redirectId + "?app=" + app
	return strings.Replace(SignupTemplateHTML, "%s", confirmationUrl, -1)
}

func EmailResetTemplate(confirmationId string) string {

	return fmt.Sprintf(
		`
			<html> 
				<body style="background-color: white;padding: 20px;">
					<h1 style="text-align: center;"> Clipsynphony </h1>
					<button><a href="%s">Verify Email</a></button>
					<div style="font-family: "Trebuchet MS";font-size: 0.85em;">
						<p>%s</p>
					</div>
				</body>
			</html>
		`, confirmationId, confirmationId)
}

func TeamInviteTemplate(confirmationId string) string {
	return fmt.Sprintf(
		`
		<html> 
			<body style="background-color: white;padding: 20px;">
				<h1 style="text-align: center;"> Clipsynphony </h1>
				<button><a href="%s">Verify Email</a></button>
				<div style="font-family: "Trebuchet MS";font-size: 0.85em;">
					<p>%s</p>
				</div>
			</body>
		</html>
	`, confirmationId, confirmationId)
}

func AdminConfirmTemplate(confirmationId string, email string, username string) string {
	return fmt.Sprintf(
		`
		<html> 
			<body style="background-color: white;padding: 20px;">
				<h1 style="text-align: center;"> Tiermedizin </h1>
				<div> The user, %s with email %s wants to signup for Tiermedizin email services.</div>
				<div> Please confirm user signup by clicking the button or link below: </div>
				<button><a href="%s"> Confirm User Signup! </a></button>
				<div style="font-family: "Trebuchet MS";font-size: 0.85em;">
					<p>%s</p>
				</div>
			</body>
		</html>
	`, username, email, confirmationId, confirmationId)
}
