-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.poke_event AS (
    ilk_id INTEGER,
    -- ilk object
    val NUMERIC,
    spot NUMERIC,
    block_height BIGINT,
    tx_idx INTEGER
    -- tx
    );

COMMENT ON COLUMN api.poke_event.ilk_id
    IS E'@omit';
COMMENT ON COLUMN api.poke_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.poke_event.tx_idx
    IS E'@omit';

CREATE FUNCTION api.max_timestamp()
    RETURNS NUMERIC AS
$$
SELECT max(block_timestamp)
FROM public.headers
$$
    LANGUAGE SQL
    STABLE;

CREATE FUNCTION api.all_poke_events(beginTime NUMERIC DEFAULT 0, endTime NUMERIC DEFAULT api.max_timestamp(), max_results INTEGER DEFAULT NULL)
    RETURNS SETOF api.poke_event AS
$body$
SELECT ilk_id, "value" AS val, spot, block_number AS block_height, tx_idx
FROM maker.spot_poke
         LEFT JOIN public.headers ON spot_poke.header_id = headers.id
WHERE block_timestamp BETWEEN beginTime AND endTime
ORDER BY block_height DESC
LIMIT all_poke_events.max_results
$body$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP FUNCTION api.all_poke_events(NUMERIC, NUMERIC, INTEGER);
DROP FUNCTION api.max_timestamp();
DROP TYPE api.poke_event CASCADE;