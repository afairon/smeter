package models

import (
	"database/sql"
	"time"

	"github.com/afairon/smeter/internal/message"
)

// AddPower saves power metric to database.
func AddPower(db *sql.DB, power *message.Power) error {
	const sql = `
		INSERT INTO
			public.power_metrics
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

	timestamp := time.Unix(power.GetTime(), 0)
	sensorID := power.GetSensorID()
	value := power.GetValue()

	// Check if sensor is active and is a sensor of type power.
	ok, err := SensorActive(db, message.SensorType_POWER, sensorID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}

	_, err = db.Exec(sql, timestamp, sensorID, value)

	return err
}
