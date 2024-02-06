package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"unhush-backend/models"
)

func LoginAndGetProfile() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		var exchangeToken models.ExchangeToken
		bindJSONError := ginContext.BindJSON(&exchangeToken)
		if bindJSONError != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": bindJSONError.Error()})
			return
		}

		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("code", exchangeToken.Code)
		data.Set("redirect_uri", exchangeToken.RedirectURL)
		data.Set("client_id", os.Getenv("CLIENT_ID"))
		data.Set("client_secret", os.Getenv("PRIMARY_CLIENT_SECRET"))

		accessTokenResponse, accessTokenError := http.PostForm(
			"https://www.linkedin.com/oauth/v2/accessToken",
			data)

		if accessTokenError != nil {
			ginContext.JSON(http.StatusInternalServerError,
				gin.H{"error": "Failed to request token"})
			log.Println(accessTokenError)
			return
		}

		defer accessTokenResponse.Body.Close()

		accessTokenBody, accessTokenBodyError := io.ReadAll(accessTokenResponse.Body)
		if accessTokenBodyError != nil {
			ginContext.JSON(http.StatusInternalServerError,
				gin.H{"error": "Failed to read response from token request"})
			log.Println(accessTokenBodyError)
			return
		}

		if accessTokenResponse.StatusCode < 200 && accessTokenResponse.StatusCode > 299 {
			ginContext.JSON(http.StatusInternalServerError,
				gin.H{"error": "Token request returned non-OK status",
					"status": accessTokenResponse.StatusCode,
					"body":   string(accessTokenBody)})
			return
		}

		var tokenResponse struct {
			AccessToken string `json:"access_token"`
		}

		unmarshallError := json.Unmarshal(accessTokenBody, &tokenResponse)
		if unmarshallError != nil {
			ginContext.JSON(http.StatusInternalServerError,
				gin.H{"error": "Failed to parse token response"})
			log.Println(unmarshallError)
			return
		}

		getRequest, getRequestError := http.NewRequest("GET",
			"https://api.linkedin.com/v2/userinfo",
			nil)

		if getRequestError != nil {
			ginContext.JSON(http.StatusInternalServerError,
				gin.H{"error": "Failed to make fetch profile request"})
			log.Println(getRequestError)
		}

		getRequest.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

		getProfileResponse, getProfileError := http.DefaultClient.Do(getRequest)
		if getProfileError != nil {
			ginContext.JSON(http.StatusInternalServerError,
				gin.H{"error": "Failed to fetch profile"})
			log.Println(getProfileError)
		}

		if getProfileResponse.StatusCode < 200 && getProfileResponse.StatusCode > 299 {
			ginContext.JSON(http.StatusInternalServerError,
				gin.H{"error": "Fetch profile returned non-OK status"})
			return
		}

		profileBody, profileBodyError := io.ReadAll(getProfileResponse.Body)
		if profileBodyError != nil {
			ginContext.JSON(http.StatusInternalServerError,
				gin.H{"error": "Failed to read response from get profile response"})
			log.Println(profileBodyError)
			return
		}

		var user models.User

		unmarshallError = json.Unmarshal(profileBody, &user)
		if unmarshallError != nil {
			ginContext.JSON(http.StatusInternalServerError,
				gin.H{"error": "Failed to parse profile body"})
			log.Println(unmarshallError)
			return
		}

		ginContext.JSON(http.StatusOK, &user)
	}
}
