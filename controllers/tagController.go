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

type TagInterface interface {
	GetTagById(c *gin.Context)
	GetAllTags(c *gin.Context)
	CreateTag(c *gin.Context)
	UpdateTagById(c *gin.Context)
	DeleteTagById(c *gin.Context)
	AddTagToTask(c *gin.Context)
	DeleteTagFromTask(c *gin.Context)
	GetAllTagsOfTask(c *gin.Context)
}

type Tag struct{}

func (x *Tag) GetTagById(c *gin.Context) {
	user_id := c.Param("user_id")
	tag_id := c.Param("tag_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var tag models.Tag

	result := database.DB.Where("user_id = ? AND id = ?", user_id, tag_id).First(&tag)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.NewErrorResponse("No records found"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(tag))
}

func (x *Tag) GetAllTags(c *gin.Context) {
	user_id := c.Param("user_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var tags []models.Tag

	result := database.DB.Where("user_id = ?", user_id).Find(&tags)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.NewErrorResponse("No records found"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(tags))
}

func (x *Tag) CreateTag(c *gin.Context) {
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

	var newTag models.Tag

	newTag.User_id = intId

	err = c.BindJSON(&newTag)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	result := database.DB.Create(&newTag)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to create tag"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(newTag))
}

func (x *Tag) UpdateTagById(c *gin.Context) {
	user_id := c.Param("user_id")
	tag_id := c.Param("tag_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var newTag models.Tag

	err := c.BindJSON(&newTag)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to start transaction"))
		return
	}

	result := tx.Model(&newTag).Where("user_id = ? AND id = ?", user_id, tag_id).Updates(&newTag)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to update user"))
		return
	}

	var updatedTag models.Tag

	result = tx.Where("user_id = ? AND id = ?", user_id, tag_id).First(&updatedTag)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to fetch updated tag"))
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, helper.NewSuccessResponse(updatedTag))
}

func (x *Tag) DeleteTagById(c *gin.Context) {
	user_id := c.Param("user_id")
	tag_id := c.Param("tag_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var user models.Tag

	result := database.DB.Where("user_id = ? AND id = ?", user_id, tag_id).Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to delete tag"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(nil))
}

func (x *Tag) AddTagToTask(c *gin.Context) {
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

	var task_tag models.Task_Tags

	task_tag.User_id = intId

	err = c.BindJSON(&task_tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	result := database.DB.Table("task_tags").Create(&task_tag)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to add tag to task"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(task_tag))
}

func (x *Tag) DeleteTagFromTask(c *gin.Context) {
	user_id := c.Param("user_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var task_tag models.Task_Tags

	err := c.BindJSON(&task_tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	result := database.DB.Table("task_tags").Where("user_id = ? AND tag_id = ? AND task_id = ?", user_id, task_tag.Tag_id, task_tag.Task_id).Delete(&task_tag)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to delete tag from task"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(nil))
}

func (x *Tag) GetAllTagsOfTask(c *gin.Context) {
	user_id := c.Param("user_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var task_tag models.Task_Tags

	err := c.BindJSON(&task_tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	var task_tags []models.Task_Tags

	result := database.DB.Table("task_tags").Where("user_id = ? AND task_id = ?", user_id, task_tag.Task_id).Find(&task_tags)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.NewErrorResponse("No records found"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(task_tags))
}
