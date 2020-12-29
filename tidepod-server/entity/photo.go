package entity

import (
	"time"
)

// Photo contains information about a single photo
type Photo struct {
	ID               int `gorm:"primary_key" yaml:"-"`
	FilePath         string
	CameraModel      string
	Latitude         float64
	Longitude        float64
	LocationString   string
	Timestamp        time.Time
	FocalLength      float64
	ApertureFStop    float64
	OriginalFilename string
	Labels           []Label `gorm:"many2many:photo_labels;"`
	Boxes            []Box
}
