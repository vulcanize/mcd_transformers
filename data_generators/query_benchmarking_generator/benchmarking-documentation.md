# Benchmarking Queries

_These instructions assume that you are inserting data and benchmarking a local database called `vulcanize_private`._

## Generating the data
1. `cd` into the `data_generators/query_benchmarking_genrator` directory
1. run `go run data_generator.go` with the following flags:
    `-pg-connection-string postgres://<your database role and password>@localhost:5432/vulcanize_private?sslmode=disable`
    `-generator-type benchmark`

## Running `pgbench`
1. `pgbench` allows for custom SQL transaction scripts to be passed in. Scripts for the custom MCD query functions are located in
`data_generators/query_benchmarking_generator/benchmark_transaction_scripts`.
1. I have not found a good way to dynamically pass ilk and urn identifiers to these scripts, so for now a manual step is
required:
    1. Ilk and urn identifiers are randomly generated strings created in the data generator. After you've inserted the
    data into the db via the generator, you'll need to find an ilk identifier and an urn identifier from the database.
    `SELECT identifier FROM maker.ilks ORDER BY random() LIMIT 1;`
    `SELECT identifier FROM maker.urns ORDER BY random() LIMIT 1;`
    1. Then, update any of the benchmark scripts you're interested in with those identifiers.
1. Once the scripts you want to benchmark are updated, you can run PostgreSQL's `pgbench` tool with a script passed in:
`pgbench vulcanize_private -f query_benchmarking_generator/benchmark_transaction_scripts/all_ilk_states.sql`
    - You can benchmark multiple scripts at once by passing in multiple files:
    ```
    pgbench vulcanize_private -f data_generators/query_benchmarking_generator/benchmark_transaction_scripts/all_ilk_states.sql \
                              -f data_generators/query_benchmarking_generator/benchmark_transaction_scripts/all_urn_states.sql
    ```
    - You can also add weights to the scripts to determine how frequently they will be run:
    ```
    pgbench vulcanize_private -f data_generators/query_benchmarking_generator/benchmark_transaction_scripts/all_ilk_states.sql@1 \
                              -f data_generators/query_benchmarking_generator/benchmark_transaction_scripts/all_urn_states.sql@2
    ```
