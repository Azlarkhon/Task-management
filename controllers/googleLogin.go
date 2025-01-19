package controllers

import (
	"encoding/json"
	"net/http"
	"task-management/config"
	"task-management/database"
	"task-management/helper"
	"task-management/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleInterface interface {
	GoogleLogin(c *gin.Context)
	GoogleCallback(c *gin.Context)
}

type Google struct{}

type UserInfo struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

var googleOauthConfig = &oauth2.Config{
	ClientID:     config.Init.ClientID,
	ClientSecret: config.Init.ClientSecret,
	RedirectURL:  config.Init.RedirectURL,
	Scopes:       []string{config.Init.ScopeProfile, config.Init.ScopeEmail},
	Endpoint:     google.Endpoint,
}

var oauthStateString = "randomstate"

func (x *Google) GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (x *Google) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid state parameter"))
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to exchange token"))
		return
	}

	client := googleOauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to get user info"))
		return
	}
	defer resp.Body.Close()

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to decode user info"))
		return
	}

	var newUser models.User

	result := database.DB.Where("email = ?", userInfo.Email).First(&newUser)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			newUser.Email = userInfo.Email
			newUser.Username = userInfo.Name

			result = database.DB.Create(&newUser)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to create user"))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Database error"))
			return
		}
	}

	jwtToken, err := helper.GenerateJWT(newUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to generate token"))
		return
	}

	redirectURL := config.Init.ServerIp + "/user/" + jwtToken
    c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
