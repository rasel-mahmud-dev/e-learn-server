package courseHandler

import (
	"e-learn/internal/database"
	"e-learn/internal/response"
	"e-learn/internal/structType"
	"e-learn/internal/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreateCourse(c *gin.Context) {

	// check auth

	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("Unauthorization"), nil)
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
                    )
		values ($1, $2, $3, $4, $5, $6, $7, $8) returning id
`

	courseId := utils.GenUUID()
	result, err := database.DB.ExecContext(
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*** Create course user mapping ***/
	result, err = database.DB.ExecContext(c, `
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
		// roll back previous step
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*** Create course category mapping ***/

	result, err = database.DB.ExecContext(c, `
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
		// roll back previous step
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(result.RowsAffected())

	c.JSON(http.StatusCreated, gin.H{
		"data": createCoursePayload,
	})
}
