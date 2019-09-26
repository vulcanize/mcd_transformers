-- +goose Up
-- +goose StatementBegin

-- Function returning the history of a given ilk as of the given block height
CREATE FUNCTION api.all_ilk_states(ilk_identifier TEXT, block_height BIGINT DEFAULT api.max_block(),
                                   max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.ilk_state AS
$$
BEGIN
    RETURN QUERY (
        WITH relevant_blocks AS (
            SELECT get_ilk_blocks_before.block_height
            FROM api.get_ilk_blocks_before(ilk_identifier, all_ilk_states.block_height)
        )
        SELECT r.*
        FROM relevant_blocks,
             LATERAL api.get_ilk(ilk_identifier, relevant_blocks.block_height) r
        LIMIT all_ilk_states.max_results OFFSET all_ilk_states.result_offset
    );
END;
$$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_ilk_states(TEXT, BIGINT, INTEGER, INTEGER);
