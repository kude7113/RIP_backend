package ds

import (
	"time"
)

type Resolutions struct {
	Resolution_ID     int    `gorm:"primaryKey"`
	Status            string `gorm:"type:varchar(255)"`
	Date_Created      time.Time
	Car_License_Plate string `gorm:"type:varchar(255)"`
	User_ID           int    `gorm:"foreignKey:User_ID"`
}

const (
	DraftStatus     = "черновик"
	DeletedStatus   = "удален"
	FormedStatus    = "сформирован"
	CompletedStatus = "завершен"
	RejectedStatus  = "отклонен"
	ApprovedStatus  = "подтвержден"
)
