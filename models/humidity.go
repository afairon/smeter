package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/afairon/smeter/message"
)

// AddHumidity saves humidity metric to database.
func AddHumidity(db *sql.DB, req *message.Humidity) error {
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

	timestamp := time.Unix(req.GetTime(), 0)
	sensorID := req.GetSensorID()
	value := req.GetValue()

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

// GetHumidity returns data point of humidity metrics.
func GetHumidity(db *sql.DB, req *message.HumidityRequest) (*sql.Rows, error) {

	sensorID := req.GetSensorID()
	from := req.GetFrom()
	to := req.GetTo()
	limit := req.GetLimit()
	offset := req.GetOffset()

	if limit > 100 || limit < 1 {
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
			public.humidity_metrics
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
			public.humidity_metrics
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
			public.humidity_metrics
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

// GetAvgHumidity returns a bucket of average humidity.
func GetAvgHumidity(db *sql.DB, req *message.HumidityRequest) (*sql.Rows, error) {
	const sql = `
		SELECT
			date_trunc($1, time) as bucket,
			sensor_id,
			avg(value)
		FROM
			public.humidity_metrics
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
	to := time.Now()
	bucket := req.GetBucket()

	row, err := db.Query(sql, strings.ToLower(bucket.String()), sensorID, from, to)
	if err != nil {
		return nil, err
	}

	return row, nil
}
