# Example of Fuel smart contract implementation and RPC proxy adapter for Fuel SDK for the Fluent environment

## Semantic units

- Fuel smart contract implementation: crates/contracts
- Fuel adoption for Fluent-EE: crates/contracts
- E2E tests: e2e
- Transaction examples: examples
- GraphQL Proxy adapter implementation for Fuel SDK

## Requirements/initial setup (not limited to, just tested with/on)

- Ubuntu (version 23.10)
- Docker (version 27)
- GNU Make (version 4.3)
- Rust+rustup installation (version will be managed automatically by this setup)
- Golang installation (check 'fuel-proxy/go.mod' file for Golang version)
- NodeJS (version v20.12.1)
- Yarn (version 1.22.22)
- Running empty Redis instance or use docker compose in the project (check for dev creds in 'fuel-proxy/config.yaml')
- This project in some local directory
- Fluent node project (https://github.com/fluentlabs-xyz/fluent) in 'fluent' directory on the same directory level as
  the current project's directory

## How to run main components and send example Fuel transaction

1. Prepare project: `make prepare`
2. Run Redis inside docker: `./docker_start.sh` (will be running in background)
3. Run Fluent node: `make run_fluent_node`
4. Run Fuel proxy: `make run_fuel_proxy`
5. Send example fuel transaction: `make send_example_fuel_tx`

## How to reset states for running services (if you want to start from genesis state)
- To reset Fluent's node state just run it again with (it always starts from genesis state): `make run_fluent_node`
- To reset Redis state (running in docker) run (remember to run this command right after resetting Fluent's node state or it will be polluted with previously written data): `./docker_recreate_redis.sh`

## How to rebuild Fuel-EE smart contract and update Fluent node with it

1. This step can be skipped, but needed to make sure that everything goes right. Run `Fluent` node with
   `make run_fluent_node` and find in logs line looking like 'Initializing genesis chain=dev genesis=SOME_HASH_HERE'.
   Remember/save 'SOME_HASH_HERE' somewhere locally
2. `Fluentbase` project (https://github.com/fluentlabs-xyz/fluentbase) in 'fluentbase' directory on the same directory
   level as
   the current project's directory
3. Run: `make build_contracts_and_update_genesis_locally` - it will update contracts and genesis in fluentbase
   directory (there must be 3 files changed in total)
4. Switch `Fluentbase` onto a new branch (you choose the name not conflicting with existing branches) and push changes
   to github
5. For `Fluent` project in the root Cargo.toml change `branch` attr for all `Fluentbase` crates onto the name you picked
   above
6. Run `cargo update` for `Fluent`
7. Rerun `Fluent` node: `make run_fluent_node`. You must see updated hash for genesis in logs, e.g.: Initializing
   genesis chain=dev genesis=NEW_HASH_HERE. Compare 'SOME_HASH_HERE' with 'NEW_HASH_HERE' - they must be different, it
   means genesis was updated (it doesn't mean everything went fine, but at least we got updated contracts inside
   genesis)
