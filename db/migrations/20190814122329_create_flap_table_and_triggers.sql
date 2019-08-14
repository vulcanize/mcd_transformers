-- +goose Up
CREATE TABLE maker.flap
(
    id SERIAL PRIMARY KEY,
    block_number BIGINT,
    bid_id NUMERIC,
--     guy TEXT,
--     tic BIGINT,
--     "end" BIGINT,
    lot NUMERIC DEFAULT 0,
    bid NUMERIC DEFAULT 0,
--     gal TEXT,
--     dealt BOOLEAN,
--     created TIMESTAMP,
--     updated TIMESTAMP
    UNIQUE (block_number, bid_id)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_bid() RETURNS TRIGGER
AS $$
    BEGIN
        INSERT INTO maker.flap(bid_id, block_number, bid) VALUES(NEW.bid_id, NEW.block_number, NEW.bid)
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
    RAISE NOTICE 'TRIGGER called on %', TG_TABLE_NAME;
    RAISE NOTICE 'flap lot NEW %', NEW;
    INSERT INTO maker.flap(bid_id, block_number, lot) VALUES(NEW.bid_id, NEW.block_number, NEW.lot)
        ON CONFLICT (bid_id, block_number) DO UPDATE SET lot = NEW.lot;
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

-- +goose Down
DROP TRIGGER flap_bid_bid ON maker.flap_bid_bid;
DROP TRIGGER flap_bid_lot ON maker.flap_bid_lot;
DROP FUNCTION maker.flap_lot();
DROP FUNCTION maker.flap_bid();
DROP TABLE maker.flap;