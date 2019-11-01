-- +goose Up
CREATE TABLE api.current_urn_state
(
    urn_identifier TEXT,
    ilk_identifier TEXT,
    ink            NUMERIC   DEFAULT NULL,
    art            NUMERIC   DEFAULT NULL,
    created        TIMESTAMP DEFAULT NULL,
    updated        TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (urn_identifier, ilk_identifier)
);

COMMENT ON TABLE api.current_urn_state IS '@omit create,update,delete';
COMMENT ON COLUMN api.current_urn_state.urn_identifier IS '@name id';
COMMENT ON COLUMN api.current_urn_state.ilk_identifier IS '@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_urn_ink() RETURNS TRIGGER
AS
$$
DECLARE
    new_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE block_number = NEW.block_number);
BEGIN
    WITH ilk AS (
        SELECT ilks.id, ilks.identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON ilks.id = urns.ilk_id
        WHERE urns.id = NEW.urn_id)
    INSERT
    INTO api.current_urn_state (urn_identifier, ilk_identifier, ink, created, updated)
    VALUES ((SELECT identifier FROM maker.urns WHERE id = NEW.urn_id),
            (SELECT identifier FROM ilk),
            NEW.ink,
            new_block_timestamp,
            new_block_timestamp)
    ON CONFLICT (urn_identifier, ilk_identifier)
        DO UPDATE
        SET ink     = (
            CASE
                WHEN current_urn_state.ink IS NULL OR current_urn_state.updated < new_block_timestamp
                    THEN NEW.ink
                ELSE current_urn_state.ink END),
            created = LEAST(new_block_timestamp, current_urn_state.created),
            updated = GREATEST(new_block_timestamp, current_urn_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_urn_art() RETURNS TRIGGER
AS
$$
DECLARE
    new_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE block_number = NEW.block_number);
BEGIN
    WITH ilk AS (
        SELECT ilks.id, ilks.identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON ilks.id = urns.ilk_id
        WHERE urns.id = NEW.urn_id)
    INSERT
    INTO api.current_urn_state (urn_identifier, ilk_identifier, art, created, updated)
    VALUES ((SELECT identifier FROM maker.urns WHERE id = NEW.urn_id),
            (SELECT identifier FROM ilk),
            NEW.art,
            new_block_timestamp,
            new_block_timestamp)
    ON CONFLICT (urn_identifier, ilk_identifier)
        DO UPDATE
        SET art     = (
            CASE
                WHEN current_urn_state.art IS NULL OR current_urn_state.updated < new_block_timestamp
                    THEN NEW.art
                ELSE current_urn_state.art END),
            updated = GREATEST(new_block_timestamp, current_urn_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


CREATE TRIGGER urn_ink
    AFTER INSERT OR UPDATE
    ON maker.vat_urn_ink
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_urn_ink();

CREATE TRIGGER urn_art
    AFTER INSERT OR UPDATE
    ON maker.vat_urn_art
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_urn_art();


-- +goose Down
DROP TRIGGER urn_art ON maker.vat_urn_art;
DROP TRIGGER urn_ink ON maker.vat_urn_ink;

DROP FUNCTION maker.insert_urn_art();
DROP FUNCTION maker.insert_urn_ink();

DROP TABLE api.current_urn_state;