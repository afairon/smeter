package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/afairon/smeter/message"
)

// AddTemperature saves temperature metric to database.
func AddTemperature(db *sql.DB, req *message.Temperature) error {
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

	timestamp := time.Unix(req.GetTime(), 0)
	sensorID := req.GetSensorID()
	value := req.GetValue()

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

// GetAvgTemperature returns a bucket of average temperature.
func GetAvgTemperature(db *sql.DB, req *message.TemperatureRequest) (*sql.Rows, error) {
	const sql = `
		SELECT
			date_trunc($1, time) as bucket,
			sensor_id,
			avg(value)
		FROM
			public.temperature_metrics
		WHERE
			sensor_id = $2
		AND
			time >= $3
		AND
			time < $4
		GROUP BY
			bucket,
			sensor_id
		ORDER BY
			bucket desc
	`

	sensorID := req.GetSensorID()
	from := time.Unix(req.GetFrom(), 0)
	to := time.Unix(req.GetTo(), 0)
	bucket := req.GetBucket()

	row, err := db.Query(sql, strings.ToLower(bucket.String()), sensorID, from, to)
	if err != nil {
		return nil, err
	}

	return row, nil
}
