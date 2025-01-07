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

type TaskInterface interface {
	CreateTask(c *gin.Context)
	GetTaskByID(c *gin.Context)
	GetAllTasks(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}

type Task struct{}

func (x *Task) CreateTask(c *gin.Context) {
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

	var newTask models.Task

	newTask.UserID = intId

	err = c.BindJSON(&newTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	result := database.DB.Create(&newTask)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to create task"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(newTask))
}

func (x *Task) GetTaskByID(c *gin.Context) {
	user_id := c.Param("user_id")
	task_id := c.Param("task_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var task models.Task

	result := database.DB.Where("user_id = ? AND id = ?", user_id, task_id).First(&task)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.NewErrorResponse("No records found"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(task))
}

func (x *Task) GetAllTasks(c *gin.Context) {
	user_id := c.Param("user_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var tasks []models.Task

	result := database.DB.Where("user_id = ?", user_id).Find(&tasks)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helper.NewErrorResponse("No records found"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(tasks))
}

func (x *Task) UpdateTask(c *gin.Context) {
	user_id := c.Param("user_id")
	task_id := c.Param("task_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var updatedTask models.Task

	err := c.Bind(&updatedTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse("Invalid input"))
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to start transaction"))
		return
	}

	result := tx.Model(&updatedTask).Where("user_id = ? AND id = ?", user_id, task_id).Updates(&updatedTask)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to update task"))
		return
	}

	var task models.Task

	result = tx.Where("user_id = ? AND id = ?", user_id, task_id).First(&task)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to fetch updated task"))
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, helper.NewSuccessResponse(task))
}

func (x *Task) DeleteTask(c *gin.Context) {
	user_id := c.Param("user_id")
	task_id := c.Param("task_id")

	_, authorized := helper.CheckAuthenticationAndAuthorization(c, user_id)
	if !authorized {
		return
	}

	var task models.Task

	result := database.DB.Where("user_id = ? AND id = ?", user_id, task_id).Delete(&task)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helper.NewErrorResponse("Failed to delete task"))
		return
	}

	c.JSON(http.StatusOK, helper.NewSuccessResponse(nil))
}
