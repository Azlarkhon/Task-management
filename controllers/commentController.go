package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"task-management/database"
	"task-management/helper"
	"task-management/models"

	"github.com/gin-gonic/gin"
)

type CommentInterface interface {
	CreateComment(c *gin.Context)
	DeleteComment(c *gin.Context)
	GetAllCommentsOfTask(c *gin.Context)
}

type Comment struct{}

func (x *Comment) CreateComment(c *gin.Context) {
	user_id := c.Param("user_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	intId, err := strconv.Atoi(user_id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var newComment models.Comment

	newComment.UserID = intId

	err = c.BindJSON(&newComment)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	result := database.DB.Create(&newComment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to create comment to task"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(newComment))
}

func (x *Comment) DeleteComment(c *gin.Context) {
	user_id := c.Param("user_id")
	comment_id := c.Param("comment_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var comment models.Comment

	result := database.DB.Where("id = ?", comment_id).Delete(&comment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to delete comment"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(nil))
}

func (x *Comment) GetAllCommentsOfTask(c *gin.Context) {
	user_id := c.Param("user_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var comment models.Comment

	err := c.BindJSON(&comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	var comments []models.Comment

	result := database.DB.Where("user_id = ? AND task_id = ?", user_id, comment.TaskID).Find(&comments)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.NewErrorResponse("No records found"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(comments))
}
