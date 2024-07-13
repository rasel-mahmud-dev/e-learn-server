package models

import (
	"database/sql"
	"e-learn/internal/database"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Profile struct {
	ID        uint64     `json:"id,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	Headline  *string    `json:"headline,omitempty"`
	Language  *string    `json:"language,omitempty"`
	Website   *string    `json:"website,omitempty"`
	Twitter   *string    `json:"twitter,omitempty"`
	Facebook  *string    `json:"facebook,omitempty"`
	YouTube   *string    `json:"youtube,omitempty"`
	Github    *string    `json:"github,omitempty"`
	AboutMe   *string    `json:"about_me,omitempty"`
	UserId    uint64     `json:"user_id,omitempty"`
}

type ProfileWithCamelCaseJSON struct {
	Profile
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	FirstName *string    `json:"firstName,omitempty"`
	LastName  *string    `json:"lastName,omitempty"`
	AboutMe   *string    `json:"aboutMe,omitempty"`
	UserId    uint64     `json:"userId,omitempty"`
}

func GetProfileById(c *gin.Context, columns []string, scanFunc func(*sql.Row, *Profile) error, userId uint64) (*Profile, error) {

	query := fmt.Sprintf("SELECT %s FROM profiles WHERE user_id = $1", strings.Join(columns, ", "))

	// Prepare a variable to hold the users data
	user := &Profile{}

	// Execute the query
	row := database.DB.QueryRowContext(c, query, userId)

	//Use the scan function to populate the users struct
	err := scanFunc(row, user)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where no rows were found
			fmt.Println("No users found with the given ID")
			return nil, errors.New("no users found")
		}
		return nil, err
	}

	return user, nil
}

func UpdateProfile(c *gin.Context, request *Profile) (*Profile, error) {

	var profileId uint64
	err := database.DB.QueryRow("SELECT id FROM profiles WHERE user_id = $1", request.UserId).Scan(&profileId)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var updates []string
	var values []interface{}
	var insertFields []string
	var insertPlaceholders []string
	valueIndex := 1

	// Check which fields to update and build the query dynamically
	if request.DeletedAt != nil {
		updates = append(updates, fmt.Sprintf("deleted_at = $%d", valueIndex))
		insertFields = append(insertFields, "deleted_at")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.DeletedAt)
		valueIndex++
	}
	if request.CreatedAt != nil {
		updates = append(updates, fmt.Sprintf("created_at = $%d", valueIndex))
		insertFields = append(insertFields, "created_at")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.CreatedAt)
		valueIndex++
	}
	if request.UpdatedAt != nil {
		updates = append(updates, fmt.Sprintf("updated_at = $%d", valueIndex))
		insertFields = append(insertFields, "updated_at")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.UpdatedAt)
		valueIndex++
	}
	if request.FirstName != nil {
		updates = append(updates, fmt.Sprintf("first_name = $%d", valueIndex))
		insertFields = append(insertFields, "first_name")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.FirstName)
		valueIndex++
	}
	if request.LastName != nil {
		updates = append(updates, fmt.Sprintf("last_name = $%d", valueIndex))
		insertFields = append(insertFields, "last_name")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.LastName)
		valueIndex++
	}
	if request.Headline != nil {
		updates = append(updates, fmt.Sprintf("headline = $%d", valueIndex))
		insertFields = append(insertFields, "headline")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.Headline)
		valueIndex++
	}
	if request.Language != nil {
		updates = append(updates, fmt.Sprintf("language = $%d", valueIndex))
		insertFields = append(insertFields, "language")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.Language)
		valueIndex++
	}
	if request.Website != nil {
		updates = append(updates, fmt.Sprintf("website = $%d", valueIndex))
		insertFields = append(insertFields, "website")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.Website)
		valueIndex++
	}
	if request.Twitter != nil {
		updates = append(updates, fmt.Sprintf("twitter = $%d", valueIndex))
		insertFields = append(insertFields, "twitter")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.Twitter)
		valueIndex++
	}
	if request.Facebook != nil {
		updates = append(updates, fmt.Sprintf("facebook = $%d", valueIndex))
		insertFields = append(insertFields, "facebook")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.Facebook)
		valueIndex++
	}
	if request.YouTube != nil {
		updates = append(updates, fmt.Sprintf("youtube = $%d", valueIndex))
		insertFields = append(insertFields, "youtube")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.YouTube)
		valueIndex++
	}
	if request.Github != nil {
		updates = append(updates, fmt.Sprintf("github = $%d", valueIndex))
		insertFields = append(insertFields, "github")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.Github)
		valueIndex++
	}
	if request.AboutMe != nil {
		updates = append(updates, fmt.Sprintf("about_me = $%d", valueIndex))
		insertFields = append(insertFields, "about_me")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.AboutMe)
		valueIndex++
	}
	if request.UserId != 0 {
		updates = append(updates, fmt.Sprintf("user_id = $%d", valueIndex))
		insertFields = append(insertFields, "user_id")
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", valueIndex))
		values = append(values, request.UserId)
		valueIndex++
	}

	// If the profile exists, update it
	if profileId != 0 {
		if len(updates) == 0 {
			return nil, fmt.Errorf("no fields to update")
		}

		// Build the SQL update statement
		query := fmt.Sprintf("UPDATE profiles SET %s WHERE id = $%d", strings.Join(updates, ", "), valueIndex)

		values = append(values, profileId) // Update the existing profile
		_, err = database.DB.ExecContext(c, query, values...)
		if err != nil {
			return nil, err
		}
	} else {
		// Build the SQL insert statement
		query := fmt.Sprintf(
			"INSERT INTO profiles (%s) VALUES (%s) RETURNING id",
			strings.Join(insertFields, ", "),
			strings.Join(insertPlaceholders, ", "),
		)

		values = append([]interface{}{}, values...)

		fmt.Println(values)
		err = database.DB.QueryRowContext(c, query, values...).Scan(&profileId)
		if err != nil {
			return nil, err
		}
		//request.ID = profileId
	}

	return request, nil

}
