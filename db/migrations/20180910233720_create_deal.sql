-- +goose Up
CREATE TABLE maker.deal (
  id               SERIAL PRIMARY KEY,
  header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  bid_id           NUMERIC NOT NULL,
  contract_address VARCHAR,
  log_idx          INTEGER NOT NULL,
  tx_idx           INTEGER NOT NULL,
  raw_log          JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN deal_checked BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
DROP TABLE maker.deal;

ALTER TABLE public.checked_headers
  DROP COLUMN deal_checked;
