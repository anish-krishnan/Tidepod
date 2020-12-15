package entity

// Photo contains information about a single photos
type Photo struct {
	ID       int `json:"id" binding:"required"`
	FilePath string
}
