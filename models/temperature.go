package models

import (
	"database/sql"
	"time"

	"github.com/afairon/smeter/internal/message"
)

// AddTemperature saves temperature metric to database.
func AddTemperature(db *sql.DB, temperature *message.Temperature) error {
	const sql = `
		INSERT INTO
			public.temperature_metrics
			(
				time,
				sensor_id,
				value
			)
		VALUES
			(
				$1,
				$2,
				$3
			)
	`

	timestamp := time.Unix(temperature.GetTime(), 0)
	sensorID := temperature.GetSensorID()
	value := temperature.GetValue()

	// Check if sensor is active and is a sensor of type temperature.
	ok, err := SensorActive(db, message.SensorType_TEMPERATURE, sensorID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}

	_, err = db.Exec(sql, timestamp, sensorID, value)

	return err
}
