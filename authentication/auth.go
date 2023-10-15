package authentication

import (
	"database/sql"
)

// PostgreSQLAuthenticator - структура для аутентификации в базе данных PostgreSQL.
type PostgreSQLAuthenticator struct {
	DB *sql.DB
}

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
