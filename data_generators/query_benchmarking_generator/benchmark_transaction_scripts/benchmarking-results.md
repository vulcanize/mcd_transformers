# Benchmarking Results


#### Data Generation 
Created 5 ilks and 500 additional blocks on on a fresh database:
- `go run data_generator.go -pg-connection-string postgres://<your local postgres user and password>@localhost:5432/vulcanize_private?sslmode=disable -generator-type benchmark -steps 500`
- seed: `1561152359659541000`


### 2019-06-24, Revision: `9044467ebbb91f717f90f68a973e135cb8fbdff3`; ilk and urn data
Notes:
- `spot_ilk_mat` and `spot_ilk_pip` were added to the inserted ilk data - this could potentially be why we're seeing an increase
in latency for `all_ilks`?
- urn data was added

#### `get_ilk`
```sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/get_ilk.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 15.599 ms
tps = 64.107708 (including connections establishing)
tps = 65.706052 (excluding connections establishing)
```
#### `all_ilks`
```sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/all_ilks.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 85.514 ms
tps = 11.693951 (including connections establishing)
tps = 11.750308 (excluding connections establishing)
```
#### `all_ilk_states`
```sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/all_ilk_states.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 5185.406 ms
tps = 0.192849 (including connections establishing)
tps = 0.192867 (excluding connections establishing)
```
#### `get_urn`
```sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/get_urn.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 12.437 ms
tps = 80.406880 (including connections establishing)
tps = 82.919411 (excluding connections establishing)
```

#### `all_urns`
```sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/all_urns.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 116.388 ms
tps = 8.591977 (including connections establishing)
tps = 8.620684 (excluding connections establishing)
```
#### `all_urn_states`
```sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/all_urn_states.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 3670.450 ms
tps = 0.272446 (including connections establishing)
tps = 0.272478 (excluding connections establishing)
```
-------------------
### 2019-06-24, with indexes [VDB-646](https://github.com/vulcanize/mcd_transformers/pull/135), just ilk data
#### `get_ilk`
_needed to update the ilk in the script_
```sql
pgbench vulcanize_private -f query_benchmarking_generator/benchmark_transaction_scripts/get_ilk.sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/get_ilk.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 15.110 ms
tps = 66.181169 (including connections establishing)
tps = 68.093301 (excluding connections establishing)
```
#### `all_ilks`
```sql
pgbench vulcanize_private -f query_benchmarking_generator/benchmark_transaction_scripts/all_ilks.sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/all_ilks.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 59.975 ms
tps = 16.673630 (including connections establishing)
tps = 16.766921 (excluding connections establishing)
```
#### `all_ilk_states`
_needed to update the ilk in the script_
```sql
pgbench vulcanize_private -f query_benchmarking_generator/benchmark_transaction_scripts/all_ilk_states.sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/all_ilk_states.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 5189.627 ms
tps = 0.192692 (including connections establishing)
tps = 0.192706 (excluding connections establishing)
```
-------------------
### 2019-06-21, Revision: `9044467ebbb91f717f90f68a973e135cb8fbdff3`; just Ilk data
#### `get_ilk`
_needed to update the ilk in the script_
```sql
pgbench vulcanize_private -f query_benchmarking_generator/benchmark_transaction_scripts/get_ilk.sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/get_ilk.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 15.590 ms
tps = 64.145548 (including connections establishing)
tps = 65.920751 (excluding connections establishing)
```
#### `all_ilks`
```sql
pgbench vulcanize_private -f query_benchmarking_generator/benchmark_transaction_scripts/all_ilks.sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/all_ilks.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 59.759 ms
tps = 16.733839 (including connections establishing)
tps = 16.832626 (excluding connections establishing)
```
#### `all_ilk_states`
_needed to update the ilk in the script_
```sql
pgbench vulcanize_private -f query_benchmarking_generator/benchmark_transaction_scripts/all_ilk_states.sql
transaction type: query_benchmarking_generator/benchmark_transaction_scripts/all_ilk_states.sql
scaling factor: 1
query mode: simple
number of clients: 1
number of threads: 1
number of transactions per client: 10
number of transactions actually processed: 10/10
latency average = 5311.753 ms
tps = 0.188262 (including connections establishing)
tps = 0.188280 (excluding connections establishing)
```
