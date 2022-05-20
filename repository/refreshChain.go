package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jschaefer-io/IDaaS/utils"
)

type RefreshChain struct {
	Id        string
	Key       string
	AccessKey string
	ExpiresAt time.Time
}

type RefreshChainRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshChainRepository {
	return &RefreshChainRepository{
		db: db,
	}
}

func (rf *RefreshChainRepository) Make(accessKey string, expiration time.Time) *RefreshChain {
	key, _ := utils.GenerateRandomString(10)
	return &RefreshChain{
		Key:       key,
		AccessKey: accessKey,
		ExpiresAt: expiration,
	}
}

func (rf *RefreshChainRepository) Find(key string, value any) (*RefreshChain, error) {
	row := rf.db.QueryRow(fmt.Sprintf(`
		SELECT
		    id, key, access_key, expires_at
		FROM refresh_chains WHERE %s = $1
	`, key), value)

	err := row.Err()
	if err != nil {
		return nil, err
	}

	token := RefreshChain{}
	err = row.Scan(
		&token.Id,
		&token.Key,
		&token.AccessKey,
		&token.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (rf *RefreshChainRepository) Get(uuid string) (*RefreshChain, error) {
	return rf.Find("id", uuid)
}

func (rf *RefreshChainRepository) Persist(token *RefreshChain) (string, error) {
	if len(token.Id) == 0 {
		var uuid string
		res := rf.db.QueryRow(`
			INSERT INTO refresh_chains (key, access_key, expires_at)
			VALUES ($1, $2, $3)
			RETURNING id
		`, token.Key, token.AccessKey, token.ExpiresAt)
		err := res.Err()
		if err != nil {
			return uuid, err
		}
		if err = res.Scan(&uuid); err != nil {
			return uuid, err
		}
		return uuid, nil
	} else {
		_, err := rf.db.Exec(`
			UPDATE refresh_chains SET
			    key = $2,
			    access_key = $3,
			    expires_at = $4
			WHERE id = $1
		`, token.Id, token.Key, token.AccessKey, token.ExpiresAt)
		if err != nil {
			return token.Id, err
		}
		return token.Id, nil
	}
}

func (rf *RefreshChainRepository) Delete(id string) error {
	_, err := rf.db.Exec(`
		DELETE FROM refresh_chains WHERE id = $1
	`, id)
	return err
}
