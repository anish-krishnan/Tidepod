package entity

import "time"

// Photo contains information about a single photos
type Photo struct {
	ID            int `json:"id" binding:"required"`
	FilePath      string
	CameraModel   string
	Latitude      float64
	Longitude     float64
	Timestamp     time.Time
	FocalLength   float64
	ApertureFStop float64
}
