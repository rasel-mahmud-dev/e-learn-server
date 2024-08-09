package topicHandler

import (
	"database/sql"
	"e-learn/internal/database"
	"e-learn/internal/models/category"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StoreTopicPreference(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.PostForm("user_id") // For UUID as a string
	ipAddress := c.PostForm("ip_address")
	deviceInfo := c.PostForm("device_info")

	rank := 0.1
	preferenceScore := 0.1

	columns := []string{
		"id",
		"title",
	}

	var err error

	bySlug, err := category.GetOne(c, columns, func(row *sql.Row, data *category.CategoryWithCamelCaseJSON) error {
		return row.Scan(&data.ID, &data.Title)
	}, "where slug = $1 AND type = $2", []any{slug, "topic"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameter."})
		return
	}

	topicID := bySlug.ID
	//defer database.DB.Close()

	// Check if the record exists
	var count int
	var updateQuery string

	if userID != "" {
		// Update query for user_id
		updateQuery = `
			UPDATE customer_topic_preference
			SET rank = rank + $5,
			    preference_score = preference_score + $6
			WHERE topic_id = $1 AND user_id = $2
		`
		err = database.DB.QueryRow(`
			SELECT COUNT(*)
			FROM customer_topic_preference
			WHERE topic_id = $1 AND user_id = $2
		`, topicID, userID).Scan(&count)

	} else {

		updateQuery = `
			UPDATE customer_topic_preference
			SET rank = rank + $4,
			    preference_score = preference_score + $5
			WHERE topic_id = $1 AND ip_address = $2 AND device_info = $3
		`
		err = database.DB.QueryRow(`
			SELECT COUNT(*)
			FROM customer_topic_preference
			WHERE topic_id = $1 AND ip_address = $2 AND device_info = $3
		`, topicID, ipAddress, deviceInfo).Scan(&count)
	}
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		_, err = database.DB.Exec(updateQuery, topicID, ipAddress, deviceInfo, rank, preferenceScore)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		insertQuery := `
			INSERT INTO customer_topic_preference
			(user_id, ip_address, device_info, topic_id, rank, preference_score)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = database.DB.Exec(insertQuery, userID, ipAddress, deviceInfo, topicID, rank, preferenceScore)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func PopularTopic(c *gin.Context) {

	query := `
		SELECT t.slug,  t.title
		FROM customer_topic_preference ctp
		JOIN categories t ON ctp.topic_id = t.id AND t.type = 'topic'
		GROUP BY t.title, t.slug
		ORDER BY SUM(ctp.preference_score) DESC
		LIMIT 100
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var topics []category.CategoryWithCamelCaseJSON
	for rows.Next() {
		var topic category.CategoryWithCamelCaseJSON
		if err := rows.Scan(&topic.Slug, &topic.Title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": topics})
}
