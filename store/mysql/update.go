package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	NoRowsEffectedDelivery = fmt.Errorf("no rows affected when updating delivery")
)

// delivered will update the delivery's status to "DELIVERED" for the provided
// delivery id
func delivered(ctx context.Context, db *sql.DB, id string) error {
	const (
		updateDelivery = `UPDATE delivery SET status = 'DELIVERED', updatedAt = CURRENT_TIMESTAMP WHERE id = ?`
	)

	rslt, err := db.ExecContext(ctx, updateDelivery, id)
	if err != nil {
		return err
	}
	r, err := rslt.RowsAffected()
	if err != nil {
		return errors.New(fmt.Sprintf("updating delivery state, %s", err.Error()))
	}
	if r <= 0 {
		return NoRowsEffectedDelivery
	}

	return nil
}

// failed will update the delivery's status to "FAILED" for the provided
// delivery id
func failed(ctx context.Context, db *sql.DB, id string) error {
	const (
		updateDelivery = `UPDATE delivery SET status = 'FAILED', attempt = attempt + 1, updatedAt = CURRENT_TIMESTAMP WHERE id = ?`
	)

	rslt, err := db.ExecContext(ctx, updateDelivery, id)
	if err != nil {
		return err
	}
	r, err := rslt.RowsAffected()
	if err != nil {
		return errors.New(fmt.Sprintf("updating delivery state, %s", err.Error()))
	}
	if r <= 0 {
		return NoRowsEffectedDelivery
	}

	return nil
}
