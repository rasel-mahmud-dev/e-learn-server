package structType

type CreateCoursePayload struct {
	Title         string   `json:"title"`
	Description   *string  `json:"description"`
	Thumbnail     string   `json:"thumbnail"`
	Price         float64  `json:"price"`
	Categories    []uint64 `json:"categories"`
	Topics        []uint64 `json:"topics"`
	SubCategories []uint64 `json:"subCategories"`
}

type CreateReviewPayload struct {
	Title         string   `json:"title"`
	Description   *string  `json:"description"`
	Thumbnail     string   `json:"thumbnail"`
	Price         float64  `json:"price"`
	Categories    []uint64 `json:"categories"`
	Topics        []uint64 `json:"topics"`
	SubCategories []uint64 `json:"subCategories"`
}
