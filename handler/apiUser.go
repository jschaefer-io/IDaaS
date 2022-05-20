package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jschaefer-io/IDaaS/repository"
	"github.com/jschaefer-io/IDaaS/utils"
	"gopkg.in/gomail.v2"
)

func (c *ApiController) AddUser() http.HandlerFunc {
	type RequestData struct {
		Redirect  string `json:"redirect"`
		Gender    string `json:"gender"`
		Email     string `json:"email"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Password  string `json:"password"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		reqData := RequestData{}
		if err := json.NewDecoder(request.Body).Decode(&reqData); err != nil {
			ApiError(writer, http.StatusBadRequest, err.Error())
			return
		}

		// validate input data
		val := utils.NewValidator()
		val.Validate("redirect", reqData.Redirect, utils.RuleRequired(), utils.RuleUrl())
		val.Validate("gender", reqData.Gender, utils.RuleRequired(), utils.RuleIsIn(repository.GenderMale, repository.GenderFemale, repository.GenderOther))
		val.Validate("email", reqData.Email, utils.RuleRequired(), utils.RuleEmail())
		val.Validate("firstname", reqData.Firstname, utils.RuleRequired())
		val.Validate("lastname", reqData.Lastname, utils.RuleRequired())
		val.Validate("password", reqData.Password, utils.RuleRequired(), utils.RulePassword())
		if !val.ErrorBag.Empty() {
			ApiError(writer, http.StatusUnprocessableEntity, val.ErrorBag.Errors())
			return
		}

		// check if a user with the given email already exists
		var created bool
		usr, err := c.components.Repositories.UserRepository.Find("email", reqData.Email)
		if err == nil {
			// if a user with the given email already exists, who is already confirmed
			// else we'll resend the confirmation email
			if usr.Confirmed {
				ApiError(writer, http.StatusUnprocessableEntity, map[string]string{
					"email": "user with this email already exists",
				})
				return
			}
		} else {
			// if no user exist, we create a new unconfirmed user
			usr = c.components.Repositories.UserRepository.Make()
			usr.Gender = repository.Gender(reqData.Gender)
			usr.Email = reqData.Email
			usr.Firstname = reqData.Firstname
			usr.Lastname = reqData.Lastname
			usr.Password = reqData.Password
			id, dbErr := c.components.Repositories.UserRepository.Persist(usr)
			usr.ID = id
			if dbErr != nil {
				fmt.Println(dbErr)
				ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
				return
			}
			created = true
		}

		// Generate Confirm-Token
		token, err := c.components.TokenManager.NewConfirmToken(usr.ID, reqData.Redirect)
		if err != nil {
			ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
			return
		}

		// Send ConfirmationE-Mail
		m := gomail.NewMessage()
		m.SetHeader("From", c.settings.Mail.From)
		m.SetHeader("To", usr.Email)
		m.SetHeader("Subject", "User-Confirmation")
		mailBody, _ := c.components.Templates.ExecuteToString("mail", "confirmation", utils.GetUserConfirmUrl(c.settings.Url, map[string]string{
			"token": token,
		}))
		m.SetBody("text/html", mailBody)
		if err = c.components.Mailer.DialAndSend(m); err != nil {
			ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
			return
		}

		// Return User
		if created {
			writer.WriteHeader(http.StatusCreated)
		}
		jsonUsr, _ := json.Marshal(usr)
		_, _ = writer.Write(jsonUsr)
	}
}

func (c *ApiController) GetUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		usr := request.Context().Value("user").(*repository.User)
		jsonUsr, _ := json.Marshal(usr)
		_, _ = writer.Write(jsonUsr)
	}
}

func (c *ApiController) UpdateUser() http.HandlerFunc {
	type RequestData struct {
		Password string `json:"password"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		reqData := RequestData{}
		err := json.NewDecoder(request.Body).Decode(&reqData)
		if err != nil {
			ApiError(writer, http.StatusBadRequest, "bad request")
			return
		}

		val := utils.NewValidator()
		val.Validate("password", reqData.Password, utils.RuleRequired(), utils.RulePassword())
		if !val.ErrorBag.Empty() {
			ApiError(writer, http.StatusUnprocessableEntity, val.ErrorBag.Errors())
			return
		}

		usr := request.Context().Value("user").(*repository.User)
		usr.Password = reqData.Password

		_, err = c.components.Repositories.UserRepository.Persist(usr)
		if err != nil {
			ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
			return
		}

		jsonUsr, _ := json.Marshal(usr)
		_, _ = writer.Write(jsonUsr)
	}
}
