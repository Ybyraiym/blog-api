package authentication

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// PostgreSQLAuthenticator - структура для аутентификации в базе данных PostgreSQL.
type PostgreSQLAuthenticator struct {
	DB *sql.DB
}

// TokenSecret - секретный ключ для подписи токена
// (Надо потом доделать ключ)
var TokenSecret = []byte("скретный_ключ")

// Authenticate проверяет логин и пароль в базе данных PostgreSQL.
func (a *PostgreSQLAuthenticator) Authenticate(username, password string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE login = $1 AND password = $2"

	var count int
	err := a.DB.QueryRow(query, username, password).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

// GenerateToken создает JWT-токен с именем пользователя
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // токен действителен 24 часа

	tokenString, err := token.SignedString(TokenSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Вспомогательная функция для получения имени пользователя из токена
func GetUsernameFromToken(tokenString string) (string, error) {
	// Проверяем, что токен не пустой
	if tokenString == "" {
		return "", errors.New("токен отсутствует")
	}

	// Проверяем токен и извлекаем из него имя пользователя
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неподдерживаемый метод подписи: %v", token.Header["alg"])
		}

		// Возвращаем секретный ключ
		return TokenSecret, nil
	})
	if err != nil {
		return "", err
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, ok := claims["username"].(string); ok {
			return username, nil
		}
	}

	return "", errors.New("неверный токен")
}
