package repository

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
	"time"
)

const ExpTime = time.Hour * 1

func GenerateJWTToken(userID int, isAdmin bool) (string, error) {
	// Создаем новый JWT токен с алгоритмом подписи HS256
	token := jwt.New(jwt.SigningMethodHS256)

	// Получаем доступ к claims
	claims := token.Claims.(jwt.MapClaims)

	// Добавляем данные пользователя
	claims["user_id"] = userID
	claims["is_admin"] = isAdmin

	// Добавляем время истечения токена
	claims["exp"] = ExpTime

	// Добавляем время создания токена (текущую метку времени)
	claims["iat"] = time.Now().Unix() // Время создания в формате Unix timestamp

	// Генерируем подпись для токена
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	// Возвращаем токен
	return tokenString, nil
}

// CheckJWTInBlacklist - проверка токена в черном списке
func (r *Repository) CheckJWTInBlacklist(userIDStr string) (string, error) {
	// Получаем токен из Redis по ключу userID
	res, err := r.rd.Get(userIDStr).Result()
	// Если токен не найден (ключ не существует) или ошибка, возвращаем пустую строку
	if err != nil {
		return "", err // Если не нашли токен в Redis, возвращаем ошибку
	}
	return res, nil
}

// IsBlackListed - проверка на наличие токена в черном списке
func (r *Repository) IsBlackListed(userID float64, token string) bool {
	userIDStr := strconv.FormatUint(uint64(userID), 10) // Преобразуем user_id в строку
	res, err := r.CheckJWTInBlacklist(userIDStr)        // Проверяем по userID
	if err != nil || res == "" {
		return false // Токен не найден в черном списке или произошла ошибка
	}

	// Если токен в черном списке, возвращаем true
	if res == token {
		return true
	}
	return false // Токен не найден в черном списке
}

// SaveJWTToken - сохранение токена в БД
func (r *Repository) SaveJWTToken(userID int, token string) error {
	// Парсим токен для извлечения времени истечения
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return err // Если не удается распарсить токен, возвращаем ошибку
	}

	// Извлекаем данные из токена
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Получаем время истечения токена (exp)
	expTime, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("exp not found in token claims")
	}

	// Рассчитываем TTL (Time To Live) для Redis
	ttl := time.Unix(int64(expTime), 0).Sub(time.Now()) // Разница между текущим временем и временем истечения

	// Конвертируем userID в строку
	idStr := strconv.FormatUint(uint64(userID), 10)

	// Сохраняем токен в Redis с TTL
	err = r.rd.Set(idStr, token, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
