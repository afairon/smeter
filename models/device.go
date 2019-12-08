package models

import (
	"database/sql"
	"strings"

	"github.com/afairon/smeter/message"
)

// AddDevice adds a new device to the database.
func AddDevice(db *sql.DB, req *message.Device) error {
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

	name := req.GetName()

	// Check if name is empty
	if strings.TrimSpace(name) == "" {
		return ErrInvalidName
	}

	_, err := db.Exec(sql, name)

	return err
}

// GetDevice retrieves device.
func GetDevice(db *sql.DB, req *message.Device) (*message.Device, error) {
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

	deviceID := req.GetID()

	row := db.QueryRow(sql, deviceID)

	if err := row.Scan(&(*req).ID, &(*req).Name, &(*req).Active); err != nil {
		return nil, err
	}

	return req, nil
}

// GetDevices returns list of devices.
func GetDevices(db *sql.DB, req *message.DevicesRequest) (*sql.Rows, error) {
	const sql = `
		SELECT
			id,
			name,
			active
		FROM
			public.devices
		ORDER BY
			id
		LIMIT
			$1
		OFFSET
			$2
	`

	limit := req.GetLimit()
	offset := req.GetOffset()

	if limit > 100 || limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	row, err := db.Query(sql, limit, offset)
	if err != nil {
		return nil, err
	}

	return row, nil
}

// DeviceActive checks if the device is active.
func DeviceActive(db *sql.DB, id int64) (bool, error) {
	const sql = `
			SELECT EXISTS
				(
					SELECT
						1
					FROM
						public.devices
					WHERE
						devices.id = $1
					AND
						devices.active = true
				)
	`

	var res bool
	row := db.QueryRow(sql, id)
	if err := row.Scan(&res); err != nil {
		return false, err
	}

	return res, nil
}
