package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/fatih/structs"
	"github.com/go-chi/chi"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sync"
)

func (s *Server) UsersGetAll() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		users := make([]models.User, 0)
		s.db.Find(&users)
		s.marshalToResponseWriter(users, writer, http.StatusOK)
	}
}

func (s *Server) UsersGetSingle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "userId")
		usr := models.User{}
		err := s.db.First(&usr, id)
		if err.Error != nil {
			s.marshalToResponseWriter(s.newErrorObject(fmt.Sprintf("user %v not found", id)), writer, http.StatusNotFound)
			return
		}
		s.marshalToResponseWriter(usr, writer, http.StatusOK)
	}
}

func (s *Server) UsersCreate() http.HandlerFunc {
	type userCreateForm struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		data := new(userCreateForm)
		if err := models.BindJson(data, request.Body); err != nil {
			s.marshalToResponseWriter(s.validationError(err), writer, http.StatusUnprocessableEntity)
			return
		}

		// Check if the given email address is unique
		uniqueCheck := make([]models.User, 0)
		s.db.Where("email = ?", data.Email).Find(&uniqueCheck)
		if len(uniqueCheck) > 0 {
			s.marshalToResponseWriter(s.newErrorObject(fmt.Sprintf("a user with this the email %s already exists", data.Email)), writer, http.StatusBadRequest)
			return
		}

		user := models.User{
			Name:  data.Name,
			Email: data.Email,
		}
		user.SetPassword(data.Password)

		if err := s.db.Create(&user); err.Error != nil {
			s.marshalToResponseWriter(s.newErrorObject(err.Error.Error()), writer, http.StatusInternalServerError)
			return
		}
		s.marshalToResponseWriter(user, writer, http.StatusCreated)
	}
}

func (s *Server) UserUpdateSingle() http.HandlerFunc {
	type userUpdateForm struct {
		Name     string `json:"name" validate:"omitempty"`
		Email    string `json:"email" validate:"omitempty,email"`
		Password string `json:"password" validate:"omitempty"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		// Prepare input data
		data := new(userUpdateForm)
		if err := models.BindJson(data, request.Body); err != nil {
			s.marshalToResponseWriter(s.validationError(err), writer, http.StatusUnprocessableEntity)
			return
		}
		updateData := make(map[string]interface{})
		for key, value := range structs.Map(data) {
			if value != "" {
				if key == "Password" {
					usr := models.User{}
					usr.SetPassword(value.(string))
					value = usr.Password
				}
				updateData[key] = value
			}
		}

		// Find User
		id := chi.URLParam(request, "userId")
		usr := models.User{}
		err := s.db.First(&usr, id)
		if err.Error != nil {
			s.marshalToResponseWriter(s.newErrorObject(fmt.Sprintf("user %v not found", id)), writer, http.StatusNotFound)
			return
		}

		// Update User
		if err := s.db.Model(&usr).Updates(updateData); err.Error != nil {
			s.marshalToResponseWriter(s.newErrorObject(err.Error.Error()), writer, http.StatusInternalServerError)
			return
		}
		s.marshalToResponseWriter(usr, writer, http.StatusOK)
	}
}

func (s *Server) UserDeleteSingle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "userId")
		del := s.db.Delete(models.User{}, id)
		if del.Error != nil {
			s.marshalToResponseWriter(s.newErrorObject(del.Error.Error()), writer, http.StatusInternalServerError)
			return
		}
		if del.RowsAffected == 0 {
			s.marshalToResponseWriter(s.newErrorObject(fmt.Sprintf("user %v not found", id)), writer, http.StatusNotFound)
			return
		}
		s.marshalToResponseWriter(nil, writer, http.StatusNoContent)
	}
}

func (s *Server) UserAuthenticate() http.HandlerFunc {
	type userLoginForm struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	var init sync.Once
	var rsaSecret *rsa.PrivateKey
	var rsaError error
	return func(writer http.ResponseWriter, request *http.Request) {
		init.Do(func() {
			rsaSecret, rsaError = crypto.ReadPrivateRsaKey()
		})
		if rsaError != nil {
			s.marshalToResponseWriter(s.newErrorObject("an unknown error occurred"), writer, http.StatusInternalServerError)
			panic(rsaError)
		}

		// check user credentials
		data := new(userLoginForm)
		usr := models.User{}
		if err := models.BindJson(data, request.Body); err != nil {
			s.marshalToResponseWriter(s.newErrorObject(err.Error()), writer, http.StatusInternalServerError)
			return
		}
		if err := s.db.Where("email = ?", data.Email).First(&usr); err.Error != nil {
			s.marshalToResponseWriter(s.newErrorObject(err.Error.Error()), writer, http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(data.Password)); err != nil {
			fmt.Println(err)
			s.marshalToResponseWriter(s.newErrorObject("permission denied"), writer, http.StatusUnauthorized)
			return
		}

		tokenString, err := crypto.NewJwt(rsaSecret, usr)
		if err != nil {
			s.marshalToResponseWriter(s.newErrorObject("an unknown error occurred"), writer, http.StatusInternalServerError)
			return
		}
		s.marshalToResponseWriter(map[string]interface{}{
			"token": tokenString,
		}, writer, http.StatusOK)
	}
}
