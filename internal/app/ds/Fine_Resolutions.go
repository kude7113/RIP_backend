package ds

type Fine_Resolutions struct {
	Fin_Res_ID    int `gorm:"primaryKey"`
	Resolution_ID int `gorm:"foreignKey:Resolution_ID"`
	Fine_ID       int `gorm:"foreignKey:Fine_ID"`
	Number        int `gorm:"default:1"`
}
