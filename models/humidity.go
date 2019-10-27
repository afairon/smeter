package models

import (
	"database/sql"
	"strings"
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

// GetAvgHumidity returns a bucket of average humidity.
func GetAvgHumidity(db *sql.DB, humidityReq *message.HumidityRequest) ([]*message.Humidity, error) {
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

	sensorID := humidityReq.GetSensorID()
	from := time.Unix(humidityReq.GetFrom(), 0)
	to := time.Unix(humidityReq.GetTo(), 0)
	bucket := humidityReq.GetBucket()

	row, err := db.Query(sql, strings.ToLower(bucket.String()), sensorID, from, to)
	if err != nil {
		return nil, err
	}

	var lstHumidities []*message.Humidity

	for row.Next() {
		humidity := message.Humidity{}
		timestamp := time.Time{}
		err = row.Scan(&timestamp, &humidity.SensorID, &humidity.Value)
		if err != nil {
			return nil, err
		}
		humidity.Time = timestamp.Unix()
		lstHumidities = append(lstHumidities, &humidity)
	}

	return lstHumidities, nil
}
