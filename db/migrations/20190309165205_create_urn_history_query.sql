-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION api.all_urn_states(ilk_identifier TEXT, urn_identifier TEXT,
                                   block_height BIGINT DEFAULT api.max_block(), max_results INTEGER DEFAULT NULL)
    RETURNS SETOF api.urn_state AS
$$
BEGIN
    RETURN QUERY (
        WITH urn_id AS (
            SELECT id
            FROM maker.urns
            WHERE urns.identifier = all_urn_states.urn_identifier
              AND urns.ilk_id = (SELECT id
                                 FROM maker.ilks
                                 WHERE ilks.identifier = all_urn_states.ilk_identifier)
        ),
             relevant_blocks AS (
                 SELECT block_number
                 FROM maker.vat_urn_ink
                 WHERE vat_urn_ink.urn_id = (SELECT * FROM urn_id)
                   AND block_number <= all_urn_states.block_height
                 UNION
                 SELECT block_number
                 FROM maker.vat_urn_art
                 WHERE vat_urn_art.urn_id = (SELECT * FROM urn_id)
                   AND block_number <= all_urn_states.block_height)
        SELECT r.*
        FROM relevant_blocks,
             LATERAL api.get_urn(ilk_identifier, urn_identifier, relevant_blocks.block_number) r
        ORDER BY relevant_blocks.block_number DESC
        LIMIT all_urn_states.max_results
    );
END;
$$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_urn_states(TEXT, TEXT, BIGINT, INTEGER);
