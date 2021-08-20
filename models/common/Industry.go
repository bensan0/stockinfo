package models

//業種
type Industry struct {
	Id      uint `gorm:"unique;autoIncrement" json:"id"`
	Comment string `gorm:"primarykey" json:"comment"`
}

func (i Industry) TableName() string {
	return "industry"
}


