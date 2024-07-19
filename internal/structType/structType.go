package structType

type CreateCoursePayload struct {
	Title         string  `json:"title"`
	Description   *string `json:"description"`
	Thumbnail     string  `json:"thumbnail"`
	Price         float64 `json:"price"`
	CategoryId    uint64  `json:"categoryId"`
	TopicId       *uint64 `json:"topicId"`
	SubCategoryId *uint64 `json:"subCategoryId"`
}
