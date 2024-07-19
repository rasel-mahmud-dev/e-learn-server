package handlers

import (
	"database/sql"
	"e-learn/internal/models/category"
	"e-learn/internal/models/subCategory"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

func GetSubCategories(c *gin.Context) {

	items, err := subCategory.GetAllBySelect(c, []string{"id", "title", "slug", "image", "description", "created_at"}, func(rows *sql.Rows, category *subCategory.SubCategoryWithCamelCaseJSON) error {
		return rows.Scan(
			&category.ID,
			&category.Title,
			&category.Slug,
			&category.Image,
			&category.Description,
			&category.CreatedAt,
		)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

func GetSubCategory(c *gin.Context) {

	slug := c.Query("slug")
	id := c.Query("id")

	if slug == "" && id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
		return
	}

	var err error
	var bySlug *subCategory.SubCategoryWithCamelCaseJSON
	if slug != "" {
		bySlug, err = subCategory.GetOneBySlug(c, []string{"id", "title", "slug", "description", "image"}, func(row *sql.Row, json *subCategory.SubCategoryWithCamelCaseJSON) error {
			return row.Scan(&json.ID, &json.Title, &json.Slug, &json.Image, &json.Description)
		}, slug)
		if err != nil {
			return
		}
	} else {
		parseInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return
		}

		bySlug, err = subCategory.GetOneById(c, []string{"id", "title", "slug", "description", "image"}, func(row *sql.Row, json *subCategory.SubCategoryWithCamelCaseJSON) error {
			return row.Scan(&json.ID, &json.Title, &json.Slug, &json.Image, &json.Description)
		}, uint64(parseInt))
		if err != nil {
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bySlug,
	})
}

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

func CreateCategories(c *gin.Context) {
	var categoriesBody []category.CategoryWithCamelCaseJSON

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&categoriesBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.BatchInsert(c, categoriesBody)

	c.JSON(http.StatusCreated, gin.H{
		"data": categoriesBody,
	})
}

func CreateSubCategories(c *gin.Context) {
	var subCategoriesBody []subCategory.SubCategoryWithCamelCaseJSON

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&subCategoriesBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subCategory.BatchInsert(c, subCategoriesBody)

	c.JSON(http.StatusCreated, gin.H{
		"data": subCategoriesBody,
	})
}

//func CreateTopics(c *gin.Context) {
//	var titles []string
//
//	// Bind JSON or form data
//	if err := c.ShouldBindJSON(&titles); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var topics []models.Topics
//	for i := range titles {
//		title := titles[i]
//		topics = append(topics, models.Topics{
//			Title: title,
//			Slug:  utils.Slugify(title),
//		})
//
//	}
//
//	// Batch insert topics
//	if err := database.DB.Create(&topics).Error; err != nil {
//
//		response.ErrorResponse(c, err, map[string]string{
//			"unique": "Item already exist.",
//		})
//		return
//	}
//
//	c.JSON(http.StatusCreated, topics)
//
//}
