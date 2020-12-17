package entity

// Label contains information about a single label
type Label struct {
	ID        int `gorm:"primary_key" yaml:"-"`
	LabelName string
}
