package repository

import (
	"RIP/internal/app/ds"
	"errors"
	"gorm.io/gorm"
	"strings"
)

func (r *Repository) GetAllFines() (*[]ds.Fines, error) {
	var deliveryItems []ds.Fines
	r.db.Model(&ds.Fines{}).Find(&deliveryItems)
	return &deliveryItems, nil
}

func (r *Repository) SearchFines(searchText string) (*[]ds.Fines, error) {
	searchText = strings.ToLower(searchText)
	var Fines []ds.Fines
	// сохраняем данные из бд в массив
	r.db.Find(&Fines)

	var filteredFines []ds.Fines
	for _, Fine := range Fines {
		fineTitle := strings.TrimSpace(strings.ToLower(Fine.Title))
		if strings.HasPrefix(fineTitle, searchText) {
			filteredFines = append(filteredFines, Fine)
		}
	}
	return &filteredFines, nil
}

func (r *Repository) GetResolutionLength(userId int) (int, error) {
	var count int64
	var req ds.Resolutions
	status := ds.DraftStatus

	if err := r.db.Where("User_id = ? AND Status = ?", userId, status).First(&req).Error; err != nil {
		return 0, err
	}

	reqID := req.Resolution_ID

	err := r.db.Model(&ds.Fine_Resolutions{}).Where("Resolution_id = ?", reqID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *Repository) ResolutionByUserID(userID int) (int, error) {
	var res ds.Resolutions
	err := r.db.Where("User_id = ? AND Status = ?", userID, ds.DraftStatus).First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return res.Resolution_ID, nil
}

func (r *Repository) GetFinesByID(id int) (*ds.Fines, error) {
	Fine := &ds.Fines{}

	err := r.db.First(Fine, "Fine_ID = ?", id).Error // find Fines with id = 1
	if err != nil {
		return nil, err
	}

	return Fine, nil
}

func (r *Repository) CreateFine(fines *ds.Fines) (*ds.Fines, error) {
	
}
