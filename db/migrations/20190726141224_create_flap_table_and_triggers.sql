-- +goose Up
CREATE TABLE maker.flap
(
    id               SERIAL PRIMARY KEY,
    block_number     BIGINT  DEFAULT NULL,
    block_hash       TEXT    DEFAULT NULL,
    contract_address TEXT    DEFAULT NULL,
    bid_id           NUMERIC DEFAULT NULL,
    guy              TEXT    DEFAULT NULL,
    tic              BIGINT  DEFAULT NULL,
    "end"            BIGINT  DEFAULT NULL,
    lot              NUMERIC DEFAULT NULL,
    bid              NUMERIC DEFAULT NULL,
    gal              TEXT    DEFAULT NULL,
--     created TIMESTAMP, -- would be nice to include this here instead of figuring it out in get_flap
--     updated TIMESTAMP
    UNIQUE (block_number, bid_id)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_bid() RETURNS TRIGGER
AS
$$
BEGIN
    WITH lot AS (
        SELECT lot
        FROM maker.flap
        WHERE lot IS NOT NULL
        ORDER BY block_number
        LIMIT 1
    ),
         "end" AS (
             SELECT "end"
             FROM maker.flap
             WHERE "end" IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         tic AS (
             SELECT tic
             FROM maker.flap
             WHERE tic IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         guy AS (
             SELECT guy
             FROM maker.flap
             WHERE guy IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         gal AS (
             SELECT gal
             FROM maker.flap
             WHERE gal IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         )
    INSERT
    INTO maker.flap(bid_id, contract_address, block_number, block_hash, bid, lot, "end", tic, guy, gal)
    VALUES (NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.bid, (SELECT lot FROM lot),
            (SELECT "end" FROM "end"), (SELECT tic FROM tic), (SELECT guy FROM guy), (SELECT gal FROM gal))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET bid = NEW.bid;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_lot() RETURNS TRIGGER
AS
$$
BEGIN
    WITH bid AS (
        SELECT bid
        FROM maker.flap
        WHERE bid IS NOT NULL
        ORDER BY block_number
        LIMIT 1
    ),
         "end" AS (
             SELECT "end"
             FROM maker.flap
             WHERE "end" IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         tic AS (
             SELECT tic
             FROM maker.flap
             WHERE tic IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         guy AS (
             SELECT guy
             FROM maker.flap
             WHERE guy IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         gal AS (
             SELECT gal
             FROM maker.flap
             WHERE gal IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         )
    INSERT
    INTO maker.flap(bid_id, contract_address, block_number, block_hash, lot, bid, "end", tic, guy, gal)
    VALUES (NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.lot, (SELECT bid FROM bid),
            (SELECT "end" FROM "end"), (SELECT tic FROM tic), (SELECT guy FROM guy), (SELECT gal FROM gal))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET lot = NEW.lot;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_guy() RETURNS TRIGGER
AS
$$
BEGIN
    WITH bid AS (
        SELECT bid
        FROM maker.flap
        WHERE bid IS NOT NULL
        ORDER BY block_number
        LIMIT 1
    ),
         lot AS (
             SELECT lot
             FROM maker.flap
             WHERE lot IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         "end" AS (
             SELECT "end"
             FROM maker.flap
             WHERE "end" IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         tic AS (
             SELECT tic
             FROM maker.flap
             WHERE tic IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         gal AS (
             SELECT gal
             FROM maker.flap
             WHERE gal IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         )
    INSERT
    INTO maker.flap(bid_id, contract_address, block_number, block_hash, guy, bid, lot, "end", tic, gal)
    VALUES (NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.guy, (SELECT bid FROM bid),
            (SELECT lot FROM lot), (SELECT "end" FROM "end"), (SELECT tic FROM tic), (SELECT gal FROM gal))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET guy = NEW.guy;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_tic() RETURNS TRIGGER
AS
$$
BEGIN
    WITH bid AS (
        SELECT bid
        FROM maker.flap
        WHERE bid IS NOT NULL
        ORDER BY block_number
        LIMIT 1
    ),
         lot AS (
             SELECT lot
             FROM maker.flap
             WHERE lot IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         guy AS (
             SELECT guy
             FROM maker.flap
             WHERE guy IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         "end" AS (
             SELECT "end"
             FROM maker.flap
             WHERE "end" IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         gal AS (
             SELECT gal
             FROM maker.flap
             WHERE gal IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         )
    INSERT
    INTO maker.flap(bid_id, contract_address, block_number, block_hash, tic, bid, lot, guy, "end", gal)
    VALUES (NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.tic, (SELECT bid FROM bid),
            (SELECT lot FROM lot), (SELECT guy FROM guy), (SELECT "end" FROM "end"), (SELECT gal FROM gal))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET tic = NEW.tic;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_gal() RETURNS TRIGGER
AS
$$
BEGIN
    WITH bid AS (
        SELECT bid
        FROM maker.flap
        WHERE bid IS NOT NULL
        ORDER BY block_number
        LIMIT 1
    ),
         lot AS (
             SELECT lot
             FROM maker.flap
             WHERE lot IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         guy AS (
             SELECT guy
             FROM maker.flap
             WHERE guy IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         "end" AS (
             SELECT "end"
             FROM maker.flap
             WHERE "end" IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         tic AS (
             SELECT tic
             FROM maker.flap
             WHERE tic IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         )
    INSERT
    INTO maker.flap(bid_id, contract_address, block_number, block_hash, gal, bid, lot, guy, "end", tic)
    VALUES (NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.gal, (SELECT bid FROM bid),
            (SELECT lot FROM lot), (SELECT guy FROM guy), (SELECT "end" FROM "end"), (SELECT tic FROM tic))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET gal = NEW.gal;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_end() RETURNS TRIGGER
AS
$$
BEGIN
    WITH bid AS (
        SELECT bid
        FROM maker.flap
        WHERE bid IS NOT NULL
        ORDER BY block_number
        LIMIT 1
    ),
         lot AS (
             SELECT lot
             FROM maker.flap
             WHERE lot IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         guy AS (
             SELECT guy
             FROM maker.flap
             WHERE guy IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         gal AS (
             SELECT gal
             FROM maker.flap
             WHERE gal IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         ),
         tic AS (
             SELECT tic
             FROM maker.flap
             WHERE tic IS NOT NULL
             ORDER BY block_number
             LIMIT 1
         )
    INSERT
    INTO maker.flap(bid_id, contract_address, block_number, block_hash, "end", bid, lot, guy, gal, tic)
    VALUES (NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW."end", (SELECT bid FROM bid),
            (SELECT lot FROM lot), (SELECT guy FROM guy), (SELECT gal FROM gal), (SELECT tic FROM tic))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET "end" = NEW."end";
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_bid_bid
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_bid
    FOR EACH ROW
EXECUTE PROCEDURE maker.flap_bid();

CREATE TRIGGER flap_bid_lot
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_lot
    FOR EACH ROW
EXECUTE PROCEDURE maker.flap_lot();

CREATE TRIGGER flap_bid_guy
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_guy
    FOR EACH ROW
EXECUTE PROCEDURE maker.flap_guy();

CREATE TRIGGER flap_bid_tic
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_tic
    FOR EACH ROW
EXECUTE PROCEDURE maker.flap_tic();

CREATE TRIGGER flap_bid_gal
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_gal
    FOR EACH ROW
EXECUTE PROCEDURE maker.flap_gal();

CREATE TRIGGER flap_bid_end
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_end
    FOR EACH ROW
EXECUTE PROCEDURE maker.flap_end();

-- +goose Down
DROP TRIGGER flap_bid_bid ON maker.flap_bid_bid;
DROP TRIGGER flap_bid_lot ON maker.flap_bid_lot;
DROP TRIGGER flap_bid_guy ON maker.flap_bid_guy;
DROP TRIGGER flap_bid_tic ON maker.flap_bid_tic;
DROP TRIGGER flap_bid_gal ON maker.flap_bid_gal;
DROP TRIGGER flap_bid_end ON maker.flap_bid_end;

DROP FUNCTION maker.flap_end();
DROP FUNCTION maker.flap_gal();
DROP FUNCTION maker.flap_tic();
DROP FUNCTION maker.flap_guy();
DROP FUNCTION maker.flap_lot();
DROP FUNCTION maker.flap_bid();
DROP TABLE maker.flap;