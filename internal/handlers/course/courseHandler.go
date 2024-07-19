package courseHandler

import (
	"database/sql"
	"e-learn/internal/database"
	"e-learn/internal/models/category"
	"e-learn/internal/structType"
	"e-learn/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetCategories(c *gin.Context) {

	categories, err := category.GetAllBySelect(c, []string{"id", "title", "slug", "image", "description", "created_at"}, func(rows *sql.Rows, category *category.CategoryWithCamelCaseJSON) error {
		return rows.Scan(
			&category.ID,
			&category.Title,
			&category.Slug,
			&category.Image,
			&category.Description,
			&category.CreatedAt,
		)
	}, "where type = 'category' ")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

func GetSubCategory(c *gin.Context) {

	slug := c.Query("slug")
	id := c.Query("id")

	if slug == "" && id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
		return
	}

	columns := []string{"id", "title", "slug", "description", "image"}

	var err error
	var bySlug *category.CategoryWithCamelCaseJSON
	if slug != "" {
		bySlug, err = category.GetOne(c, columns, func(row *sql.Row, json *category.CategoryWithCamelCaseJSON) error {
			return row.Scan(&json.ID, &json.Title, &json.Slug, &json.Image, &json.Description)
		}, "where slug = $1 AND type = $2", []any{slug, "subcategory"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
			return
		}
	} else {
		bySlug, err = category.GetOne(c, columns, func(row *sql.Row, json *category.CategoryWithCamelCaseJSON) error {
			return row.Scan(&json.ID, &json.Title, &json.Slug, &json.Image, &json.Description)
		}, "where id = $1 AND type = $2", []any{id, "subcategory"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bySlug,
	})
}

//func UpdateSubCategory(c *gin.Context) {
//
//	id := c.Param("id")
//
//	if id == "" {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
//		return
//	}
//
//	//var err error
//	var bySlug *category.SubCategoryWithCamelCaseJSON
//
//	c.JSON(http.StatusOK, gin.H{
//		"data": bySlug,
//	})
//}

//func GetTopics(c *gin.Context) {
//
//	var users []models.Topics
//	result := database.DB.Find(&users)
//	if result.Error != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
//		return
//	}
//
//	var response []models.Topics
//	for _, user := range users {
//		response = append(response, models.Topics{
//			ID:    user.ID,
//			Title: user.Title,
//			Slug:  user.Slug,
//		})
//	}
//
//	c.JSON(http.StatusOK, response)
//
//}
//
//func CreateSubCategories(c *gin.Context) {
//	var titles []string
//
//	// Bind JSON or form data
//	if err := c.ShouldBindJSON(&titles); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var subcategories []models.SubCategory
//	for i := range titles {
//		title := titles[i]
//		subcategories = append(subcategories, models.SubCategory{
//			Title: title,
//			Slug:  utils.Slugify(title),
//		})
//	}
//
//	// Batch insert subcategories
//	if err := database.DB.Create(&subcategories).Error; err != nil {
//		panic(err)
//	}
//
//	c.JSON(http.StatusCreated, subcategories)
//
//}

func CreateCourse(c *gin.Context) {
	var createCoursePayload structType.CreateCoursePayload

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&createCoursePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createCourseSql := `insert into courses(
                    title,
                    slug,
                	thumbnail, 
                    description, 
                    publish_date, 
                    price,
                    created_at
                    )
		values ($1, $2, $3, $4, $5, $6, $7) returning id
`

	result, err := database.DB.ExecContext(
		c,
		createCourseSql,
		createCoursePayload.Title,
		utils.Slugify(createCoursePayload.Title),
		createCoursePayload.Thumbnail,
		createCoursePayload.Description,
		nil,
		createCoursePayload.Price,
		time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(result)

	c.JSON(http.StatusCreated, gin.H{
		"data": createCoursePayload,
	})
}
