package handler

import (
	"net/http"
	"net/url"

	"github.com/jschaefer-io/IDaaS/utils"
)

func (c *WebController) HandleUserConfirm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Check that token exists
		tokens := request.URL.Query()["token"]
		if len(tokens) == 0 {
			c.Error(writer, http.StatusUnauthorized)
			return
		}

		// Validate given confirm token
		claims, err := c.components.TokenManager.ValidateWithTokenType(tokens[0], utils.TokenTypeConfirm)
		if err != nil {
			c.Error(writer, http.StatusUnauthorized)
			return
		}

		// Get User from token
		usr, err := c.components.Repositories.UserRepository.Get(claims["user"].(string))
		if err != nil {
			c.Error(writer, http.StatusInternalServerError)
			return
		}

		// confirm the user
		if !usr.Confirmed {
			usr.Confirmed = true
			_, err = c.components.Repositories.UserRepository.Persist(usr)
			if err != nil {
				c.Error(writer, http.StatusInternalServerError)
				return
			}
		}

		// Redirect to target with access token
		redirectUrl, _ := url.Parse(claims["redirect"].(string))
		query := redirectUrl.Query()
		query.Set("jio-confirm", "1")
		redirectUrl.RawQuery = query.Encode()
		http.Redirect(writer, request, redirectUrl.String(), http.StatusSeeOther)
	}
}
