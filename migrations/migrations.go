package migrations

import (
	"database/sql"

	"github.com/remind101/migrate"
)

var Migrations = []migrate.Migration{
	{
		ID: 1,
		Up: func(tx *sql.Tx) error {
			_, err := tx.Exec(migration1)
			return err
		},
		Down: func(tx *sql.Tx) error {
			return nil
		},
	},
	{
		ID: 2,
		Up: func(tx *sql.Tx) error {
			_, err := tx.Exec(migration2)
			return err
		},
		Down: func(tx *sql.Tx) error {
			return nil
		},
	},
	{
		ID: 3,
		Up: func(tx *sql.Tx) error {
			_, err := tx.Exec(migration3)
			return err
		},
		Down: func(tx *sql.Tx) error {
			return nil
		},
	},
	{
		ID: 4,
		Up: func(tx *sql.Tx) error {
			_, err := tx.Exec(migration4)
			return err
		},
		Down: func(tx *sql.Tx) error {
			return nil
		},
	},
	{
		ID: 5,
		Up: func(tx *sql.Tx) error {
			_, err := tx.Exec(migration5)
			return err
		},
		Down: func(tx *sql.Tx) error {
			return nil
		},
	},
	{
		ID: 6,
		Up: func(tx *sql.Tx) error {
			_, err := tx.Exec(migration6)
			return err
		},
		Down: func(tx *sql.Tx) error {
			return nil
		},
	},
	{
		ID: 7,
		Up: func(tx *sql.Tx) error {
			_, err := tx.Exec(migration7)
			return err
		},
		Down: func(tx *sql.Tx) error {
			return nil
		},
	},
}
