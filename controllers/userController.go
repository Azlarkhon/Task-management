package controllers

import (
	"errors"
	"net/http"
	"task-management/database"
	"task-management/helper"
	"task-management/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserInterface interface {
	GetUserByID(c *gin.Context)
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type User struct{}

func (x *User) GetUserByID(c *gin.Context) {
	id := c.Param("user_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, id)
	if !authorized {
		return
	}

	var newUser models.User

	result := database.DB.Where("id = ?", id).First(&newUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, helper.NewErrorResponse("User not found"))
		} else {
			c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Database error"))
		}
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(newUser))
}

func (x *User) SignUp(c *gin.Context) {
	var newUser models.User

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	result := database.DB.Where("email = ?", newUser.Email).First(&newUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Email already registered"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to hash password"))
	}
	newUser.Password = string(hashedPassword)

	result = database.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to create user"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(newUser))
}

func (x *User) Login(c *gin.Context) {
	var inputUser models.User
	var storedUser models.User

	err := c.BindJSON(&inputUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	result := database.DB.Where("email = ?", inputUser.Email).First(&storedUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, helper.NewErrorResponse("Invalid email or password"))
		} else {
			c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Database error"))
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(inputUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, helper.NewErrorResponse("Invalid email or password"))
		return
	}

	token, err := helper.GenerateJWT(storedUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       storedUser.ID,
		"token":      token,
		"is_success": true,
	})
}

func (x *User) UpdateUser(c *gin.Context) {
	id := c.Param("user_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, id)
	if !authorized {
		return
	}

	var newUser models.User

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to start transaction"))
		return
	}

	result := tx.Model(&newUser).Where("id = ?", id).Updates(&newUser)

	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to update user"))
		return
	}

	var updatedUser models.User
	err = tx.First(&updatedUser, "id = ?", id).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to fetch updated user"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(updatedUser))
}

func (x *User) DeleteUser(c *gin.Context) {
	id := c.Param("user_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, id)
	if !authorized {
		return
	}

	var user models.User

	result := database.DB.Where("id = ?", id).Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to delete user"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(nil))
}
