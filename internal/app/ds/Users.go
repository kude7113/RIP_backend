package ds

type Users struct {
	User_ID  int    `json:"id" gorm:"primaryKey"`
	Login    string `json:"login" gorm:"type:varchar(255)"`
	Password string `json:"password" gorm:"type:varchar(255)"`
	IsAdmin  bool   `json:"is_admin" gorm:"default:false"`
}
