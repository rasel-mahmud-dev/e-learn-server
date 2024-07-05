package handlers

import (
	"e-learn/internal/database"
	"e-learn/internal/models"
	"e-learn/internal/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCourses(c *gin.Context) {

	var users []models.Course
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var response []models.Course
	for _, user := range users {
		response = append(response, models.Course{
			ID:        user.ID,
			Thumbnail: "",
		})
	}

	c.JSON(http.StatusOK, response)

}

func CreateCourse(c *gin.Context) {
	var data models.Course

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCourse := models.Course{
		Thumbnail: data.Thumbnail,
		//Title:       data.Title,
		//Slug:        utils.Slugify(data.Title),
		//Description: data.Description,
		//AuthorID:    data.AuthorID,
		//PublishDate: time.Now(),
		//CreatedAt:   time.Now(),
		//UpdatedAt:   time.Now(),
		//Price:       data.Price,
		//Categories:  [],
	}

	// Batch insert topics
	if err := database.DB.Create(&newCourse).Error; err != nil {

		response.ErrorResponse(c, err, map[string]string{
			"unique": "Item already exist.",
		})
		return
	}

	c.JSON(http.StatusCreated, newCourse)
}
