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

	_, err := db.Exec(sql, timestamp, sensorID, value)

	return err
}
