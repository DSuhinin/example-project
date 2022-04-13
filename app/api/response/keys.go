package response

import "github.com/DSuhinin/passbase-test-task/app/service/keys/model"

// Key represents Key response object
type Key struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

// NewKey creates new response object for `POST /keys` endpoint.
func NewKey(key *model.Key) Key {
	return Key{
		ID:  key.ID,
		Key: key.Value,
	}
}

// NewKeys creates new response object for `GET /keys` endpoint.
func NewKeys(keys []model.Key) []Key {
	result := make([]Key, 0, len(keys))
	for _, key := range keys {
		result = append(result, NewKey(&key))
	}
	return result
}
