-- +goose Up
CREATE TYPE api.bid_event AS (
    bid_id NUMERIC,
    lot NUMERIC,
    bid_amount NUMERIC,
    act api.bid_act,
    block_height BIGINT,
    tx_idx INTEGER,
    contract_address TEXT
    );

COMMENT ON COLUMN api.bid_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.bid_event.tx_idx
    IS E'@omit';
COMMENT ON COLUMN api.bid_event.contract_address
    IS E'@omit';

CREATE FUNCTION api.all_bid_events()
    RETURNS SETOF api.bid_event AS
$$
SELECT *
FROM api.all_flip_bid_events()
UNION
SELECT *
FROM api.all_flap_bid_events()
UNION
SELECT *
FROM api.all_flop_bid_events()
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
DROP FUNCTION api.all_bid_events();
DROP TYPE api.bid_event CASCADE;
