package fixtures

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Fixtures represents object to load/unload dynamically fixtures.
type Fixtures struct {
	connection *sqlx.DB
}

// NewFixtures creates new Fixtures instance
func NewFixtures(connection *sqlx.DB) *Fixtures {
	return &Fixtures{
		connection: connection,
	}
}

// UnloadFixtures unloads postgres test fixtures.
func (f Fixtures) UnloadFixtures() error {
	tableList := []string{
		"public.keys",
	}

	for _, table := range tableList {
		// nolint
		_, err := f.connection.Exec(fmt.Sprintf(`
			DELETE
			FROM %s
		`, table))
		if err != nil {
			return err
		}
	}

	return nil
}
