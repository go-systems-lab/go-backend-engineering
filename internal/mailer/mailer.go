package mailer

import "embed"

const (
	FromName               = "SocialGo"
	MaxRetries             = 3
	UserInvitationTemplate = "user_invitation.tmpl"
)

//go:embed templates
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}
