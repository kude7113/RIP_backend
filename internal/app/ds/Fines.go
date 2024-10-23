package ds

type Fines struct {
	Fine_ID int    `gorm:"primaryKey;autoIncrement" json:"fineID"`
	Title   string `gorm:"type:varchar(255)" json:"title"`
	FullInf string `gorm:"type:varchar(255)" json:"fullInf"`
	Price   int    `gorm:"type:int" json:"price"`
	Imge    string `gorm:"type:varchar(255)" json:"imge"`
	DopInf  string `gorm:"type:varchar(255)" json:"dopInf"`
}
