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

// GetDevice retrieves device.
func GetDevice(db *sql.DB, device *message.Device) (*message.Device, error) {
	const sql = `
		SELECT
			id,
			name,
			active
		FROM
			public.devices
		WHERE
			devices.id = $1
	`

	deviceID := device.GetID()

	row := db.QueryRow(sql, deviceID)

	if err := row.Scan(&(*device).ID, &(*device).Name, &(*device).Active); err != nil {
		return nil, err
	}

	return device, nil
}
