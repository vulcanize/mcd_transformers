package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/mcd_transformers/data_generators/probabilistic_data_generator"
	"github.com/vulcanize/mcd_transformers/data_generators/query_benchmarking_generator"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"math/rand"
	"os"
	"time"
)

var (
	node = core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}
)

func main() {
	stepsPtr := flag.Int("steps", 100, "number of interactions to generate")
	seedPtr := flag.Int64("seed", -1,
		"optional seed for repeatability. Running same seed several times will lead to database constraint violations.")
	const defaultGeneratorType = "probabilistic"
	generatorType := flag.String("generator-type", defaultGeneratorType, "type of data generator to run")
	const defaultConnectionString = "postgres://vulcanize:vulcanize@localhost:5432/vulcanize_private?sslmode=disable"
	connectionStringPtr := flag.String("pg-connection-string", defaultConnectionString,
		"postgres connection string")

	flag.Parse()

	db := dbSetup(*connectionStringPtr)
	pg := postgres.DB{
		DB:     db,
		Node:   node,
		NodeID: 0,
	}

	printSeedInfo(*seedPtr)
	printDBWarning(*connectionStringPtr)

	startTime := time.Now()

	var runErr error
	switch *generatorType {
	case "probabilistic":
		fmt.Println("probabilistic")
		generatorState := probabilistic_data_generator.NewProbabilisticDataGeneratorState(&pg)
		runErr = generatorState.Run(*stepsPtr)

	case "benchmark":
		fmt.Println("benchmark")
		generatorState := query_benchmarking_generator.NewBenchmarkingDataGeneratorState(&pg)
		runErr = generatorState.GenerateDataForIlkQueryTesting(1, *stepsPtr)
	}

	if runErr != nil {
		fmt.Println("Error occurred while running generator: ", runErr.Error())
		fmt.Println("Exiting without writing any data to DB.")
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	speed := float64(*stepsPtr) / duration.Seconds()
	fmt.Printf("Simulated %v interactions in %v. (%.f/s)\n",
		*stepsPtr, duration.Round(time.Duration(time.Second)).String(), speed)
}

func dbSetup(connectionString string) *sqlx.DB {
	db, connectErr := sqlx.Connect("postgres", connectionString)
	if connectErr != nil {
		fmt.Println("Could not connect to DB: ", connectErr)
		os.Exit(1)
	}
	return db
}

func printSeedInfo(seed int64) {
	if seed != -1 {
		rand.Seed(seed)
		fmt.Println("\nUsing passed seed. If data from this seed is already in the DB, there will be database constraint errors.")
	} else {
		seed := time.Now().UnixNano()
		rand.Seed(seed)
		fmt.Printf("\nUsing current time as seed: %v. Pass this with '-seed' to reproduce results on a fresh DB.\n", seed)
	}
}

func printDBWarning(connectionString string) {
	fmt.Println("\nRunning this will write mock data to the DB you specified, possibly contaminating real data:")
	fmt.Println(connectionString)
	fmt.Println("------------------------------")
	fmt.Print("Do you want to continue? (y/n)")

	var input string
	_, err := fmt.Scanln(&input)
	if input != "y" || err != nil {
		os.Exit(0)
	}
}