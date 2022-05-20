package handler

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/jschaefer-io/IDaaS/repository"
	"github.com/jschaefer-io/IDaaS/utils"
	"gopkg.in/gomail.v2"
)

type resetStartData struct {
	Mail     string
	Redirect string
	Errors   map[string]string
}

type resetFormData struct {
	Token  string
	Errors map[string]string
}

type resetSuccessTemplate struct {
	Step string
	Data map[string]any
}

func (c *WebController) PasswordResetStart() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		redirect, _ := c.getRedirectFromHeader(request)
		_ = c.components.Templates.Execute(writer, "reset", "start", resetStartData{
			Mail:     "",
			Redirect: redirect,
		})
	}
}

func (c *WebController) HandlePasswordResetStart() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		email := request.FormValue("email")
		redirect := request.FormValue("redirect")

		// Validate Input
		val := utils.NewValidator()
		val.Validate("mail", email, utils.RuleRequired(), utils.RuleEmail())
		if !val.ErrorBag.Empty() {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			_ = c.components.Templates.Execute(writer, "reset", "start", resetStartData{
				Mail:     email,
				Redirect: redirect,
				Errors:   val.ErrorBag.Errors(),
			})
			return
		}

		usr, err := c.components.Repositories.UserRepository.Find("email", email)
		if err != nil {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			_ = c.components.Templates.Execute(writer, "reset", "start", resetStartData{
				Mail:     email,
				Redirect: redirect,
				Errors: map[string]string{
					"system": "no user associated with this email",
				},
			})
			return
		}

		genericError := func() {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = c.components.Templates.Execute(writer, "reset", "start", resetStartData{
				Mail:     email,
				Redirect: redirect,
				Errors: map[string]string{
					"system": "internal server error",
				},
			})
		}

		// Generate Reset-Token
		resetToken, err := c.components.TokenManager.NewResetToken(usr.ID, usr.UpdateAt, redirect)
		if err != nil {
			genericError()
			return
		}

		// Send E-Mail
		m := gomail.NewMessage()
		m.SetHeader("From", c.settings.Mail.From)
		m.SetHeader("To", usr.Email)
		m.SetHeader("Subject", "RulePassword Reset")
		mailBody, _ := c.components.Templates.ExecuteToString("mail", "reset", utils.GetResetConfirmUrl(c.settings.Url, map[string]string{
			"token": resetToken,
		}))
		m.SetBody("text/html", mailBody)
		if err = c.components.Mailer.DialAndSend(m); err != nil {
			genericError()
			return
		}

		// Show success screen
		_ = c.components.Templates.Execute(writer, "reset", "success", resetSuccessTemplate{
			Step: "start",
		})

	}
}

func (c *WebController) resolveResetToken(fromQuery bool, request *http.Request) (string, *repository.User, jwt.MapClaims, error) {
	hasToken := false
	var token string
	if fromQuery {
		tokens := request.URL.Query()["token"]
		if len(tokens) != 0 {
			token = tokens[0]
			hasToken = true
		}
	} else {
		token = request.FormValue("token")
		if len(token) > 0 {
			hasToken = true
		}
	}
	if !hasToken {
		return "", nil, nil, errors.New("no token provided")
	}

	// validate token
	tokenClaims, err := c.components.TokenManager.ValidateWithTokenType(token, utils.TokenTypeReset)
	if err != nil {
		return token, nil, tokenClaims, errors.New("token expired")
	}

	// resolve token to a User
	usr, err := c.components.Repositories.UserRepository.Get(tokenClaims["user"].(string))
	if err != nil {
		return token, usr, tokenClaims, errors.New("associated user not found")
	}

	// check if token key matches
	check, err := utils.HashStringMd5(usr.UpdateAt)
	if err != nil || check != tokenClaims["key"].(string) {
		return token, usr, tokenClaims, errors.New("token key mismatch")
	}
	return token, usr, tokenClaims, nil
}

func (c *WebController) PasswordResetForm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		resetToken, _, _, err := c.resolveResetToken(true, request)
		if err != nil {
			_ = c.components.Templates.Execute(writer, "error", "server", nil)
			return
		}

		_ = c.components.Templates.Execute(writer, "reset", "form", resetFormData{
			Token: resetToken,
		})
	}
}

func (c *WebController) HandlePasswordResetForm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		resetToken, usr, claims, err := c.resolveResetToken(false, request)
		if err != nil {
			_ = c.components.Templates.Execute(writer, "error", "server", nil) // @todo expired error
			return
		}

		// Validate Input
		password := request.FormValue("password")
		passwordRepeat := request.FormValue("password-repeat")
		val := utils.NewValidator()
		val.Validate("pw", password, utils.RuleRequired(), utils.RulePassword())
		val.Validate("pw-repeat", passwordRepeat, utils.RuleRequired(), utils.RuleEqual(password))
		if !val.ErrorBag.Empty() {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			_ = c.components.Templates.Execute(writer, "reset", "form", resetFormData{
				Token:  resetToken,
				Errors: val.ErrorBag.Errors(),
			})
			return
		}

		usr.Password = password
		if _, err = c.components.Repositories.UserRepository.Persist(usr); err != nil {
			_ = c.components.Templates.Execute(writer, "error", "server", nil)
			return
		}

		_ = c.components.Templates.Execute(writer, "reset", "success", resetSuccessTemplate{
			Step: "reset",
			Data: map[string]any{
				"user": usr,
				"loginUrl": utils.GetLoginUrl(c.settings.Url, map[string]string{
					"redirect": claims["redirect"].(string),
				}),
			},
		})
	}
}
