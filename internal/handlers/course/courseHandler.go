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
	"github.com/gin-gonic/gin"
	"net/http"
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
		"(select title from courses_categories c where courses.course_id = c.course_id) as categories",
		"(select title from courses_sub_categories sc where courses.course_id = sc.course_id) as sub_categories",
		"(select title from courses_topics ct where courses.course_id = ct.course_id) as topics",
		"(select id from authors_courses ac where ac.course_id = courses.course_id) as authors",
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
			&course.CategoryList,
			&course.SubCategoryList,
			&course.TopicList,
			&course.AuthorList,
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
