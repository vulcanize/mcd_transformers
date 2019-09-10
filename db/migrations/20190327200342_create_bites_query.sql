-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.bite_event AS (
    ilk_identifier TEXT,
    -- ilk object
    urn_identifier TEXT,
    -- urn object
    bid_id NUMERIC,
    -- bid object
    ink NUMERIC,
    art NUMERIC,
    tab NUMERIC,
    block_height BIGINT,
    tx_idx INTEGER
    -- tx
    );

COMMENT ON COLUMN api.bite_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.bite_event.tx_idx
    IS E'@omit';

CREATE FUNCTION api.all_bites(ilk_identifier TEXT, max_results INTEGER DEFAULT NULL)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier, identifier AS urn_identifier, bite_identifier AS bid_id, ink, art, tab, block_number, tx_idx
FROM maker.bite
         LEFT JOIN maker.urns ON bite.urn_id = urns.id
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY urn_identifier, block_number DESC
LIMIT all_bites.max_results
$$
    LANGUAGE sql
    STABLE;


CREATE FUNCTION api.urn_bites(ilk_identifier TEXT, urn_identifier TEXT, max_results INTEGER DEFAULT NULL)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_bites.urn_identifier)

SELECT ilk_identifier, urn_bites.urn_identifier, bite_identifier AS bid_id, ink, art, tab, block_number, tx_idx
FROM maker.bite
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE bite.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT urn_bites.max_results
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_bites(TEXT, TEXT, INTEGER);
DROP FUNCTION api.all_bites(TEXT, INTEGER);
DROP TYPE api.bite_event CASCADE;
