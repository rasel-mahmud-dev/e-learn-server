package handlers

import (
	"e-learn/internal/database"
	"e-learn/internal/models"
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
	"net/http"
	"time"
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
			ID:          user.ID,
			Thumbnail:   user.Thumbnail,
			Title:       user.Title,
			Slug:        user.Slug,
			Description: user.Description,
			AuthorID:    user.AuthorID,
			PublishDate: user.PublishDate,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Price:       user.Price,
		})
	}

	c.JSON(http.StatusOK, response)

}

func GetCourseDetail(c *gin.Context) {
	slug := c.Param("slug")
	var course models.Course

	row := database.DB.Table("courses").
		Select("id,  title, created_at, thumbnail, description, publish_date, price, slug").
		Where("slug = ?", slug).
		Row()

	err := row.Scan(&course.ID, &course.Title, &course.CreatedAt, &course.Thumbnail, &course.Description, &course.PublishDate, &course.Price, &course.Slug)
	if err != nil {
		response.ErrorResponse(c, err, map[string]string{})
		return
	}

	c.JSON(http.StatusOK, course)

}

func CreateCourse(c *gin.Context) {
	var data models.Course

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tim3e := time.Now()

	newCourse := models.Course{
		Thumbnail:   data.Thumbnail,
		Title:       data.Title,
		Slug:        utils.Slugify(data.Title),
		Description: data.Description,
		AuthorID:    data.AuthorID,
		PublishDate: &tim3e,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Price:       data.Price,
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

func CreateCourseBatch(c *gin.Context) {
	var data []models.Course

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tim3e := time.Now()

	var newCourses []models.Course

	for _, datum := range data {
		newCourse := models.Course{
			Thumbnail:   datum.Thumbnail,
			Title:       datum.Title,
			Slug:        utils.Slugify(datum.Title),
			Description: datum.Description,
			AuthorID:    datum.AuthorID,
			PublishDate: &tim3e,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Price:       datum.Price,
			//Categories:  [],
		}
		newCourses = append(newCourses, newCourse)
	}

	// Batch insert topics
	if err := database.DB.Create(&newCourses).Error; err != nil {
		response.ErrorResponse(c, err, map[string]string{
			"unique": "Item already exist.",
		})
		return
	}

	c.JSON(http.StatusCreated, newCourses)
}
