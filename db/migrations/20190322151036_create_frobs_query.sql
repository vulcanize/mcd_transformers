-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.frob_event AS (
  ilk_name     TEXT,
  -- ilk object
  urn_id       TEXT,
  dink         NUMERIC,
  dart         NUMERIC,
  block_height BIGINT
  -- tx
);


CREATE FUNCTION api.frobs_for_urn(ilk_name TEXT, urn TEXT)
  RETURNS SETOF api.frob_event AS
$body$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1),
    urn AS (
      SELECT id FROM maker.urns
      WHERE ilk_id = (SELECT id FROM ilk)
        AND guy = $2
    )

  SELECT $1 AS ilk_name, $2 AS urn_id, dink, dart, block_number AS block_height
  FROM maker.vat_frob LEFT JOIN headers ON vat_frob.header_id = headers.id
  WHERE vat_frob.urn_id = (SELECT id FROM urn)
  ORDER BY block_number DESC
$body$
LANGUAGE sql STABLE;


CREATE FUNCTION api.all_frobs(ilk_name TEXT)
  RETURNS SETOF api.frob_event AS
$$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1)

  SELECT $1 AS ilk_name, guy AS urn_id, dink, dart, block_number AS block_height
  FROM maker.vat_frob
  LEFT JOIN maker.urns ON vat_frob.urn_id = urns.id
  LEFT JOIN headers    ON vat_frob.header_id = headers.id
  WHERE urns.ilk_id = (SELECT id FROM ilk)
  ORDER BY guy, block_number DESC
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.frobs_for_urn(TEXT, TEXT);
DROP FUNCTION api.all_frobs(TEXT);
DROP TYPE api.frob_event CASCADE;