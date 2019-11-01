-- +goose Up
CREATE TABLE api.ilk_state_history
(
    ilk_identifier TEXT,
    block_number   BIGINT,
    rate           NUMERIC   DEFAULT NULL,
    art            NUMERIC   DEFAULT NULL,
    spot           NUMERIC   DEFAULT NULL,
    line           NUMERIC   DEFAULT NULL,
    dust           NUMERIC   DEFAULT NULL,
    chop           NUMERIC   DEFAULT NULL,
    lump           NUMERIC   DEFAULT NULL,
    flip           TEXT      DEFAULT NULL,
    rho            NUMERIC   DEFAULT NULL,
    duty           NUMERIC   DEFAULT NULL,
    pip            TEXT      DEFAULT NULL,
    mat            NUMERIC   DEFAULT NULL,
    created        TIMESTAMP DEFAULT NULL,
    updated        TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (ilk_identifier, block_number)
);

COMMENT ON COLUMN api.ilk_state_history.ilk_identifier IS '@name id';


CREATE FUNCTION ilk_rate_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT rate
FROM maker.vat_ilk_rate
WHERE vat_ilk_rate.ilk_id = ilk_rate_at_block.ilk_id
  AND vat_ilk_rate.block_number < ilk_rate_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_art_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT art
FROM maker.vat_ilk_art
WHERE vat_ilk_art.ilk_id = ilk_art_at_block.ilk_id
  AND vat_ilk_art.block_number < ilk_art_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_spot_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT spot
FROM maker.vat_ilk_spot
WHERE vat_ilk_spot.ilk_id = ilk_spot_at_block.ilk_id
  AND vat_ilk_spot.block_number < ilk_spot_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_line_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT line
FROM maker.vat_ilk_line
WHERE vat_ilk_line.ilk_id = ilk_line_at_block.ilk_id
  AND vat_ilk_line.block_number < ilk_line_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_dust_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT dust
FROM maker.vat_ilk_dust
WHERE vat_ilk_dust.ilk_id = ilk_dust_at_block.ilk_id
  AND vat_ilk_dust.block_number < ilk_dust_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_chop_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT chop
FROM maker.cat_ilk_chop
WHERE cat_ilk_chop.ilk_id = ilk_chop_at_block.ilk_id
  AND cat_ilk_chop.block_number < ilk_chop_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_lump_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT lump
FROM maker.cat_ilk_lump
WHERE cat_ilk_lump.ilk_id = ilk_lump_at_block.ilk_id
  AND cat_ilk_lump.block_number < ilk_lump_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_flip_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS TEXT AS
$$
SELECT flip
FROM maker.cat_ilk_flip
WHERE cat_ilk_flip.ilk_id = ilk_flip_at_block.ilk_id
  AND cat_ilk_flip.block_number < ilk_flip_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_rho_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT rho
FROM maker.jug_ilk_rho
WHERE jug_ilk_rho.ilk_id = ilk_rho_at_block.ilk_id
  AND jug_ilk_rho.block_number < ilk_rho_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_duty_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT duty
FROM maker.jug_ilk_duty
WHERE jug_ilk_duty.ilk_id = ilk_duty_at_block.ilk_id
  AND jug_ilk_duty.block_number < ilk_duty_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_pip_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS TEXT AS
$$
SELECT pip
FROM maker.spot_ilk_pip
WHERE spot_ilk_pip.ilk_id = ilk_pip_at_block.ilk_id
  AND spot_ilk_pip.block_number < ilk_pip_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_mat_at_block(ilk_id INTEGER, block_number BIGINT) RETURNS NUMERIC AS
$$
SELECT mat
FROM maker.spot_ilk_mat
WHERE spot_ilk_mat.ilk_id = ilk_mat_at_block.ilk_id
  AND spot_ilk_mat.block_number < ilk_mat_at_block.block_number
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION ilk_time_created(ilk_id INTEGER) RETURNS TIMESTAMP AS
$$
SELECT api.epoch_to_datetime(block_timestamp)
FROM public.headers
         LEFT JOIN maker.vat_init ON vat_init.header_id = headers.id
WHERE vat_init.ilk_id = ilk_time_created.ilk_id
ORDER BY headers.block_number
LIMIT 1
$$
    LANGUAGE sql;


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_rate(new_diff maker.vat_ilk_rate) RETURNS maker.vat_ilk_rate
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            new_diff.rate,
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET rate = new_diff.rate;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_rates() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_rate(NEW);
    --PERFORM maker.update_later_rates();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_rate
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_rate
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_rates();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_art(new_diff maker.vat_ilk_art) RETURNS maker.vat_ilk_art
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.art,
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET art = new_diff.art;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_arts() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_art(NEW);
    --PERFORM maker.update_later_arts();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_art
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_art
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_arts();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_spot(new_diff maker.vat_ilk_spot) RETURNS maker.vat_ilk_spot
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.spot,
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET spot = new_diff.spot;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_spots() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_spot(NEW);
    --PERFORM maker.update_later_spots();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_spot
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_spot
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_spots();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_line(new_diff maker.vat_ilk_line) RETURNS maker.vat_ilk_line
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.line,
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET line = new_diff.line;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_lines() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_line(NEW);
    --PERFORM maker.update_later_lines();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_line
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_line
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_lines();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_dust(new_diff maker.vat_ilk_dust) RETURNS maker.vat_ilk_dust
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.dust,
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET dust = new_diff.dust;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_dusts() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_dust(NEW);
    --PERFORM maker.update_later_dusts();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_dust
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_dust
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_dusts();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_chop(new_diff maker.cat_ilk_chop) RETURNS maker.cat_ilk_chop
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.chop,
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET chop = new_diff.chop;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_chops() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_chop(NEW);
    --PERFORM maker.update_later_rates();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_chop
    AFTER INSERT OR UPDATE
    ON maker.cat_ilk_chop
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_chops();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_lump(new_diff maker.cat_ilk_lump) RETURNS maker.cat_ilk_lump
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.lump,
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET lump = new_diff.lump;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_lumps() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_lump(NEW);
    --PERFORM maker.update_later_rates();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_lump
    AFTER INSERT OR UPDATE
    ON maker.cat_ilk_lump
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_lumps();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flip(new_diff maker.cat_ilk_flip) RETURNS maker.cat_ilk_flip
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.flip,
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET flip = new_diff.flip;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_flips() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flip(NEW);
    --PERFORM maker.update_later_rates();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_flip
    AFTER INSERT OR UPDATE
    ON maker.cat_ilk_flip
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_flips();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_rho(new_diff maker.jug_ilk_rho) RETURNS maker.jug_ilk_rho
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.rho,
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET rho = new_diff.rho;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_rhos() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_rho(NEW);
    --PERFORM maker.update_later_rates();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_rho
    AFTER INSERT OR UPDATE
    ON maker.jug_ilk_rho
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_rhos();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_duty(new_diff maker.jug_ilk_duty) RETURNS maker.jug_ilk_duty
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.duty,
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET duty = new_diff.duty;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_duties() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_duty(NEW);
    --PERFORM maker.update_later_rates();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_duty
    AFTER INSERT OR UPDATE
    ON maker.jug_ilk_duty
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_duties();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_pip(new_diff maker.spot_ilk_pip) RETURNS maker.spot_ilk_pip
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.pip,
            ilk_mat_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET pip = new_diff.pip;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_pips() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_pip(NEW);
    --PERFORM maker.update_later_rates();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_pip
    AFTER INSERT OR UPDATE
    ON maker.spot_ilk_pip
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_pips();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_mat(new_diff maker.spot_ilk_mat) RETURNS maker.spot_ilk_mat
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE hash = new_diff.block_hash AND block_number = new_diff.block_number);
BEGIN
    INSERT
    INTO api.ilk_state_history (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                                pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            new_diff.block_number,
            ilk_rate_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_art_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_spot_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_line_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_dust_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_chop_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_lump_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_flip_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_rho_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_duty_at_block(new_diff.ilk_id, new_diff.block_number),
            ilk_pip_at_block(new_diff.ilk_id, new_diff.block_number),
            new_diff.mat,
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET mat = new_diff.mat;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_mats() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_mat(NEW);
    --PERFORM maker.update_later_rates();
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_mat
    AFTER INSERT OR UPDATE
    ON maker.spot_ilk_mat
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_mats();


-- +goose Down
DROP TRIGGER ilk_mat ON maker.spot_ilk_mat;
DROP TRIGGER ilk_pip ON maker.spot_ilk_pip;
DROP TRIGGER ilk_duty ON maker.jug_ilk_duty;
DROP TRIGGER ilk_rho ON maker.jug_ilk_rho;
DROP TRIGGER ilk_flip ON maker.cat_ilk_flip;
DROP TRIGGER ilk_lump ON maker.cat_ilk_lump;
DROP TRIGGER ilk_chop ON maker.cat_ilk_chop;
DROP TRIGGER ilk_dust ON maker.vat_ilk_dust;
DROP TRIGGER ilk_line ON maker.vat_ilk_line;
DROP TRIGGER ilk_spot ON maker.vat_ilk_spot;
DROP TRIGGER ilk_art ON maker.vat_ilk_art;
DROP TRIGGER ilk_rate ON maker.vat_ilk_rate;

DROP FUNCTION maker.insert_ilk_mats();
DROP FUNCTION maker.insert_ilk_pips();
DROP FUNCTION maker.insert_ilk_duties();
DROP FUNCTION maker.insert_ilk_rhos();
DROP FUNCTION maker.insert_ilk_flips();
DROP FUNCTION maker.insert_ilk_lumps();
DROP FUNCTION maker.insert_ilk_chops();
DROP FUNCTION maker.insert_ilk_dusts();
DROP FUNCTION maker.insert_ilk_lines();
DROP FUNCTION maker.insert_ilk_spots();
DROP FUNCTION maker.insert_ilk_arts();
DROP FUNCTION maker.insert_ilk_rates();

DROP FUNCTION maker.insert_new_mat(maker.spot_ilk_mat);
DROP FUNCTION maker.insert_new_pip(maker.spot_ilk_pip);
DROP FUNCTION maker.insert_new_duty(maker.jug_ilk_duty);
DROP FUNCTION maker.insert_new_rho(maker.jug_ilk_rho);
DROP FUNCTION maker.insert_new_flip(maker.cat_ilk_flip);
DROP FUNCTION maker.insert_new_lump(maker.cat_ilk_lump);
DROP FUNCTION maker.insert_new_chop(maker.cat_ilk_chop);
DROP FUNCTION maker.insert_new_dust(maker.vat_ilk_dust);
DROP FUNCTION maker.insert_new_line(maker.vat_ilk_line);
DROP FUNCTION maker.insert_new_spot(maker.vat_ilk_spot);
DROP FUNCTION maker.insert_new_art(maker.vat_ilk_art);
DROP FUNCTION maker.insert_new_rate(maker.vat_ilk_rate);

DROP FUNCTION ilk_time_created(INTEGER);
DROP FUNCTION ilk_mat_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_pip_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_duty_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_rho_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_flip_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_lump_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_chop_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_dust_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_line_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_spot_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_art_at_block(INTEGER, BIGINT);
DROP FUNCTION ilk_rate_at_block(INTEGER, BIGINT);

DROP TABLE api.ilk_state_history;