package models

import (
	"RIP/internal/app/ds"
	"mime/multipart"
	"time"
)

type FinesListWithRes struct {
	Fines    *[]ds.Fines `json:"fines"`
	ResCount int         `json:"resCount"`
	ResID    int         `json:"resId"`
}

type FineWithCount struct {
	Fine       *ds.Fines `json:"fines"`
	Count      int       `json:"count"`
	Fin_res_ID int       `json:"fin_res_id"`
}

type ResForAll struct {
	Resolution_ID     int
	Status            string
	Date_Created      time.Time
	Date_Formed       time.Time
	Date_Done         time.Time
	Car_License_Plate string
	User              string
	Head_Of_Depart    string
}

type ResWithFines struct {
	Res   *ds.Resolutions
	Fines *[]FineWithCount
}

type UploadImageRequest struct {
	Image *multipart.FileHeader // Поле для файла
}
