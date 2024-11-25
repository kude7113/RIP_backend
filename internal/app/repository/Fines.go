package repository

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/models"
	"RIP/internal/app/storage"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
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

func (r *Repository) ResolutionsList(dateFrom, dateTo *time.Time, status string) (*[]models.ResForAll, error) {
	var resolutions []models.ResForAll

	// Базовый запрос: выбираем все записи из таблицы `resolutions` и присоединяем таблицу `users`.
	query := r.db.Table("resolutions").
		Select(`resolutions.resolution_id, resolutions.status, resolutions.date_created, resolutions.date_formed, 
                resolutions.date_done, resolutions.car_license_plate, 
                users_login.login as user, head_login.login as head_of_depart`).
		Joins("LEFT JOIN users AS users_login ON users_login.user_id = resolutions.user_id").
		Joins("LEFT JOIN users AS head_login ON head_login.user_id = resolutions.head_of_depart_id")

	// Добавляем фильтрацию по `date_from`, если параметр не равен nil
	if dateFrom != nil {
		query = query.Where("resolutions.date_created >= ?", *dateFrom)
	}

	// Добавляем фильтрацию по `date_to`, если параметр не равен nil
	if dateTo != nil {
		query = query.Where("resolutions.date_created <= ?", *dateTo)
	}

	// Добавляем фильтрацию по `status`, если параметр не пуст
	if status != "" {
		query = query.Where("resolutions.status = ?", status)
	}

	// Выполняем запрос
	if err := query.Find(&resolutions).Error; err != nil {
		return nil, err
	}

	return &resolutions, nil
}

func (r *Repository) GetResByID(resID int) (*models.ResWithFines, error) {
	var res ds.Resolutions

	if err := r.db.Model(&ds.Resolutions{}).Where("resolution_id = ?", resID).First(&res).Error; err != nil {
		r.logger.Infof("Found penis")
		return nil, err
	}
	var finesWithCount []models.FineWithCount
	var finesResolution []ds.Fine_Resolutions
	// Получаем все записи Fine_Resolution для заданного resID
	err := r.db.Where("resolution_id = ?", resID).Find(&finesResolution).Error
	if err != nil {
		return nil, err
	}

	// Цикл по всем найденным штрафам в постановлении
	for _, fineRes := range finesResolution {
		var fine ds.Fines
		// Для каждого штрафа ищем его полную информацию
		err := r.db.Where("fine_id = ?", fineRes.Fine_ID).First(&fine).Error
		if err != nil {
			return nil, err
		}

		// Создаём объект FinesWithCount и добавляем его в результат
		fineCount := models.FineWithCount{
			Fine:  &fine,
			Count: fineRes.Number, // Используем поле Number из finesResolution
		}
		finesWithCount = append(finesWithCount, fineCount)
	}

	result := models.ResWithFines{
		Res:   &res,
		Fines: &finesWithCount,
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

	var totalPrice int
	var finesResolution []ds.Fine_Resolutions
	// Получаем все записи Fine_Resolution для заданного resID
	err = r.db.Where("resolution_id = ?", resID).Find(&finesResolution).Error
	if err != nil {
		return nil, err
	}

	// Цикл по всем найденным штрафам в постановлении
	for _, fineRes := range finesResolution {
		var fine ds.Fines
		// Для каждого штрафа ищем его полную информацию
		err := r.db.Where("fine_id = ?", fineRes.Fine_ID).First(&fine).Error
		if err != nil {
			return nil, err
		}

		totalPrice += fine.Price * fineRes.Number
	}

	if rand.Intn(2) == 0 {
		// С вероятностью 50% умножаем число на 0.7
		totalPrice *= int(float64(totalPrice) * 0.7)
		result.Sale = true
	}
	result.Total_Price = totalPrice
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

func (r *Repository) UpdateFRCount(fid int, newNumber int) (*ds.Fine_Resolutions, error) {
	var result ds.Fine_Resolutions

	// Ищем запись с указанным ID
	err := r.db.Where("fin_res_id = ?", fid).First(&result).Error
	if err != nil {
		return nil, err
	}

	// Обновляем только поле Number
	result.Number = newNumber
	if err = r.db.Save(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Repository) UploadImageAndUpdateURL(fineID int, fileName string, file multipart.File, fileSize int64) (string, error) {
	// Инициализация Minio хранилища
	minioStorage, err := storage.NewMinioStorage(
		os.Getenv("MINIO_ENDPOINT_URL"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_SECURE") == "true",
	)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Minio client: %w", err)
	}

	// Получение существующей записи о штрафе
	delivery, err := r.GetFinesByID(fineID)
	if err != nil {
		return "", fmt.Errorf("failed to get delivery by ID: %w", err)
	}
	if delivery == nil {
		return "", fmt.Errorf("delivery not found")
	}

	// Удаление предыдущего изображения из Minio, если оно существует
	if delivery.Imge != "" {
		previousFileName := filepath.Base(delivery.Imge) // Получаем имя файла из URL
		err = minioStorage.DeleteImg(os.Getenv("MINIO_BUCKET_NAME"), previousFileName)
		if err != nil {
			return "", fmt.Errorf("failed to delete previous image: %w", err)
		}
	}

	// Загрузка нового файла в Minio
	err = minioStorage.LoadImg(os.Getenv("MINIO_BUCKET_NAME"), fileName, file, fileSize)
	if err != nil {
		return "", fmt.Errorf("failed to load image to Minio: %w", err)
	}

	// Генерация URL нового изображения
	imageURL := "http://" + os.Getenv("MINIO_ENDPOINT_URL") + "/" + os.Getenv("MINIO_BUCKET_NAME") + "/" + fileName

	// Обновление URL изображения в базе данных
	query := "UPDATE fines SET imge = $1 WHERE fine_id = $2"
	result := r.db.Exec(query, imageURL, fineID)
	if result.Error != nil {
		return "", fmt.Errorf("failed to update image URL in database: %w", result.Error)
	}

	r.logger.Info("Rows affected:", result.RowsAffected)

	return imageURL, nil
}

func (r *Repository) CreateUser(name, password string) (*ds.Users, error) {

	user := &ds.Users{
		Login:    name,
		Password: password,
	}
	// логин не может повторяться
	if err := r.db.Where("login = ?", name).First(&user).Error; err == nil {
		return nil, fmt.Errorf("user with login %s already exists", name)
	}
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpdateUser(user *ds.Users) (*ds.Users, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
