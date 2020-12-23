package entity

type Box struct {
	ID       int `gorm:"primary_key" yaml:"-"`
	PhotoID  int
	FaceID   int
	FilePath string
	MinX     int
	MinY     int
	MaxX     int
	MaxY     int
	Photo    Photo
	Face     Face
}
