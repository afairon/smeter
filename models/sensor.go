package models

import (
	"database/sql"
	"strings"

	"github.com/afairon/smeter/message"
)

// AddSensor adds a new sensor to the database.
func AddSensor(db *sql.DB, req *message.Sensor) error {
	const sql = `
		INSERT INTO
			public.sensors
			(
				device_id,
				type,
				name
			)
		VALUES
			(
				$1,
				$2,
				$3
			)
	`

	deviceID := req.GetDeviceID()
	sensorType := req.GetType()
	name := req.GetName()

	// Check if name is empty
	if strings.TrimSpace(name) == "" {
		return ErrInvalidName
	}

	_, err := db.Exec(sql, deviceID, sensorType, name)

	return err
}

// GetSensor retrieves sensor.
func GetSensor(db *sql.DB, req *message.Sensor) (*message.Sensor, error) {
	const sql = `
		SELECT
			id,
			device_id,
			type,
			name,
			active
		FROM
			public.sensors
		WHERE
			sensors.id = $1
	`

	sensorID := req.GetID()

	row := db.QueryRow(sql, sensorID)

	if err := row.Scan(&(*req).ID, &(*req).DeviceID, &(*req).Type, &(*req).Name, &(*req).Active); err != nil {
		return nil, err
	}

	return req, nil
}

// GetSensors returns list of sensors.
func GetSensors(db *sql.DB, req *message.SensorRequest) (*sql.Rows, error) {
	const sql = `
		SELECT
			id,
			device_id,
			type,
			name,
			active
		FROM
			public.sensors
		WHERE
			sensors.device_id = $1
		ORDER BY
			id
		LIMIT
			$2
		OFFSET
			$3
	`

	deviceID := req.GetDeviceID()
	limit := req.GetLimit()
	offset := req.GetOffset()

	if limit > 100 || limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	row, err := db.Query(sql, deviceID, limit, offset)
	if err != nil {
		return nil, err
	}

	return row, nil
}

// SensorActive checks if the sensor is active.
func SensorActive(db *sql.DB, sensorType message.SensorType, id int64) (bool, error) {
	const sql = `
			SELECT EXISTS
				(
					SELECT
						1
					FROM
						public.sensors
					WHERE
						sensors.type = $1
					AND
						sensors.id = $2
					AND
						sensors.active = true
				)
	`

	var res bool
	row := db.QueryRow(sql, sensorType, id)
	if err := row.Scan(&res); err != nil {
		return false, err
	}

	return res, nil
}
