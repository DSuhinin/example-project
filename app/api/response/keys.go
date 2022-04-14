package response

import "github.com/DSuhinin/passbase-test-task/app/service/keys/model"

// Key represents Key response object
type Key struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// NewKey creates new response object for `POST /keys` endpoint.
func NewKey(key *model.Key) Key {
	return Key{
		ID:    key.ID,
		Value: key.Value,
	}
}

// NewKeys creates new response object for `GET /keys` endpoint.
func NewKeys(keys []model.Key) []Key {
	result := make([]Key, 0, len(keys))
	for _, key := range keys {
		key := key
		result = append(result, NewKey(&key))
	}
	return result
}
