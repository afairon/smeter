package models

import (
	"database/sql"
	"strings"

	"github.com/afairon/smeter/internal/message"
)

// AddSensor adds a new sensor to the database.
func AddSensor(db *sql.DB, sensor *message.Sensor) error {
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

	deviceID := sensor.GetDeviceID()
	sensorType := sensor.GetType()
	name := sensor.GetName()

	// Check if name is empty
	if strings.TrimSpace(name) == "" {
		return ErrInvalidName
	}

	_, err := db.Exec(sql, deviceID, sensorType, name)

	return err
}

// GetSensor retrieves sensor.
func GetSensor(db *sql.DB, sensor *message.Sensor) (*message.Sensor, error) {
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

	sensorID := sensor.GetID()

	row := db.QueryRow(sql, sensorID)

	if err := row.Scan(&(*sensor).ID, &(*sensor).DeviceID, &(*sensor).Type, &(*sensor).Name, &(*sensor).Active); err != nil {
		return nil, err
	}

	return sensor, nil
}
