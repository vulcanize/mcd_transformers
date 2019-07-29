`Dockerfile` will build an alpine image containing:
- An $GOPATH with vulcanizedb, mcd_transformers, and goose
- An app directory with the vulcanizedb binary, startup_script.sh, and a (configurable) config.toml
Build with (e.g. from the project directory) `docker build ./ -t vulcanize_mcd_transformers:0.0.1 --build-arg USER`


## To use the container:
1. Setup a postgres database matching your config (e.g. `vulcanize_public`)
1. Determine values for the following _required_ environment variables:
    - `CLIENT_IPCPATH`
    - `DATABASE_NAME`
    - `DATABASE_HOSTNAME`
    - `DATABASE_PORT`
    - `DATABASE_USER`
    - `DATABASE_PASSWORD`
    - `FILESYSTEM_STORAGEDIFFSPATH`
1. Run with required environment variables: `docker run -e CLIENT_IPCPATH="" -e DATABASE_NAME="" -e DATABASE_HOSTNAME="" -e DATABASE_PORT="" -e DATABASE_USER="" -e DATABASE_PASSWORD="" -e FILESYSTEM_STORAGEDIFFSPATH="" vulcanize_mcd_transformers:0.0.1`.
    - This triggers `headerSync` + `composeAndExecute`.
    - NOTE: contract addresses are currently configured in `environments/example.toml` to point at the 0.2.9 deployment to Kovan.
       You can optionally replace any address with an environment variable, e.g. `-e CONTRACT_ADDRESS_MCD_FLIP_REP_A="0x1234"`.

NOTE: this file is written for execution on OS X, making use of `host.docker.internal` to access Postgres from the host.
For execution on linux, replace instances of `host.docker.internal` with `localhost` and run with `--network="host"`.

