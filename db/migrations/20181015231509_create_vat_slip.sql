-- +goose Up
CREATE TABLE maker.vat_slip
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    usr       TEXT,
    wad       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_slip_header_index
    ON maker.vat_slip (header_id);

CREATE INDEX vat_slip_ilk_index
    ON maker.vat_slip (ilk_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_slip INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vat_slip_header_index;
DROP INDEX maker.vat_slip_ilk_index;

DROP TABLE maker.vat_slip;

ALTER TABLE public.checked_headers
    DROP COLUMN vat_slip;
