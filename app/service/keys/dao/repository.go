package dao

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/DSuhinin/passbase-test-task/core/errors"

	"github.com/DSuhinin/passbase-test-task/app/service/keys/model"
)

// KeysRepositoryProvider provides an interface to work with Keys data.
type KeysRepositoryProvider interface {
	// CreateKey creates new Value entity.
	CreateKey() (*model.Key, error)
	// GetKey returns existing Value by it's ID.
	GetKey(keyID int) (*model.Key, error)
	// GetKeyByValue returns existing Value by it's value.
	GetKeyByValue(value string) (*model.Key, error)
	// GetKeys returns list of existing Keys.
	GetKeys() ([]model.Key, error)
	// RegenerateKey regenerates existing Value by it's ID.
	RegenerateKey(keyID int) (string, error)
	// DeleteKey deletes existing Value by it's ID.
	DeleteKey(keyID int) error
}

// KeysRepository represents object to work with Keys related data.
type KeysRepository struct {
	db *sqlx.DB
}

// NewKeysRepository creates new instance of KeysRepository to work with Keys data.
func NewKeysRepository(db *sqlx.DB) *KeysRepository {
	return &KeysRepository{db: db}
}

// CreateKey creates new Value entity.
func (r KeysRepository) CreateKey() (*model.Key, error) {
	m := model.Key{
		Value:     uuid.NewV4().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.QueryRow(`
		INSERT INTO public.keys(
			value,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3
		) RETURNING id`,
		m.Value,
		m.CreatedAt,
		m.UpdatedAt,
	).Scan(&m.ID); err != nil {
		return nil, errors.Wrap(err, "error creating new key")
	}

	return &m, nil
}

// GetKey returns existing Value by it's ID.
func (r KeysRepository) GetKey(keyID int) (*model.Key, error) {
	result := model.Key{}
	if err := r.db.Get(&result, `
		SELECT
			*
		FROM public.keys AS k
		WHERE k.id = $1`,
		keyID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "error getting key with id: %d", keyID)
	}
	return &result, nil
}

// GetKeyByValue returns existing Value by it's value.
func (r KeysRepository) GetKeyByValue(value string) (*model.Key, error) {
	result := model.Key{}
	if err := r.db.Get(&result, `
		SELECT
			*
		FROM public.keys AS k
		WHERE k.value = $1`,
		value,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "error getting key with value: %s", value)
	}
	return &result, nil
}

// GetKeys returns list of existing Keys.
func (r KeysRepository) GetKeys() ([]model.Key, error) {
	var result []model.Key
	if err := r.db.Select(&result, `
		SELECT
			*
		FROM public.keys AS k
		ORDER BY k.created_at ASC`,
	); err != nil {
		return nil, errors.Wrap(err, "error getting keys")
	}
	return result, nil
}

// RegenerateKey regenerates existing Value by it's ID.
func (r KeysRepository) RegenerateKey(keyID int) (string, error) {
	value := uuid.NewV4().String()
	if err := r.db.QueryRow(`
		UPDATE public.keys
		SET value = $1
		WHERE id = $2`,
		value,
		keyID,
	).Err(); err != nil {
		return "", errors.Wrap(err, "error regenerating key")
	}
	return value, nil
}

// DeleteKey deletes existing Value by it's ID.
func (r KeysRepository) DeleteKey(keyID int) error {
	result, err := r.db.Exec(`
		DELETE 
		FROM public.keys
		WHERE id = $1`,
		keyID,
	)

	if err != nil {
		return errors.Wrapf(err, " error deleting key by it's id: %d", keyID)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "error getting affected rows")
	}

	if affected == 0 {
		return errors.Wrapf(err, "key with id: %d not found", keyID)
	}

	return nil
}
