package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jschaefer-io/IDaaS/repository"
	"github.com/jschaefer-io/IDaaS/utils"
)

func (c *ApiController) TokenCheck() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		token := request.Header.Get("Authorization")
		token = strings.ReplaceAll(token, "Bearer ", "")
		_, err := c.components.TokenManager.ValidateWithTokenType(token, utils.TokenTypeBearer)
		if err != nil {
			ApiError(writer, http.StatusUnauthorized, "access denied")
			return
		}
		respData, _ := json.Marshal(errorResponse{
			Code:  200,
			Error: false,
		})
		_, _ = writer.Write(respData)
	}
}

func (c *ApiController) createTokenPair(userId string, key string, token *repository.RefreshChain) (string, string, error) {
	// generate bearer token
	bearerToken, err := c.components.TokenManager.NewBearerToken(userId)
	if err != nil {
		return "", "", err
	}

	// create and persist new refresh-chain or use the given one, if it exists
	var refreshChain *repository.RefreshChain
	var tokenId string
	if token == nil {
		expiration, _ := c.components.TokenManager.GetTypeExpiration(utils.TokenTypeRefresh)
		refreshChain = c.components.Repositories.RefreshChainRepository.Make(key, expiration)
		tokenId, err = c.components.Repositories.RefreshChainRepository.Persist(refreshChain)
		if err != nil {
			return "", "", err
		}
	} else {
		refreshChain = token
		tokenId = refreshChain.Id
	}

	// build token string and return to response
	refreshToken, err := c.components.TokenManager.NewRefreshToken(userId, tokenId, refreshChain.Key)
	if err != nil {
		return "", "", err
	}

	return bearerToken, refreshToken, nil
}

func (c *ApiController) GrantTokens() http.HandlerFunc {
	type RequestData struct {
		AccessToken string `json:"accessToken"`
	}
	type ResponseData struct {
		BearerToken  string `json:"bearerToken"`
		RefreshToken string `json:"refreshToken"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		reqData := RequestData{}
		err := json.NewDecoder(request.Body).Decode(&reqData)
		if err != nil {
			ApiError(writer, http.StatusBadRequest, "invalid request")
			return
		}

		// validate access token
		claims, err := c.components.TokenManager.ValidateWithTokenType(reqData.AccessToken, utils.TokenTypeAccess)
		if err != nil {
			ApiError(writer, http.StatusUnauthorized, "access denied")
			return
		}

		userId := claims["user"].(string)
		key := claims["key"].(string)

		// reject token if a token with the given access_key is already active
		if _, err := c.components.Repositories.RefreshChainRepository.Find("access_key", key); err == nil {
			ApiError(writer, http.StatusUnauthorized, "access denied")
			return
		}

		bearerToken, refreshToken, err := c.createTokenPair(userId, key, nil)
		if err != nil {
			fmt.Println(err)
			ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
			return
		}

		response, _ := json.Marshal(ResponseData{
			BearerToken:  bearerToken,
			RefreshToken: refreshToken,
		})
		_, _ = writer.Write(response)
	}
}

func (c *ApiController) TokenRefresh() http.HandlerFunc {
	type RequestData struct {
		RefreshToken string `json:"refreshToken"`
	}
	type ResponseData struct {
		BearerToken  string `json:"bearerToken"`
		RefreshToken string `json:"refreshToken"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		reqData := RequestData{}
		err := json.NewDecoder(request.Body).Decode(&reqData)
		if err != nil {
			ApiError(writer, http.StatusBadRequest, "bad request")
			return
		}

		val := utils.NewValidator()
		val.Validate("refreshToken", reqData.RefreshToken, utils.RuleRequired())
		if !val.ErrorBag.Empty() {
			ApiError(writer, http.StatusUnprocessableEntity, val.ErrorBag.Errors())
			return
		}

		claims, err := c.components.TokenManager.ValidateWithTokenType(reqData.RefreshToken, utils.TokenTypeRefresh)
		if err != nil {
			ApiError(writer, http.StatusUnauthorized, "access denied")
			return
		}

		fmt.Println(claims)
		refreshChain, err := c.components.Repositories.RefreshChainRepository.Get(claims["id"].(string))
		if err != nil {
			fmt.Println(err)
			fmt.Println(err == sql.ErrNoRows)
			ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
			//ApiError(writer, http.StatusUnauthorized, "access denied")
			return
		}

		// if the reference keys are not equal, we
		// assume a compromised token and invalidate the
		// entire refresh-chain by deleting it
		if refreshChain.Key != claims["key"].(string) {
			if err = c.components.Repositories.RefreshChainRepository.Delete(refreshChain.Id); err != nil {
				ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
			} else {
				ApiError(writer, http.StatusUnauthorized, "access denied")
			}
			return
		}

		// Generate new key and persist it to the database
		key, err := utils.GenerateRandomString(10)
		if err != nil {
			ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
			return
		}

		refreshChain.Key = key
		if _, err = c.components.Repositories.RefreshChainRepository.Persist(refreshChain); err != nil {
			fmt.Println(err)
			ApiError(writer, http.StatusInternalServerError, "an unexpected error occurred")
			return
		}

		bearerToken, refreshToken, err := c.createTokenPair(claims["user"].(string), key, refreshChain)
		jsonTokens, _ := json.Marshal(ResponseData{
			BearerToken:  bearerToken,
			RefreshToken: refreshToken,
		})
		_, _ = writer.Write(jsonTokens)
	}
}
