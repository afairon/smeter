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
