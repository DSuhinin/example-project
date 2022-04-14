package keys

import (
	"github.com/DSuhinin/passbase-test-task/core/errors"

	"github.com/DSuhinin/passbase-test-task/app/service/keys/dao"
	"github.com/DSuhinin/passbase-test-task/app/service/keys/model"
)

// ServiceProvider provides Keys related operations.
type ServiceProvider interface {
	// CreateKey creates new Value entity.
	CreateKey() (*model.Key, error)
	// GetKey return existing Value entity by it's ID.
	GetKey(keyID int) (*model.Key, error)
	// GetKeys returns the list of existing Keys.
	GetKeys() ([]model.Key, error)
	// RegenerateKey regenerates  Value by it's ID.
	RegenerateKey(keyID int) (*model.Key, error)
	// DeleteKey deletes existing Value entity by it's ID.
	DeleteKey(keyID int) error
}

// Service implements ServiceProvider interface.
type Service struct {
	keysRepository dao.KeysRepositoryProvider
}

// NewService creates new Keys service instance.
func NewService(keysRepository dao.KeysRepositoryProvider) *Service {
	return &Service{
		keysRepository: keysRepository,
	}
}

// CreateKey creates new Value entity.
func (s Service) CreateKey() (*model.Key, error) {
	key, err := s.keysRepository.CreateKey()
	if err != nil {
		return nil, errors.InternalServerError.WithError(err)
	}

	return key, nil
}

// GetKey return existing Value entity by it's ID.
func (s Service) GetKey(keyID int) (*model.Key, error) {
	key, err := s.keysRepository.GetKey(keyID)
	if err != nil {
		return nil, errors.InternalServerError.WithError(err)
	}

	if key == nil {
		return nil, errors.EntityNotFoundError("keys")
	}

	return key, nil
}

// GetKeys returns the list of existing Keys.
func (s Service) GetKeys() ([]model.Key, error) {
	keyList, err := s.keysRepository.GetKeys()
	if err != nil {
		return nil, err
	}

	return keyList, nil
}

// RegenerateKey regenerates  Value by it's ID.
func (s Service) RegenerateKey(keyID int) (*model.Key, error) {
	key, err := s.GetKey(keyID)
	if err != nil {
		return nil, err
	}

	value, err := s.keysRepository.RegenerateKey(keyID)
	if err != nil {
		return nil, errors.InternalServerError.WithError(err)
	}
	key.Value = value
	return key, nil
}

// DeleteKey deletes existing Value entity by it's ID.
func (s Service) DeleteKey(keyID int) error {
	key, err := s.GetKey(keyID)
	if err != nil {
		return err
	}

	if err := s.keysRepository.DeleteKey(key.ID); err != nil {
		return errors.InternalServerError.WithError(err)
	}

	return nil
}
