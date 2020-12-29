package entity

// Box is a rectangular region within a photo containing a Face.
// A box belongs to exactly one photo, and one person's face
// One photo can contain many boxes (faces), and a person/face
// can contain many boxes
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
