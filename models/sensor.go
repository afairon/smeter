package models

import (
	"database/sql"
	"strings"

	"github.com/afairon/smeter/message"
)

// AddSensor adds a new sensor to the database.
func AddSensor(db *sql.DB, req *message.Sensor) error {
	const sql = `
		INSERT INTO
			public.sensors
			(
				device_id,
				type,
				name
			)
		VALUES
			(
				$1,
				$2,
				$3
			)
		RETURNING
			id
	`

	deviceID := req.GetDeviceID()
	sensorType := req.GetType()
	name := req.GetName()

	// Check if name is empty
	if strings.TrimSpace(name) == "" {
		return ErrInvalidName
	}

	err := db.QueryRow(sql, deviceID, sensorType, name).Scan(&req.ID)

	return err
}

// GetSensor retrieves sensor.
func GetSensor(db *sql.DB, req *message.Sensor) (*message.Sensor, error) {
	const sql = `
		SELECT
			id,
			device_id,
			type,
			name,
			active
		FROM
			public.sensors
		WHERE
			sensors.id = $1
	`

	sensorID := req.GetID()

	row := db.QueryRow(sql, sensorID)

	if err := row.Scan(&(*req).ID, &(*req).DeviceID, &(*req).Type, &(*req).Name, &(*req).Active); err != nil {
		return nil, err
	}

	return req, nil
}

// GetSensors returns list of sensors.
func GetSensors(db *sql.DB, req *message.SensorsRequest) (*sql.Rows, error) {

	deviceID := req.GetDeviceID()
	sensorStatus := req.GetStatus()
	limit := req.GetLimit()
	offset := req.GetOffset()

	if limit > 100 || limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	sql := `
	SELECT
		id,
		device_id,
		type,
		name,
		active
	FROM
		public.sensors
	WHERE
		sensors.device_id = $1
	ORDER BY
		id
	LIMIT
		$2
	OFFSET
		$3
	`

	if deviceID < 1 {
		sql = `
			SELECT
				id,
				device_id,
				type,
				name,
				active
			FROM
				public.sensors
			ORDER BY
				id
			LIMIT
				$1
			OFFSET
				$2
			`
		if sensorStatus == message.Status_ACTIVE {
			sql = `
			SELECT
				id,
				device_id,
				type,
				name,
				active
			FROM
				public.sensors
			WHERE
				sensors.active = true
			ORDER BY
				id
			LIMIT
				$1
			OFFSET
				$2
			`
		}
		if sensorStatus == message.Status_INACTIVE {
			sql = `
			SELECT
				id,
				device_id,
				type,
				name,
				active
			FROM
				public.sensors
			WHERE
				sensors.active = false
			ORDER BY
				id
			LIMIT
				$1
			OFFSET
				$2
			`
		}
		row, err := db.Query(sql, limit, offset)
		if err != nil {
			return nil, err
		}

		return row, nil
	}

	if sensorStatus == message.Status_ACTIVE {
		sql = `
		SELECT
			id,
			device_id,
			type,
			name,
			active
		FROM
			public.sensors
		WHERE
			sensors.device_id = $1,
			sensors.active = true
		ORDER BY
			id
		LIMIT
			$2
		OFFSET
			$3
		`
	}

	if sensorStatus == message.Status_INACTIVE {
		sql = `
		SELECT
			id,
			device_id,
			type,
			name,
			active
		FROM
			public.sensors
		WHERE
			sensors.device_id = $1,
			sensors.active = false
		ORDER BY
			id
		LIMIT
			$2
		OFFSET
			$3
		`
	}

	row, err := db.Query(sql, deviceID, limit, offset)
	if err != nil {
		return nil, err
	}

	return row, nil
}

// UpdateSensor updates device ID and sensor name for a given ID.
func UpdateSensor(db *sql.DB, req *message.Sensor) error {

	sensorID := req.GetID()
	deviceID := req.GetDeviceID()
	sensorName := req.GetName()
	sensorStatus := req.GetActive()

	if sensorID < 1 || (deviceID < 1 && strings.TrimSpace(sensorName) == "") {
		return ErrInvalidReq
	}

	sql := `
	UPDATE
		public.sensors
	SET
		sensors.device_id = $1,
		sensors.name = $2,
		sensors.active = $3
	WHERE
		sensors.id = $4
	`

	if deviceID < 1 {
		sql = `
		UPDATE
			public.sensors
		SET
			sensors.name = $1,
			sensors.active = $2
		WHERE
			sensors.id = $3
		`
		_, err := db.Exec(sql, sensorName, sensorStatus, sensorID)

		return err
	}

	if strings.TrimSpace(sensorName) == "" {
		sql = `
		UPDATE
			public.sensors
		SET
			sensors.device_id = $1,
			sensors.active = $2
		WHERE
			sensors.id = $3
		`
		_, err := db.Exec(sql, deviceID, sensorStatus, sensorID)

		return err
	}

	_, err := db.Exec(sql, deviceID, sensorName, sensorStatus, sensorID)

	return err
}

// DeleteSensor deletes sensor for a given ID.
func DeleteSensor(db *sql.DB, req *message.Sensor) error {

	sensorID := req.GetID()

	if sensorID < 1 {
		return ErrInvalidReq
	}

	sql := `
	DELETE
		public.sensors
	WHERE
		sensors.id = $1
	`

	_, err := db.Exec(sql, sensorID)

	return err
}

// CountSensor counts number of sensors.
func CountSensor(db *sql.DB, req *message.SensorCountRequest) (*message.SensorCount, error) {

	sql := `
	SELECT
		COUNT (*)
	FROM
		public.sensors
	`

	row := db.QueryRow(sql)

	res := &message.SensorCount{}

	if err := row.Scan(&(*res).Count); err != nil {
		return nil, err
	}

	return res, nil
}

// SensorActive checks if the sensor is active.
func SensorActive(db *sql.DB, sensorType message.SensorType, id int64) (bool, error) {
	const sql = `
			SELECT EXISTS
				(
					SELECT
						1
					FROM
						public.sensors
					WHERE
						sensors.type = $1
					AND
						sensors.id = $2
					AND
						sensors.active = true
				)
	`

	var res bool
	row := db.QueryRow(sql, sensorType, id)
	if err := row.Scan(&res); err != nil {
		return false, err
	}

	return res, nil
}
