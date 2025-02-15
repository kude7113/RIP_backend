package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/skip2/go-qrcode"

	"RIP/internal/app/ds" // Замените на актуальный путь импорта вашего пакета ds
)

// getString возвращает значение строки, если оно не пустое, иначе — "не указан".
func getString(s string) string {
	if s == "" {
		return "не указан"
	}
	return s
}

// getInt возвращает строковое представление числа, если оно не равно 0, иначе — "не указан".
func getInt(n int) string {
	if n == 0 {
		return "не указан"
	}
	return fmt.Sprintf("%d", n)
}

// getTime возвращает отформатированное значение времени, если оно не нулевое, иначе — "не указан".
func getTime(t time.Time) string {
	if t.IsZero() {
		return "не указан"
	}
	// Формат: ГГГГ-ММ-ДД ЧЧ:ММ:СС
	return t.Format("2006-01-02 15:04:05")
}

// GenerateResolutionQR формирует строку со всеми полями структуры ds.Resolutions,
// подставляя "не указан" для пустых значений, генерирует QR‑код и возвращает его в формате base64.
func GenerateResolutionQR(r ds.Resolutions) (string, error) {
	info := fmt.Sprintf(
		"Постановление №: %s\n"+
			"Пользователь ID: %s\n"+
			"Руководитель отдела ID: %s\n"+
			"Статус: %s\n"+
			"Дата создания: %s\n"+
			"Дата формирования: %s\n"+
			"Дата выполнения: %s\n"+
			"Номер машины: %s\n"+
			"Итоговая цена: %s\n"+
			"Скидка: %t\n"+
			"QR: %s",
		getInt(r.Resolution_ID),
		getInt(r.User_ID),
		getInt(r.Head_Of_Depart_ID),
		getString(r.Status),
		getTime(r.Date_Created),
		getTime(r.Date_Formed),
		getTime(r.Date_Done),
		getString(r.Car_License_Plate),
		getInt(r.Total_Price),
		r.Sale,
		getString(r.Qr),
	)

	// Генерация QR‑кода с размером 256x256 и уровнем коррекции ошибок Medium.
	png, err := qrcode.Encode(info, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	// Преобразование PNG‑изображения в строку base64.
	base64Image := base64.StdEncoding.EncodeToString(png)
	return base64Image, nil
}
