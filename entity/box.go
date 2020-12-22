package entity

type Box struct {
	ID      int `gorm:"primary_key" yaml:"-"`
	PhotoID int
	FaceID  int
	MinX    int
	MinY    int
	MaxX    int
	MaxY    int
	Photo   Photo
	Face    Face
}
