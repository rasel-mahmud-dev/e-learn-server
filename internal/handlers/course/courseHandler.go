package courseHandler

import (
	"database/sql"
	"e-learn/internal/database"
	"e-learn/internal/models/course"
	"e-learn/internal/response"
	"e-learn/internal/structType"
	"e-learn/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetInstructorCourses(c *gin.Context) {

	// check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("Unauthorization"), nil)
		return
	}

	columns := []string{
		"courses.id as id",
		"ac.course_id as course_id",
		"title",
		"slug",
		"thumbnail",
		"price",
		"created_at",
		"(select jsonb_agg(DISTINCT cs.category_id) from courses_categories cs where courses.course_id = cs.course_id) as categories",
		"(select jsonb_agg(DISTINCT sc.category_id) from courses_sub_categories sc where courses.course_id = sc.course_id) as sub_categories",
		"(select jsonb_agg(DISTINCT ct.topic_id)  from courses_topics ct where courses.course_id = ct.course_id)         as topics",
		"(select jsonb_agg(DISTINCT ac.author_id) from authors_courses ac where ac.course_id = courses.course_id)           as authors",
	}

	authJoin := `join authors_courses ac on courses.course_id = ac.course_id where ac.author_id = $1`

	courses, err := course.GetAllBySelect(c, columns, func(rows *sql.Rows, course *course.Course) error {
		return rows.Scan(
			&course.ID,
			&course.CourseID,
			&course.Title,
			&course.Slug,
			&course.Thumbnail,
			&course.Price,
			&course.CreatedAt,
			&course.CategoryListJson,
			&course.SubCategoryListJson,
			&course.TopicListJson,
			&course.AuthorListJson,
		)
	}, authJoin, []any{authUser.UserId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": courses,
	})
}

func GetCourses(c *gin.Context) {

	// Retrieve topics from query parameters
	topics := c.QueryArray("topic")
	durations := c.QueryArray("duration")

	// Construct placeholders for the query
	placeholders := make([]string, len(topics))
	for i := range topics {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	placeholderString := strings.Join(placeholders, ", ")

	// Construct placeholders for the durations
	durationPlaceholders := make([]string, len(durations))
	for i := range durations {
		durationPlaceholders[i] = fmt.Sprintf("$%d", len(topics)+i+1)
	}
	durationPlaceholderString := strings.Join(durationPlaceholders, ", ")

	query := fmt.Sprintf(`
    SELECT c.*
    FROM public.courses c
    JOIN public.courses_topics ct ON c.course_id = ct.course_id
    JOIN public.categories cat ON ct.topic_id = cat.id AND cat.type = 'topic'
    WHERE cat.slug IN (%s)  AND c.duration IN (%s)
    `, placeholderString, durationPlaceholderString)

	// Convert topics and durations to a slice of interface{}
	params := make([]interface{}, len(topics)+len(durations))
	for i, topic := range topics {
		params[i] = topic
	}
	for i, duration := range durations {
		params[len(topics)+i] = duration
	}
	fmt.Println(params)

	// Execute the query with topics as an array
	rows, err := database.DB.QueryContext(c, query, params...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching courses"})
		return
	}
	defer rows.Close()

	var courses []course.Course
	for rows.Next() {
		var course course.Course
		err := rows.Scan(
			&course.ID,
			&course.CourseID,
			&course.CreatedAt,
			&course.UpdatedAt,
			&course.DeletedAt,
			&course.Thumbnail,
			&course.Title,
			&course.Slug,
			&course.Description,
			&course.PublishDate,
			&course.Price,
			&course.Duration,
			&course.NumLectures,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing courses"})
			return
		}
		courses = append(courses, course)
	}
	//
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing courses"})
		return
	}

	// Respond with the courses data
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

func GetInstructorCourseDetail(c *gin.Context) {

	// check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("Unauthorization"), nil)
		return
	}

	slug := c.Param("slug")

	columns := []string{
		"courses.id as id",
		"ac.course_id as course_id",
		"title",
		"slug",
		"description",
		"thumbnail",
		"price",
		"created_at",
		"(select jsonb_agg(DISTINCT cs.category_id) from courses_categories cs where courses.course_id = cs.course_id) as categories",
		"(select jsonb_agg(DISTINCT sc.category_id) from courses_sub_categories sc where courses.course_id = sc.course_id) as sub_categories",
		"(select jsonb_agg(DISTINCT ct.topic_id)  from courses_topics ct where courses.course_id = ct.course_id)         as topics",
		"(select jsonb_agg(DISTINCT ac.author_id) from authors_courses ac where ac.course_id = courses.course_id)           as authors",
	}

	authJoin := `join authors_courses ac on courses.course_id = ac.course_id where ac.author_id = $1 AND courses.slug = $2`

	course, err := course.GetOne(c, columns, func(row *sql.Row, course *course.Course) error {
		return row.Scan(
			&course.ID,
			&course.CourseID,
			&course.Title,
			&course.Slug,
			&course.Description,
			&course.Thumbnail,
			&course.Price,
			&course.CreatedAt,
			&course.CategoryListJson,
			&course.SubCategoryListJson,
			&course.TopicListJson,
			&course.AuthorListJson,
		)
	}, authJoin, []any{authUser.UserId, slug})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = json.Unmarshal([]byte(*course.CategoryListJson), &course.CategoryList)
	course.CategoryListJson = nil
	err = json.Unmarshal([]byte(*course.SubCategoryListJson), &course.SubCategoryList)
	course.SubCategoryListJson = nil
	err = json.Unmarshal([]byte(*course.TopicListJson), &course.TopicList)
	course.TopicListJson = nil
	err = json.Unmarshal([]byte(*course.AuthorListJson), &course.AuthorList)
	course.AuthorListJson = nil

	c.JSON(http.StatusOK, gin.H{
		"data": course,
	})
}

func CreateCourse(c *gin.Context) {

	// check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		return
	}

	var createCoursePayload structType.CreateCoursePayload

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&createCoursePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createCourseSql := `insert into courses(
	                   course_id,
	                   title,
	                   slug,
	               	thumbnail,
	                   description,
	                   publish_date,
	                   price,
	                   created_at
	               ) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id
	`

	tx, err := database.DB.BeginTx(c, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	courseId := utils.GenUUID()
	result, err := tx.ExecContext(
		c,
		createCourseSql,
		courseId,
		createCoursePayload.Title,
		utils.Slugify(createCoursePayload.Title),
		createCoursePayload.Thumbnail,
		createCoursePayload.Description,
		nil,
		createCoursePayload.Price,
		time.Now(),
	)
	if err != nil {
		_ = tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*** Create course user mapping ***/
	result, err = tx.ExecContext(c, `
			insert into authors_courses(
								course_id,
								author_id
								)
					values ($1, $2) returning id
					`,
		courseId,
		authUser.UserId,
	)

	if err != nil {
		// rollback previous step
		_ = tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*** Create course category mapping ***/

	if len(createCoursePayload.Categories) > 0 {
		for _, category := range createCoursePayload.Categories {
			result, err = tx.ExecContext(c, `
			insert into courses_categories(course_id, category_id) values ($1, $2)`,
				courseId,
				category,
			)
			if err != nil {
				// rollback previous step
				_ = tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if len(createCoursePayload.SubCategories) > 0 {
		for _, subCategory := range createCoursePayload.SubCategories {
			result, err = tx.ExecContext(c, `
			insert into courses_sub_categories(course_id, category_id )values ($1, $2)`,
				courseId,
				subCategory,
			)
			if err != nil {
				// rollback previous step
				_ = tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if len(createCoursePayload.Topics) > 0 {
		for _, topic := range createCoursePayload.Topics {
			result, err = tx.ExecContext(c, `
			insert into courses_topics(course_id, topic_id) values ($1, $2)`,
				courseId,
				topic,
			)
			if err != nil {
				// rollback previous step
				_ = tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	_, err = result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"data": createCoursePayload,
	})
}

func SearchCourse(c *gin.Context) {
	type SearchCoursePayload struct {
		Value string `json:"value"`
	}

	// Check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		return
	}

	var searchCoursePayload SearchCoursePayload
	err := c.ShouldBindJSON(&searchCoursePayload)
	if err != nil {
		response.ErrorResponse(c, errors.New("Body data missing."), nil)
		return
	}

	// Start transaction
	tx, err := database.DB.BeginTx(c, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var keywordID int
	err = tx.QueryRowContext(
		c,
		"SELECT keyword_id FROM keywords WHERE keyword = $1",
		searchCoursePayload.Value,
	).Scan(&keywordID)

	if err == sql.ErrNoRows {
		// Keyword does not exist, insert it
		err = tx.QueryRowContext(
			c,
			"INSERT INTO keywords(keyword, type_id) VALUES ($1, $2) RETURNING keyword_id",
			searchCoursePayload.Value,
			1,
		).Scan(&keywordID)

		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if err != nil {
		_ = tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var currentRank, currentPreferenceScore float64
	err = tx.QueryRowContext(
		c,
		`SELECT rank, preference_score FROM customer_keyword_metadata 
		 WHERE user_id = $1 AND keyword_id = $2 AND type_id = $3`,
		authUser.UserId,
		keywordID,
		1,
	).Scan(&currentRank, &currentPreferenceScore)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate new values
	newRank := currentRank + 1
	newPreferenceScore := currentPreferenceScore + 1

	if err == sql.ErrNoRows {
		// Insert new record
		_, err = tx.ExecContext(
			c,
			`INSERT INTO customer_keyword_metadata
					(user_id, keyword_id, type_id, rank, preference_score, created_at) 
				VALUES ($1, $2, $3, $4, $5, $6)`,
			authUser.UserId,
			keywordID,
			1,
			newRank,            // New calculated rank
			newPreferenceScore, // New calculated preference score
			time.Now(),         // Created at
		)
	} else {
		// Update existing record
		_, err = tx.ExecContext(
			c,
			`UPDATE customer_keyword_metadata 
				SET rank = $4, preference_score = $5, created_at = $6 
				WHERE user_id = $1 AND keyword_id = $2 AND type_id = $3`,
			authUser.UserId,
			keywordID,
			1,
			newRank,            // New calculated rank
			newPreferenceScore, // New calculated preference score
			time.Now(),         // Updated at
		)
	}

	if err != nil {
		_ = tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": "Keyword and user preference successfully recorded",
	})
}

func GetUserSearchSuggestions(c *gin.Context) {
	type SearchSuggestionPayload struct {
		Query string `json:"query"`
	}

	var payload SearchSuggestionPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		response.ErrorResponse(c, errors.New("Body data missing."), nil)
		return
	}

	// Check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorized"), nil)
		return
	}

	// Fetch user-specific suggestions from the database
	query := `%` + payload.Query + `%`
	rows, err := database.DB.QueryContext(
		c,
		`SELECT k.keyword FROM keywords k  where k.keyword LIKE $1`,
		query,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		suggestions = append(suggestions, keyword)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
	})
}

func ClearUserSearch(c *gin.Context) {
	type ClearSearchPayload struct {
		Keyword string `json:"keyword"`
	}

	var payload ClearSearchPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		response.ErrorResponse(c, errors.New("Body data missing."), nil)
		return
	}

	// Check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorized"), nil)
		return
	}

	// Update the is_cleared field in customer_keyword_metadata
	_, err = database.DB.ExecContext(
		c,
		`UPDATE customer_keyword_metadata ckm
		 SET is_cleared = TRUE
		 FROM keywords k
		 WHERE ckm.keyword_id = k.keyword_id AND ckm.user_id = $1 AND k.keyword = $2`,
		authUser.UserId,
		payload.Keyword,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Search cleared successfully",
	})
}
