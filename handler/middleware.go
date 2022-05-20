package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jschaefer-io/IDaaS/repository"
	"github.com/jschaefer-io/IDaaS/utils"
)

func Logger(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := middleware.RequestLogger(&middleware.DefaultLogFormatter{
			Logger:  logger,
			NoColor: runtime.GOOS != "windows",
		})
		return handler(next)
	}
}

func Html(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		next.ServeHTTP(writer, request)
	})
}

func Json(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

func TokenAuth(tokenManager *utils.TokenManager, userRepo *repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := request.Header.Get("Authorization")
			if len(token) == 0 {
				ApiError(writer, http.StatusUnauthorized, "access denied")
				return
			}

			claims, err := tokenManager.ValidateWithTokenType(strings.ReplaceAll(token, "Bearer ", ""), utils.TokenTypeBearer)
			if err != nil {
				fmt.Println(err)
				ApiError(writer, http.StatusUnauthorized, "access denied")
				return
			}

			usr, err := userRepo.Get(claims["user"].(string))
			if err != nil {
				ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
			}
			ctx := context.WithValue(request.Context(), "user", usr)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
