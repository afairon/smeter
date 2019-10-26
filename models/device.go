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

// GetDevices returns list of devices.
func GetDevices(db *sql.DB, devicesReq *message.DevicesRequest) ([]*message.Device, error) {
	const sql = `
		SELECT
			id,
			name,
			active
		FROM
			public.devices
		LIMIT
			$1
		OFFSET
			$2
	`

	limit := devicesReq.GetLimit()
	offset := devicesReq.GetOffset()

	if limit > 100 || limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var devices []*message.Device

	row, err := db.Query(sql, limit, offset)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		device := message.Device{}
		err = row.Scan(&device.ID, &device.Name, &device.Active)
		if err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}

	return devices, nil
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
