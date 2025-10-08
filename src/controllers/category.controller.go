package controllers

import (
	"net/http"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{DB: db}
}

func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var payload models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	category := models.Category{
		Name: payload.Name,
		Icon: payload.Icon,
	}

	result := cc.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("category created successfully", category))
}

func (cc *CategoryController) GetAllCategories(c *gin.Context) {
	var categories []models.Category
	result := cc.DB.Find(&categories)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("categories fetched successfully", categories))
}
