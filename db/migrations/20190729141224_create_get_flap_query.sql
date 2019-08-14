-- +goose Up
CREATE TYPE api.relevant_flap_block AS (
    block_height BIGINT,
    block_hash TEXT,
    bid_id NUMERIC
    );

CREATE FUNCTION api.get_flap_blocks_before(bid_id NUMERIC, contract_address TEXT, block_height BIGINT)
    RETURNS SETOF api.relevant_flap_block AS
$$
SELECT block_number AS block_height, block_hash, kicks AS bid_id
FROM maker.flap_kicks
WHERE block_number <= get_flap_blocks_before.block_height
  AND kicks = get_flap_blocks_before.bid_id
  AND flap_kicks.contract_address = get_flap_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, bid_id
FROM maker.flap
WHERE block_number <= get_flap_blocks_before.block_height
  AND flap.bid_id = get_flap_blocks_before.bid_id
  AND flap.contract_address = get_flap_blocks_before.contract_address
ORDER BY block_height DESC
$$
    LANGUAGE sql
    STABLE;

CREATE TYPE api.flap AS (
    bid_id NUMERIC,
    guy TEXT,
    tic BIGINT,
    "end" BIGINT,
    lot NUMERIC,
    bid NUMERIC,
    gal TEXT,
    dealt BOOLEAN,
    created TIMESTAMP,
    updated TIMESTAMP
    );

CREATE FUNCTION api.get_flap(bid_id NUMERIC, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flap
AS
$$
WITH address AS (
    SELECT contract_address
    FROM maker.flap
    WHERE flap.bid_id = get_flap.bid_id
      AND block_number <= block_height
    LIMIT 1
),
     storage_values AS (
         SELECT bid_id, guy, tic, "end", lot, bid, gal
         FROM maker.flap
         WHERE bid_id = get_flap.bid_id
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deal AS (
         SELECT deal, bid_id
         FROM maker.deal
                  LEFT JOIN public.headers ON deal.header_id = headers.id
         WHERE deal.bid_id = get_flap.bid_id
           AND deal.contract_address IN (SELECT * FROM address)
           AND headers.block_number <= block_height
         ORDER BY bid_id, block_number DESC
         LIMIT 1
     ),
     relevant_blocks AS (
         SELECT *
         FROM api.get_flap_blocks_before(bid_id, (SELECT * FROM address), get_flap.block_height)
     ),
     created AS (
         SELECT DISTINCT ON (relevant_blocks.bid_id, relevant_blocks.block_height) relevant_blocks.block_height,
                                                                                   relevant_blocks.block_hash,
                                                                                   relevant_blocks.bid_id,
                                                                                   api.epoch_to_datetime(headers.block_timestamp) AS datetime
         FROM relevant_blocks
                  LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
         ORDER BY relevant_blocks.block_height ASC
         LIMIT 1
     ),
     updated AS (
         SELECT DISTINCT ON (relevant_blocks.bid_id, relevant_blocks.block_height) relevant_blocks.block_height,
                                                                                   relevant_blocks.block_hash,
                                                                                   relevant_blocks.bid_id,
                                                                                   api.epoch_to_datetime(headers.block_timestamp) AS datetime
         FROM relevant_blocks
                  LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
         ORDER BY relevant_blocks.block_height DESC
         LIMIT 1
     )

SELECT get_flap.bid_id,
       storage_values.guy,
       storage_values.tic,
       storage_values."end",
       storage_values.lot,
       storage_values.bid,
       storage_values.gal,
       CASE (SELECT COUNT(*) FROM deal)
           WHEN 0 THEN FALSE
           ELSE TRUE
           END AS dealt,
       created.datetime,
       updated.datetime
FROM maker.flap
    LEFT JOIN storage_values ON storage_values.bid_id = get_flap.bid_id
         JOIN created ON created.bid_id = flap.bid_id
         JOIN updated ON updated.bid_id = flap.bid_id
$$
    LANGUAGE sql
    STABLE;
-- +goose Down
DROP FUNCTION api.get_flap_blocks_before(NUMERIC, TEXT, BIGINT);
DROP TYPE api.relevant_flap_block CASCADE;
DROP FUNCTION api.get_flap(NUMERIC, BIGINT);
DROP TYPE api.flap CASCADE;