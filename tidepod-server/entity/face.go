package entity

// Face contains information about a Face (ie. a person)
type Face struct {
	ID    int `gorm:"primary_key" yaml:"-"`
	Name  string
	Boxes []Box
}
