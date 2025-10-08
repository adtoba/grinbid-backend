package controllers

import (
	"net/http"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{DB: db}
}

func (ac *AdminController) GetAllUsers(c *gin.Context) {
	var users []models.User
	result := ac.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("users fetched successfully", users))
}

func (ac *AdminController) GetUser(c *gin.Context) {
	var user models.User
	result := ac.DB.First(&user, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("user fetched successfully", user))
}

func (ac *AdminController) BlockUser(c *gin.Context) {
	var user models.User
	result := ac.DB.Model(&user).Where("id = ?", c.Param("id")).Update("is_blocked", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("user blocked successfully", user))
}

func (ac *AdminController) UnblockUser(c *gin.Context) {
	var user models.User
	result := ac.DB.Model(&user).Where("id = ?", c.Param("id")).Update("is_blocked", false)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("user unblocked successfully", user))
}
