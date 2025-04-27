package model

type User struct {
	Id   int64  `json:"id"   gorm:"primaryKey;autoIncrement;column:id;type:bigint(20);not null"`
	Name string `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Age  int    `json:"age"  gorm:"column:age;type:int(4);not null"`
}

func (t *User) TableName() string {
	return "t_user"
}
