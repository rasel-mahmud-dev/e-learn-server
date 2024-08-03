package handlers

import (
	"database/sql"
	"e-learn/internal/database"
	"e-learn/internal/models/category"
	"encoding/json"
	"errors"
	"fmt"
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

func GetTopics(c *gin.Context) {

	categories, err := category.GetAllBySelect(c, []string{"id", "title", "slug", "image", "description", "created_at"}, func(rows *sql.Rows, category *category.CategoryWithCamelCaseJSON) error {
		return rows.Scan(
			&category.ID,
			&category.Title,
			&category.Slug,
			&category.Image,
			&category.Description,
			&category.CreatedAt,
		)
	}, "where type = 'topic' ")

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

	columns := []string{"id", "title", "slug", "description", "image",
		"(select JSON_AGG(category_id::text) from subcategory_categories as category_ids where sub_category_id = categories.id)",
	}

	var err error
	var bySlug *category.CategoryWithCamelCaseJSON
	if slug != "" {
		bySlug, err = category.GetOne(c, columns, func(row *sql.Row, data *category.CategoryWithCamelCaseJSON) error {
			var CategoryIds []byte
			err := row.Scan(&data.ID, &data.Title, &data.Slug, &data.Image, &data.Description, &CategoryIds)
			if err != nil {
				return err
			}
			if len(CategoryIds) == 0 {
				data.CategoryIds = nil
			} else {
				if err := json.Unmarshal(CategoryIds, &data.CategoryIds); err != nil {
					return err
				}
			}
			return nil
		}, "where slug = $1 AND type = $2", []any{slug, "subcategory"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
			return
		}

	} else {
		bySlug, err = category.GetOne(c, columns, func(row *sql.Row, data *category.CategoryWithCamelCaseJSON) error {
			var CategoryIds []byte
			err := row.Scan(&data.ID, &data.Title, &data.Slug, &data.Image, &data.Description, &CategoryIds)
			if err != nil {
				return err
			}
			if len(CategoryIds) == 0 {
				data.CategoryIds = nil
			} else {
				if err := json.Unmarshal(CategoryIds, &data.CategoryIds); err != nil {
					return err
				}
			}
			return nil
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

func GetTopic(c *gin.Context) {

	slug := c.Query("slug")
	id := c.Query("id")

	if slug == "" && id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
		return
	}

	columns := []string{
		"id",
		"title",
		"slug",
		"(select JSON_AGG(sub_category_id::text) from topic_subcategories as sub_category_ids where topic_id = categories.id)",
	}

	var err error
	var bySlug *category.CategoryWithCamelCaseJSON
	if slug != "" {
		bySlug, err = category.GetOne(c, columns, func(row *sql.Row, data *category.CategoryWithCamelCaseJSON) error {
			var subCategoryIds []byte

			row.Scan(&data.ID, &data.Title, &data.Slug, &subCategoryIds)

			if len(subCategoryIds) == 0 {
				data.SubCategoryIds = nil
			} else {
				if err := json.Unmarshal(subCategoryIds, &data.SubCategoryIds); err != nil {
					return err
				}
			}
			return nil

		}, "where slug = $1 AND type = $2", []any{slug, "topic"})
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
			return
		}

	} else {
		bySlug, err = category.GetOne(c, columns, func(row *sql.Row, json *category.CategoryWithCamelCaseJSON) error {
			return row.Scan(&json.ID, &json.Title, &json.Slug, &json.Image, &json.Description, &json.SubCategoryIds)
		}, "where id = $1 AND type = $2", []any{id, "topic"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bySlug,
	})
}

func UpdateTopic(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
		return
	}

	payload := struct {
		Title          string   `json:"title"`
		SubCategoryIds []string `json:"subCategories"`
	}{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction."})
		return
	}

	var topicID int
	err = tx.QueryRow("SELECT id FROM categories WHERE slug = $1 AND type = 'topic'", slug).Scan(&topicID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topic."})
		}
		return
	}

	// Update the topic title
	_, err = tx.Exec("UPDATE categories SET title = $1 WHERE id = $2", payload.Title, topicID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update topic."})
		return
	}

	// Clear existing subcategory associations
	_, err = tx.Exec("DELETE FROM topic_subcategories WHERE topic_id = $1", topicID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing subcategories."})
		return
	}

	// Insert new subcategory associations
	stmt, err := tx.Prepare("INSERT INTO topic_subcategories (sub_category_id, topic_id) VALUES ($1, $2)")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement."})
		return
	}
	defer stmt.Close()

	for _, subCategoryID := range payload.SubCategoryIds {
		_, err = stmt.Exec(subCategoryID, topicID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert subcategory association."})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": payload,
	})
}

func UpdateSubCategory(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
		return
	}

	payload := struct {
		Title      string   `json:"title"`
		Categories []string `json:"categories"`
	}{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction."})
		return
	}

	var subCategoryID int
	err = tx.QueryRow("SELECT id FROM categories WHERE slug = $1 AND type = 'subcategory'", slug).Scan(&subCategoryID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory not found."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topic."})
		}
		return
	}

	// Update the topic title
	_, err = tx.Exec("UPDATE categories SET title = $1 WHERE id = $2", payload.Title, subCategoryID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update topic."})
		return
	}

	// Clear existing subcategory associations
	_, err = tx.Exec("DELETE FROM subcategory_categories WHERE sub_category_id = $1", subCategoryID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing subcategories."})
		return
	}

	// Insert new subcategory associations
	stmt, err := tx.Prepare("INSERT INTO subcategory_categories (sub_category_id, category_id) VALUES ($1, $2)")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement."})
		return
	}
	defer stmt.Close()

	for _, categoryID := range payload.Categories {
		_, err = stmt.Exec(subCategoryID, categoryID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert subcategory association."})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": payload,
	})
}

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
