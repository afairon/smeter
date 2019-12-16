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
		RETURNING
			id
	`

	name := req.GetName()

	// Check if name is empty
	if strings.TrimSpace(name) == "" {
		return ErrInvalidName
	}

	err := db.QueryRow(sql, name).Scan(&req.ID)

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

	deviceStatus := req.GetStatus()
	limit := req.GetLimit()
	offset := req.GetOffset()

	if limit > 100 || limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	sql := `
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

	if deviceStatus == message.Status_ACTIVE {
		sql = `
		SELECT
			id,
			name,
			active
		FROM
			public.devices
		WHERE
			devices.active = true
		ORDER BY
			id
		LIMIT
			$1
		OFFSET
			$2
	`
		row, err := db.Query(sql, limit, offset)
		if err != nil {
			return nil, err
		}

		return row, nil
	}

	if deviceStatus == message.Status_INACTIVE {
		sql = `
		SELECT
			id,
			name,
			active
		FROM
			public.devices
		WHERE
			devices.active = false
		ORDER BY
			id
		LIMIT
			$1
		OFFSET
			$2
	`
		row, err := db.Query(sql, limit, offset)
		if err != nil {
			return nil, err
		}

		return row, nil
	}

	row, err := db.Query(sql, limit, offset)
	if err != nil {
		return nil, err
	}

	return row, nil
}

// UpdateDevice updates device name for a given ID.
func UpdateDevice(db *sql.DB, req *message.Device) error {

	deviceID := req.GetID()
	deviceName := req.GetName()
	deviceStatus := req.GetActive()

	if deviceID < 1 || strings.TrimSpace(deviceName) == "" {
		return ErrInvalidReq
	}

	sql := `
	UPDATE
		public.devices
	SET
		devices.name = $1,
		devices.active = $2
	WHERE
		devices.id = $3
	`

	_, err := db.Exec(sql, deviceName, deviceStatus, deviceID)

	return err
}

// DeleteDevice deletes device for a given ID.
func DeleteDevice(db *sql.DB, req *message.Device) error {

	deviceID := req.GetID()

	if deviceID < 1 {
		return ErrInvalidReq
	}

	sql := `
	DELETE
	FROM
		public.devices
	WHERE
		devices.id = $1
	`

	_, err := db.Exec(sql, deviceID)

	return err
}

// CountDevice counts number of devices.
func CountDevice(db *sql.DB, req *message.DeviceCountRequest) (*message.DeviceCount, error) {

	deviceStatus := req.GetStatus()

	sql := `
	SELECT
		COUNT (*)
	FROM
		public.devices
	`

	if deviceStatus == message.Status_ACTIVE {
		sql = `
		SELECT
			COUNT (*)
		FROM
			public.devices
		WHERE
			devices.active = true
		`
	}
	if deviceStatus == message.Status_INACTIVE {
		sql = `
		SELECT
			COUNT (*)
		FROM
			public.devices
		WHERE
			devices.active = false
		`
	}

	row := db.QueryRow(sql)

	res := &message.DeviceCount{}

	if err := row.Scan(&(*res).Count); err != nil {
		return nil, err
	}

	return res, nil
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
