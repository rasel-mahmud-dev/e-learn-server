package roleHandler

import (
	"database/sql"
	"e-learn/internal/database"
	"e-learn/internal/models/role"
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetRoles(c *gin.Context) {

	// check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("Unauthorization"), nil)
		return
	}

	columns := []string{
		"id",
		"name",
		"slug",
		"role_id",
		"status",
		"description",
		"created_at",
		"deleted_at",
	}

	roles, err := role.GetAllBySelect(c, columns, func(rows *sql.Rows, role *role.Role) error {
		return rows.Scan(
			&role.ID,
			&role.Name,
			&role.Slug,
			&role.RoleId,
			&role.Status,
			&role.Description,
			&role.CreatedAt,
			&role.DeletedAt,
		)
	}, "")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": roles,
	})
}

func GetUsersRoles(c *gin.Context) {

	// check auth
	//authUser := utils.GetAuthUser(c)
	//if authUser == nil {
	//	response.ErrorResponse(c, errors.New("Unauthorization"), nil)
	//	return
	//}

	usersRoles, err := role.GetUsersRolesMapping(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": usersRoles,
	})
}

func CreateRole(c *gin.Context) {

	// check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("Unauthorization"), nil)
		return
	}

	var createRolePayload role.Role

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&createRolePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createRoleSql := `insert into roles(
                    role_id,
                    name,
                    slug,
                    description, 
                    status, 
                    created_at
                )
		values ($1, $2, $3, $4, $5, $6) 
`

	result, err := database.DB.ExecContext(
		c,
		createRoleSql,
		utils.GenUUID(),
		createRolePayload.Name,
		utils.Slugify(createRolePayload.Name),
		createRolePayload.Description,
		"active",
		time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(result.RowsAffected())

	c.JSON(http.StatusCreated, gin.H{
		"data": createRolePayload,
	})
}
