package main

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/go-chi/chi"
	"github.com/jschaefer-io/IDaaS/models"
	"net/http"
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
			Name:     data.Name,
			Email:    data.Email,
			Password: data.Password, // @todo hash it
		}
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
