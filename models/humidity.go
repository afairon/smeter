package models

import (
	"database/sql"
	"time"

	"github.com/afairon/smeter/internal/message"
)

// AddHumidity saves humidity metric to database.
func AddHumidity(db *sql.DB, humidity *message.Humidity) error {
	const sql = `
		INSERT INTO
			public.humidity_metrics
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

	timestamp := time.Unix(humidity.GetTime(), 0)
	sensorID := humidity.GetSensorID()
	value := humidity.GetValue()

	// Check if sensor is active and is a sensor of type humidity.
	ok, err := SensorActive(db, message.SensorType_HUMIDITY, sensorID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}

	_, err = db.Exec(sql, timestamp, sensorID, value)

	return err
}
