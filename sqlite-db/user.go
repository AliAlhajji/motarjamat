package sqlitedb

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/AliAlhajji/Motarjamat/utils"
)

type userRepo struct {
	db *sqlx.DB
}

func InitUserRepo() (*userRepo, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	return &userRepo{
		db: db,
	}, nil
}

func (r *userRepo) CreateUser(uuid string, username string, email string, name string) (int64, error) {
	q := `INSERT INTO user(uuid, username, email, name, join_date, role) VALUES (
		:uuid, :username, :email, :name, :join_date, :role
	)`

	user := models.User{
		UUID:     uuid,
		Username: username,
		Email:    email,
		Name:     name,
		JoinDate: time.Now(),
		Role:     models.RoleUser,
	}

	_, err := r.db.NamedExec(q, &user)

	if err != nil {
		return 0, err
	}

	return 0, nil
}

// Verify the provided credentials are correct
func (r *userRepo) VerifyCredentials(usernameOrEmail string, password string) (string, error) {
	//Get the current password of the account
	q := `SELECT uuid, password FROM user WHERE username=$1 OR email=$2`

	var user models.User

	err := r.db.Get(&user, q, usernameOrEmail, usernameOrEmail)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	// Verify the provided password is the same as the one in the database
	passwordOK := utils.VerifyPassword(user.Password, password)
	if passwordOK {
		return user.UUID, nil
	}

	return "", nil
}

func (r *userRepo) GetUserByUsername(username string) (*models.User, error) {
	return nil, nil

}

func (r *userRepo) GetUserByID(id int64) (*models.User, error) {
	return nil, nil

}

func (r *userRepo) GetUserByUUID(uuid string) (*models.User, error) {
	q := `SELECT * FROM user WHERE uuid=$1`
	var user models.User

	err := r.db.Get(&user, q, uuid)
	if err != nil {
		return nil, err
	}

	return &user, nil

}
