package repository

import (
	"RIP/internal/app/ds"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"strings"
	"time"
)

func (r *Repository) GetAllFines() (*[]ds.Fines, error) {
	var deliveryItems []ds.Fines
	err := r.db.Model(&ds.Fines{}).Find(&deliveryItems).Error
	if err != nil {
		return nil, err
	}
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
	result := r.db.Create(fines)
	if result.Error != nil {
		return nil, result.Error
	}
	return fines, nil
}

func (r *Repository) UpdateFine(fine *ds.Fines) (*ds.Fines, error) {
	validate := validator.New()
	if err := validate.Struct(fine); err != nil {
		return nil, err
	}

	// Обновляем все поля, кроме поля Image
	if err := r.db.Model(&ds.Fines{}).Omit("Image").Where("fine_id = ?", fine.Fine_ID).Updates(fine).Error; err != nil {
		return nil, err
	}

	return fine, nil
}

func (r *Repository) DeleteFine(id int) error {
	err := r.db.Delete(&ds.Fines{}, "Fine_ID = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddFinesToResolution(userID, fineID int) error {
	// Поиск существующей заявки пользователя со статусом 'черновик'
	var draftRequest ds.Resolutions
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&draftRequest).Error

	// Если черновик не найден, создаём новый
	if errors.Is(err, gorm.ErrRecordNotFound) {
		draftRequest = ds.Resolutions{
			User_ID:      userID,
			Status:       ds.DraftStatus,
			Date_Created: time.Now(),
		}

		// Создание новой записи
		err = r.db.Create(&draftRequest).Error
		if err != nil {
			return fmt.Errorf("error creating new draft request: %w", err)
		}

		r.logger.Infof("Created new draft request ID: %d for user ID: %d", draftRequest.Resolution_ID, userID)
	} else if err != nil {
		// Если произошла ошибка запроса, возвращаем её
		return fmt.Errorf("error fetching draft request: %w", err)
	} else {
		r.logger.Infof("Found existing draft request ID: %d for user ID: %d", draftRequest.Resolution_ID, userID)
	}

	// Добавляем элемент в существующую заявку (или новую)
	fineRes := ds.Fine_Resolutions{
		Fine_ID:       fineID,
		Resolution_ID: draftRequest.Resolution_ID,
	}

	// Вставляем в базу данных
	err = r.db.Create(&fineRes).Error
	if err != nil {
		return fmt.Errorf("error linking fine to draft request: %w", err)
	}

	r.logger.Infof("Fine ID: %d successfully added to Resolution ID: %d", fineID, draftRequest.Resolution_ID)

	return nil
}

func (r *Repository) ResolutionsList() (*[]ds.Resolutions, error) {
	var result []ds.Resolutions

	err := r.db.Model(&ds.Resolutions{}).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Repository) GetResByID(id int) (*ds.Resolutions, error) {
	var result ds.Resolutions

	err := r.db.First(&result, "resolution_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Repository) UpdateRes(resolution *ds.Resolutions) (*ds.Resolutions, error) {
	validate := validator.New()
	if err := validate.Struct(resolution); err != nil {
		return nil, err
	}

	if err := r.db.Model(&ds.Resolutions{}).Where("resolution_id = ?", resolution.Resolution_ID).Updates(resolution).Error; err != nil {
		return nil, err
	}

	return resolution, nil
}

func (r *Repository) SetStatusByUser(userID int) (*ds.Resolutions, error) {
	// проверко является ли user владельцем данного постановления
	var result ds.Resolutions
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&result).Error
	if err != nil {
		return nil, err
	}

	result.Status = ds.FormedStatus
	if err := r.db.Save(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Repository) SetStatusByAdmin(resID int, status string) (*ds.Resolutions, error) {
	var result ds.Resolutions

	err := r.db.Where("resolution_id  = ? AND status = ?", resID, ds.FormedStatus).First(&result).Error
	if err != nil {
		return nil, err
	}

	result.Status = status
	if err := r.db.Save(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) DeleteResolution(resID int) (*ds.Resolutions, error) {
	var result ds.Resolutions

	err := r.db.Where("resolution_id  = ?", resID).First(&result).Error
	if err != nil {
		return nil, err
	}

	result.Status = ds.DeletedStatus
	if err := r.db.Save(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) DeleteFR(id int) error {
	err := r.db.Delete(&ds.Fine_Resolutions{}, "fin_res_id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}
