BEGIN;

/* List of devices */
CREATE TABLE IF NOT EXISTS public.devices
(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name TEXT UNIQUE NOT NULL,
    active BOOL NOT NULL DEFAULT true
);

/* List of sensors attached to device
    Type:
        - 0: Unknown
        - 1: Power
        - 2: Temperature
        - 3: Humidity
*/
CREATE TABLE IF NOT EXISTS public.sensors
(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    device_id INT8 NOT NULL REFERENCES public.devices(id) ON DELETE CASCADE,
    type INT NOT NULL DEFAULT 1,
    name TEXT UNIQUE NOT NULL,
    active BOOL NOT NULL DEFAULT true
);

CREATE INDEX ON public.sensors(type, device_id, id);

/* Measurements of power */
CREATE TABLE IF NOT EXISTS public.power_metrics
(
    time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sensor_id INT8 NOT NULL REFERENCES public.sensors(id) ON DELETE CASCADE,
    value DOUBLE PRECISION NOT NULL
);

CREATE INDEX ON public.power_metrics(sensor_id, time DESC);

/* Measurements of temperature */
CREATE TABLE IF NOT EXISTS public.temperature_metrics
(
    time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sensor_id INT8 NOT NULL REFERENCES public.sensors(id) ON DELETE CASCADE,
    value DOUBLE PRECISION NOT NULL
);

CREATE INDEX ON public.temperature_metrics(sensor_id, time DESC);

/* Measurements of humidity */
CREATE TABLE IF NOT EXISTS public.humidity_metrics
(
    time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sensor_id INT8 NOT NULL REFERENCES public.sensors(id) ON DELETE CASCADE,
    value DOUBLE PRECISION NOT NULL
);

CREATE INDEX ON public.humidity_metrics(sensor_id, time DESC);

COMMIT;
