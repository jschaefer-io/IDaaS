package handler

import (
	"net/http"
	"net/url"

	"github.com/jschaefer-io/IDaaS/utils"
)

type loginData struct {
	Mail     string
	Redirect string
	ResetUrl string
	Errors   map[string]string
}

func (c *WebController) LoginForm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		redirect, err := c.getRedirectFromHeader(request)
		if err != nil {
			redirectUrl, _ := utils.AddQueryToUrl(c.settings.Redirect.Default, map[string]string{
				"jio-error": "redirect_url_error",
			})
			http.Redirect(writer, request, redirectUrl, http.StatusTemporaryRedirect)
			return
		}
		_ = c.components.Templates.Execute(writer, "login", "form", loginData{
			Mail:     "",
			Redirect: redirect,
			ResetUrl: utils.GetResetUrl(c.settings.Url, map[string]string{
				"redirect": redirect,
			}),
			Errors: nil,
		})
	}
}

func (c *WebController) HandleLoginForm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		email := request.FormValue("email")
		password := request.FormValue("password")
		redirect := request.FormValue("redirect")

		// Generate reset Url since we need it multiple times
		resetUrl := utils.GetResetUrl(c.settings.Url, map[string]string{
			"redirect": redirect,
		})

		// Validate Input
		val := utils.NewValidator()
		val.Validate("mail", email, utils.RuleRequired(), utils.RuleEmail())
		val.Validate("pw", password, utils.RuleRequired())
		if !val.ErrorBag.Empty() {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			_ = c.components.Templates.Execute(writer, "login", "form", loginData{
				Mail:     email,
				Redirect: redirect,
				ResetUrl: resetUrl,
				Errors:   val.ErrorBag.Errors(),
			})
			return
		}

		// Execute login
		usr, err := c.components.Repositories.UserRepository.Find("email", email)
		if err != nil || !utils.CheckPassword(password, usr.PasswordHash) {
			writer.WriteHeader(http.StatusUnauthorized)
			_ = c.components.Templates.Execute(writer, "login", "form", loginData{
				Mail:     email,
				Redirect: redirect,
				ResetUrl: resetUrl,
				Errors: map[string]string{
					"system": "permission denied",
				},
			})
			return
		}

		// Check if the user is unconfirmed
		if !usr.Confirmed {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			_ = c.components.Templates.Execute(writer, "login", "form", loginData{
				Mail:     email,
				Redirect: redirect,
				ResetUrl: resetUrl,
				Errors: map[string]string{
					"system": "user is not confirmed, check your emails",
				},
			})
			return
		}

		// Generate the JWT Access-Token
		accessToken, err := c.components.TokenManager.NewAccessToken(usr.ID)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = c.components.Templates.Execute(writer, "login", "form", loginData{
				Mail:     email,
				Redirect: redirect,
				ResetUrl: resetUrl,
				Errors: map[string]string{
					"system": "an error occurred, try again later",
				},
			})
			return
		}

		if !utils.SliceContains(c.settings.Redirect.Whitelist, redirect) {
			redirect = c.settings.Redirect.Default
		}

		// Redirect to target with access token
		redirectUrl, _ := url.Parse(redirect)
		query := redirectUrl.Query()
		query.Set("access-token", accessToken)
		redirectUrl.RawQuery = query.Encode()
		http.Redirect(writer, request, redirectUrl.String(), http.StatusSeeOther)
	}
}
