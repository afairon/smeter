package models

import (
	"database/sql"
	"strings"
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

// GetPower returns power metrics.
func GetPower(db *sql.DB, powerReq *message.PowerRequest) ([]*message.Power, error) {
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

	sensorID := powerReq.GetSensorID()
	from := time.Unix(powerReq.GetFrom(), 0)
	to := time.Unix(powerReq.GetTo(), 0)

	row, err := db.Query(sql, sensorID, from, to)
	if err != nil {
		return nil, err
	}

	var lstPowers []*message.Power

	for row.Next() {
		power := message.Power{}
		timestamp := time.Time{}
		err = row.Scan(&timestamp, &power.SensorID, &power.Value)
		if err != nil {
			return nil, err
		}
		power.Time = timestamp.Unix()
		lstPowers = append(lstPowers, &power)
	}

	return lstPowers, nil
}

// GetAvgConsumption returns a bucket of average energy consumption.
func GetAvgConsumption(db *sql.DB, consumptionReq *message.ConsumptionRequest) ([]*message.Energy, error) {
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

	sensorID := consumptionReq.GetSensorID()
	from := time.Unix(consumptionReq.GetFrom(), 0)
	to := time.Unix(consumptionReq.GetTo(), 0)
	bucket := consumptionReq.GetBucket()

	row, err := db.Query(sql, strings.ToLower(bucket.String()), sensorID, from, to)
	if err != nil {
		return nil, err
	}

	var lstEnergies []*message.Energy

	for row.Next() {
		energy := message.Energy{}
		timestamp := time.Time{}
		err = row.Scan(&timestamp, &energy.SensorID, &energy.Value)
		if err != nil {
			return nil, err
		}
		energy.Time = timestamp.Unix()
		lstEnergies = append(lstEnergies, &energy)
	}

	return lstEnergies, nil
}
