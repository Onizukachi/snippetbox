package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Onizukachi/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	user := models.User{}
	stm := "SELECT id, name, email, hashed_password, created, active FROM users where email = ? and active = true"
	err := m.DB.QueryRow(stm, email).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Created, &user.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return 0, models.ErrInvalidCredentials
	}

	return user.ID, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
