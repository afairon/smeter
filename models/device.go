package models

import (
	"database/sql"
	"strings"

	"github.com/afairon/smeter/internal/message"
)

// AddDevice adds a new device to the database.
func AddDevice(db *sql.DB, device *message.Device) error {
	const sql = `
		INSERT INTO
			public.devices
			(
				name
			)
		VALUES
			(
				$1
			)
	`

	name := device.GetName()

	// Check if name is empty
	if strings.TrimSpace(name) == "" {
		return ErrInvalidName
	}

	_, err := db.Exec(sql, name)

	return err
}
