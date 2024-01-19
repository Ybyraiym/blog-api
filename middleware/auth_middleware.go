// package middleware

// import (
// 	"context"
// 	"errors"
// 	"net/http"
// 	"strings"

// 	"blog-api/authentication"

// 	"github.com/sirupsen/logrus"
// )

// var log = logrus.New()

// // AuthMiddleware - middleware для проверки токена аутентификации
// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("AuthMiddleware: Checking authentication")

// 		// Получаем токен из запроса
// 		token, err := getTokenFromRequest(r)
// 		// if err != nil {
// 		// 	http.Error(w, "Ошибка получения токена", http.StatusUnauthorized)
// 		// 	return
// 		// }
// 		if err != nil {
// 			// Пропускаем проверку, если токен отсутствует
// 			log.Info("AuthMiddleware: Token not found, skipping authentication check")
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// Если токен получен успешно, извлекаем информацию о пользователе и добавляем токен обратно в запрос
// 		if username, err := authentication.GetUsernameFromToken(token); err == nil {
// 			r.Header.Set("Authorization", "Bearer "+token)
// 			r = r.WithContext(context.WithValue(r.Context(), "username", username))
// 			log.Info("AuthMiddleware: Authentication successful")
// 		} else {
// 			log.Error("AuthMiddleware: Authentication failed -", err.Error())
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// // Вспомогательная функция для извлечения токена из запроса
// func getTokenFromRequest(r *http.Request) (string, error) {
// 	// Получаем токен из заголовка Authorization
// 	tokenString := r.Header.Get("Authorization")
// 	if tokenString == "" {
// 		return "", errors.New("токен отсутствует в запросе")
// 	}

// 	// Возвращаем токен без префикса "Bearer "
// 	return strings.TrimPrefix(tokenString, "Bearer "), nil
// }

package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"blog-api/authentication"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// AuthMiddleware проверяет токен аутентификации в запросе и предоставляет доступ к группе эндпоинтов /api.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Checking authentication")

		// Получаем токен из запроса
		token, err := getTokenFromRequest(r)
		if err != nil {
			log.Info("AuthMiddleware: Token not found, skipping authentication check")
			// Если токен отсутствует, пропускаем проверку и передаем управление следующему обработчику
			next.ServeHTTP(w, r)
			return
		}

		// Если токен получен успешно, извлекаем информацию о пользователе и добавляем токен обратно в запрос
		if username, err := authentication.GetUsernameFromToken(token); err == nil {
			// Проверка наличия правильного префикса в токене
			if strings.HasPrefix(r.URL.Path, "/api") {
				// Добавляем токен и имя пользователя в контекст запроса
				r.Header.Set("Authorization", "Bearer "+token)
				r = r.WithContext(context.WithValue(r.Context(), "username", username))
				log.Info("AuthMiddleware: Authentication successful")
			}
		} else {
			log.Error("AuthMiddleware: Authentication failed -", err.Error())
		}

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// Вспомогательная функция для извлечения токена из запроса
func getTokenFromRequest(r *http.Request) (string, error) {
	// Получаем токен из заголовка Authorization
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("токен отсутствует в запросе")
	}

	// Возвращаем токен без префикса "Bearer "
	return strings.TrimPrefix(tokenString, "Bearer "), nil
}
