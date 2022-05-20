package server

import (
	"log"

	"github.com/jschaefer-io/IDaaS/repository"

	"github.com/go-chi/chi/v5"
	"github.com/jschaefer-io/IDaaS/utils"
	"github.com/jschaefer-io/IDaaS/view"
	"gopkg.in/gomail.v2"
)

type Components struct {
	Logger       *log.Logger
	Router       chi.Router
	Mailer       *gomail.Dialer
	TokenManager *utils.TokenManager
	Templates    *view.TemplateList
	Repositories *Repositories
}

type Repositories struct {
	UserRepository         *repository.UserRepository
	RefreshChainRepository *repository.RefreshChainRepository
}

func NewComponents(settings *Settings, repositories *Repositories) *Components {
	return &Components{
		Mailer: gomail.NewDialer(
			settings.Mail.Host,
			settings.Mail.Port,
			settings.Mail.Username,
			settings.Mail.Password,
		),
		TokenManager: utils.NewTokenManager(settings.Token.Secret),
		Templates:    view.NewTemplateList(settings.Static.Dir),
		Repositories: repositories,
	}
}
