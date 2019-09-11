-- +goose Up
CREATE TABLE maker.tend
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC,
    bid        NUMERIC,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    log_idx    INTEGER NOT NULL,
    tx_idx     INTEGER NOT NULL,
    raw_log    JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX tend_header_index
    ON maker.tend (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN tend INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.tend_header_index;

DROP TABLE maker.tend;

ALTER TABLE public.checked_headers
    DROP COLUMN tend;
