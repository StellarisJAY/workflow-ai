package model

type User struct {
	UserId   int64  `json:"id" gorm:"primary_key;column:user_id;type:bigint"`
	Username string `json:"name" gorm:"column:username;type:varchar(50);not null"`
	Password string `json:"password" gorm:"column:password;type:varchar(255);not null"`
}

func (User) TableName() string {
	return "wf_user"
}
