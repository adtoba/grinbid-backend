package controllers

import (
	"net/http"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListingController struct {
	DB *gorm.DB
}

func NewListingController(db *gorm.DB) *ListingController {
	return &ListingController{DB: db}
}
func (lc *ListingController) CreateListing(c *gin.Context) {
	var payload models.CreateListingRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	userID := c.MustGet("user_id").(string)

	var user models.User
	res := lc.DB.First(&user, "id = ?", userID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}

	listing := models.Listing{
		Title:       payload.Title,
		Description: payload.Description,
		Price:       payload.Price,
		Image:       payload.Image,
		Status:      "active",
		Location:    payload.Location,
		Condition:   payload.Condition,
		CategoryID:  payload.CategoryID,
		User:        user.ToSimpleUserResponse(),
		UserID:      userID,
	}

	result := lc.DB.Create(&listing)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusCreated, models.SuccessResponse("listing created successfully", listing))
}

func (lc *ListingController) GetAllListings(c *gin.Context) {
	var listings []models.Listing

	userID := c.MustGet("user_id").(string)

	result := lc.DB.Find(&listings).Where("status = ? AND user_id != ?", "active", userID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("listings fetched successfully", listings))
}

func (lc *ListingController) GetAllListingsByUserID(c *gin.Context) {
	var listings []models.Listing
	result := lc.DB.Find(&listings, "user_id = ?", c.Param("user_id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("listings fetched successfully", listings))
}

func (lc *ListingController) GetMyListings(c *gin.Context) {
	var listings []models.Listing
	result := lc.DB.Find(&listings, "user_id = ?", c.MustGet("user_id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("listings fetched successfully", listings))
}

func (lc *ListingController) GetListing(c *gin.Context) {
	var listing models.Listing
	result := lc.DB.First(&listing, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("listing fetched successfully", listing))
}

func (lc *ListingController) GetListingByCategory(c *gin.Context) {
	var listings []models.Listing
	categoryID := c.Param("category_id")

	userID := c.MustGet("user_id").(string)

	result := lc.DB.Find(&listings, "category_id = ? AND status = ? AND user_id != ?", categoryID, "active", userID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("listings fetched successfully", listings))
}
