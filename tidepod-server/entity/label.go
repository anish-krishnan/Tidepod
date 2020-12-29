package entity

// Label is used for photo categorization
type Label struct {
	ID        int `gorm:"primary_key" yaml:"-"`
	LabelName string
	Photos    []Photo `gorm:"many2many:photo_labels;"`
}
