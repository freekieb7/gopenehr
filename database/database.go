package database

import (
	"context"
)

type Database struct {
	*Postgres
}

func NewDatabase(postgres *Postgres) Database {
	return Database{
		postgres,
	}
}

func (db *Database) Migrate(ctx context.Context) error {
	// todo
	return nil
}
