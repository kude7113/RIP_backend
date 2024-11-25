package ds

import (
	"time"
)

type Resolutions struct {
	Resolution_ID     int    `gorm:"primaryKey"`
	User_ID           int    `gorm:"foreignKey:User_ID"`
	Head_Of_Depart_ID int    `gorm:"foreignKey:User_ID"`
	Status            string `gorm:"type:varchar(255)"`
	Date_Created      time.Time
	Date_Formed       time.Time
	Date_Done         time.Time
	Car_License_Plate string `gorm:"type:varchar(255)"`
	Total_Price       int    `gorm:"type:int"`
	Sale              bool   `gorm:"type:bool"`
}

const (
	DraftStatus     = "черновик"
	DeletedStatus   = "удален"
	FormedStatus    = "сформирован"
	CompletedStatus = "завершен"
	RejectedStatus  = "отклонен"
	ApprovedStatus  = "подтвержден"
)
