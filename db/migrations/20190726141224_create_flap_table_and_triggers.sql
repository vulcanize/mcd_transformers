-- +goose Up
CREATE TABLE maker.flap
(
    id SERIAL PRIMARY KEY,
    block_number BIGINT DEFAULT 0,
    block_hash TEXT DEFAULT '',
    contract_address TEXT DEFAULT '',
    bid_id NUMERIC DEFAULT 0,
    guy TEXT DEFAULT '',
    tic BIGINT DEFAULT 0,
    "end" BIGINT DEFAULT 0,
    lot NUMERIC DEFAULT 0,
    bid NUMERIC DEFAULT 0,
    gal TEXT DEFAULT '',
--     dealt BOOLEAN, -- not sure how to populate this with triggers but get flap can figure it out
--     created TIMESTAMP, -- would be nice to include this here instead of figuring it out in get_flap
--     updated TIMESTAMP
    UNIQUE (block_number, bid_id)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_bid() RETURNS TRIGGER
AS $$
    BEGIN
        INSERT INTO maker.flap(bid_id, contract_address, block_number, block_hash, bid) VALUES(NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.bid)
            ON CONFLICT (bid_id, block_number) DO UPDATE SET bid = NEW.bid;
        return NEW;
    END
$$
LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_lot() RETURNS TRIGGER
AS $$
BEGIN
    INSERT INTO maker.flap(bid_id, contract_address, block_number, block_hash, lot) VALUES(NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.lot)
        ON CONFLICT (bid_id, block_number) DO UPDATE SET lot = NEW.lot;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_guy() RETURNS TRIGGER
AS $$
BEGIN
    INSERT INTO maker.flap(bid_id, contract_address, block_number, block_hash, guy) VALUES(NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.guy)
    ON CONFLICT (bid_id, block_number) DO UPDATE SET guy = NEW.guy;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_tic() RETURNS TRIGGER
AS $$
BEGIN
    INSERT INTO maker.flap(bid_id, contract_address, block_number, block_hash, tic) VALUES(NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.tic)
    ON CONFLICT (bid_id, block_number) DO UPDATE SET tic = NEW.tic;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_gal() RETURNS TRIGGER
AS $$
BEGIN
    INSERT INTO maker.flap(bid_id, contract_address, block_number, block_hash, gal) VALUES(NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW.gal)
    ON CONFLICT (bid_id, block_number) DO UPDATE SET gal = NEW.gal;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_end() RETURNS TRIGGER
AS $$
BEGIN
    INSERT INTO maker.flap(bid_id, contract_address, block_number, block_hash, "end") VALUES(NEW.bid_id, NEW.contract_address, NEW.block_number, NEW.block_hash, NEW."end")
    ON CONFLICT (bid_id, block_number) DO UPDATE SET "end" = NEW."end";
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_bid_bid AFTER INSERT OR UPDATE
    ON maker.flap_bid_bid
    FOR EACH ROW EXECUTE PROCEDURE maker.flap_bid();

CREATE TRIGGER flap_bid_lot AFTER INSERT OR UPDATE
    ON maker.flap_bid_lot
    FOR EACH ROW EXECUTE PROCEDURE maker.flap_lot();

CREATE TRIGGER flap_bid_guy AFTER INSERT OR UPDATE
    ON maker.flap_bid_guy
    FOR EACH ROW EXECUTE PROCEDURE maker.flap_guy();

CREATE TRIGGER flap_bid_tic AFTER INSERT OR UPDATE
    ON maker.flap_bid_tic
    FOR EACH ROW EXECUTE PROCEDURE maker.flap_tic();

CREATE TRIGGER flap_bid_gal AFTER INSERT OR UPDATE
    ON maker.flap_bid_gal
    FOR EACH ROW EXECUTE PROCEDURE maker.flap_gal();

CREATE TRIGGER flap_bid_end AFTER INSERT OR UPDATE
    ON maker.flap_bid_end
    FOR EACH ROW EXECUTE PROCEDURE maker.flap_end();

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