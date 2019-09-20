-- +goose Up
CREATE TABLE maker.vat_flux
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    src       TEXT,
    dst       TEXT,
    wad       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_flux_header_index
    ON maker.vat_flux (header_id);

CREATE INDEX vat_flux_ilk_index
    ON maker.vat_flux (ilk_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_flux INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vat_flux_header_index;
DROP INDEX maker.vat_flux_ilk_index;

DROP TABLE maker.vat_flux;

ALTER TABLE public.checked_headers
    DROP COLUMN vat_flux;
