package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/afairon/smeter/message"
)

// AddPower saves power metric to database.
func AddPower(db *sql.DB, req *message.Power) error {
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

	timestamp := time.Unix(req.GetTime(), 0)
	sensorID := req.GetSensorID()
	value := req.GetValue()

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

// GetPower returns power metrics.
func GetPower(db *sql.DB, req *message.PowerRequest) (*sql.Rows, error) {
	const sql = `
		SELECT
			time,
			sensor_id,
			value
		FROM
			public.power_metrics
		WHERE
			sensor_id = $1
		AND
			time >= $2
		AND
			time < $3
		ORDER BY
			time desc
	`

	sensorID := req.GetSensorID()
	from := time.Unix(req.GetFrom(), 0)
	to := time.Unix(req.GetTo(), 0)

	row, err := db.Query(sql, sensorID, from, to)
	if err != nil {
		return nil, err
	}

	return row, nil
}

// GetAvgConsumption returns a bucket of average energy consumption.
func GetAvgConsumption(db *sql.DB, req *message.ConsumptionRequest) (*sql.Rows, error) {
	const sql = `
		SELECT
			date_trunc($1, hourly.hour_bucket) as bucket,
			hourly.sensor_id,
			sum(hourly.avg)
		FROM
			(
				SELECT
					date_trunc('hour', time) as hour_bucket,
					sensor_id,
					avg(value)
				FROM
					public.power_metrics
				WHERE
					sensor_id = $2
				AND
					time >= $3
				AND
					time < $4
				GROUP BY
					hour_bucket, sensor_id
			)
		AS
			hourly
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
