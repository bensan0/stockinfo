package models

//業種
type Industry struct {
	Id      uint `gorm:"primarykey"`
	Comment string
}

func (i Industry) TableName() string {
	return "industry"
}


