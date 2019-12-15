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

	sensorID := req.GetSensorID()
	from := req.GetFrom()
	to := req.GetTo()
	limit := req.GetLimit()
	offset := req.GetOffset()

	if limit > 60 || limit < 1 {
		limit = 60
	}
	if offset < 0 {
		offset = 0
	}

	if from == to {
		sql := `
		SELECT
			time,
			sensor_id,
			value
		FROM
			public.power_metrics
		WHERE
			sensor_id = $1
		ORDER BY
			time desc
		LIMIT
			$2
		OFFSET
			$3
	`
		row, err := db.Query(sql, sensorID, limit, offset)
		if err != nil {
			return nil, err
		}

		return row, nil
	}

	fromEpoch := time.Unix(from, 0)
	toEpoch := time.Unix(to, 0)

	if from > to {
		sql := `
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
		ORDER BY
			time desc
		LIMIT
			$3
		OFFSET
			$4
	`
		row, err := db.Query(sql, sensorID, fromEpoch, limit, offset)
		if err != nil {
			return nil, err
		}

		return row, nil
	}

	sql := `
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
		LIMIT
			$4
		OFFSET
			$5
	`

	row, err := db.Query(sql, sensorID, fromEpoch, toEpoch, limit, offset)
	if err != nil {
		return nil, err
	}

	return row, nil
}

// GetAvgConsumption returns a bucket of average energy consumption.
func GetAvgConsumption(db *sql.DB, req *message.ConsumptionRequest) (*sql.Rows, error) {

	sensorID := req.GetSensorID()
	from := req.GetFrom()
	to := req.GetTo()
	bucket := req.GetBucket()
	limit := req.GetLimit()
	offset := req.GetOffset()

	if limit > 100 || limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if from == to {
		sql := `
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
		LIMIT
			$3
		OFFSET
			$4
	`
		row, err := db.Query(sql, strings.ToLower(bucket.String()), sensorID, limit, offset)
		if err != nil {
			return nil, err
		}

		return row, nil
	}

	fromEpoch := time.Unix(from, 0)
	toEpoch := time.Unix(to, 0)

	if from > to {
		sql := `
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
		LIMIT
			$4
		OFFSET
			$5
	`
		row, err := db.Query(sql, strings.ToLower(bucket.String()), sensorID, fromEpoch, limit, offset)
		if err != nil {
			return nil, err
		}

		return row, nil
	}

	sql := `
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
		LIMIT
			$5
		OFFSET
			$6
	`

	row, err := db.Query(sql, strings.ToLower(bucket.String()), sensorID, fromEpoch, toEpoch, limit, offset)
	if err != nil {
		return nil, err
	}

	return row, nil
}
