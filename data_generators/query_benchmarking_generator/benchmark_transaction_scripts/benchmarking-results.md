# Benchmarking Results


#### Data Generation 
Created 5 ilks and 500 additional blocks on on a fresh database:
- `go run data_generator.go -pg-connection-string postgres://elizabethengelman@localhost:5432/vulcanize_private?sslmode=disable -generator-type benchmark -steps 500`
- seed: `1561152359659541000`

### Revision: `9044467ebbb91f717f90f68a973e135cb8fbdff3`

```sql
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


```sql
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

```sql
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

### Revision: `9044467ebbb91f717f90f68a973e135cb8fbdff3`

