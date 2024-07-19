package handlers

import (
	"database/sql"
	"e-learn/internal/models/category"
	"github.com/gin-gonic/gin"
	"net/http"
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

func GetSubCategories(c *gin.Context) {

	items, err := category.GetAllBySelect(
		c,
		[]string{"id", "title", "slug", "image", "description", "created_at"},
		func(rows *sql.Rows, category *category.CategoryWithCamelCaseJSON) error {
			return rows.Scan(
				&category.ID,
				&category.Title,
				&category.Slug,
				&category.Image,
				&category.Description,
				&category.CreatedAt,
			)
		}, "where type = 'subcategory' ")

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

func CreateCategories(c *gin.Context) {
	var categoriesBody []category.CategoryWithCamelCaseJSON

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&categoriesBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.BatchInsert(c, categoriesBody, "category")

	c.JSON(http.StatusCreated, gin.H{
		"data": categoriesBody,
	})
}

func CreateSubCategories(c *gin.Context) {
	var subCategoriesBody []category.CategoryWithCamelCaseJSON

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&subCategoriesBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.BatchInsert(c, subCategoriesBody, "subcategory")

	c.JSON(http.StatusCreated, gin.H{
		"data": subCategoriesBody,
	})
}

func CreateTopics(c *gin.Context) {
	var topicsBody []category.CategoryWithCamelCaseJSON

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&topicsBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.BatchInsert(c, topicsBody, "topic")

	c.JSON(http.StatusCreated, gin.H{
		"data": topicsBody,
	})
}
