package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jschaefer-io/IDaaS/utils"
)

type Gender string

const (
	GenderMale   Gender = "m"
	GenderFemale Gender = "f"
	GenderOther  Gender = "d"
)

func (g Gender) String() string {
	return string(g)
}

type User struct {
	ID           string    `json:"id"`
	Gender       Gender    `json:"gender"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"-"`
	Confirmed    bool      `json:"confirmed"`
	UpdateAt     time.Time `json:"-"`
	CreatedAt    time.Time `json:"-"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Make() *User {
	return &User{
		Gender: GenderOther,
	}
}

func (u *UserRepository) prepareUser(user *User) error {
	if len(user.Password) > 0 {
		hash, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.PasswordHash = hash
	}
	return nil
}

func (u *UserRepository) Find(key string, value any) (*User, error) {
	row := u.db.QueryRow(fmt.Sprintf(`
		SELECT
		    id, gender, firstname, lastname, email, password, confirmed, updated_at, created_at
		FROM users WHERE %s = $1
	`, key), value)

	err := row.Err()
	if err != nil {
		return nil, err
	}

	user := User{}
	err = row.Scan(
		&user.ID,
		&user.Gender,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.PasswordHash,
		&user.Confirmed,
		&user.UpdateAt,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) Get(uuid string) (*User, error) {
	return u.Find("id", uuid)
}

func (u *UserRepository) Persist(usr *User) (string, error) {
	if err := u.prepareUser(usr); err != nil {
		return "", err
	}
	if len(usr.ID) == 0 {
		res := u.db.QueryRow(`
			INSERT INTO users (gender, firstname, lastname, email, password, confirmed)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, usr.Gender, usr.Firstname, usr.Lastname, usr.Email, usr.PasswordHash, usr.Confirmed)
		if res.Err() != nil {
			return "", res.Err()
		}
		var id string
		err := res.Scan(&id)
		return id, err
	} else {
		_, err := u.db.Exec(`
			UPDATE users
			SET gender = $2, firstname = $3, lastname = $4, email = $5, password = $6, confirmed = $7, updated_at = $8
			WHERE id = $1
		`, usr.ID, usr.Gender, usr.Firstname, usr.Lastname, usr.Email, usr.PasswordHash, usr.Confirmed, time.Now())
		return usr.ID, err
	}
}

func (u *UserRepository) Delete(id string) error {
	_, err := u.db.Exec(`
		DELETE FROM refresh_chains WHERE id = $1
	`, id)
	return err
}
