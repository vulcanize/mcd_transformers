package main

import (
	"flag"

	"fmt"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/history"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	startingBlockNumber := flag.Int("starting-number", -1, "First block to fill from")
	flag.Parse()
	config := cmd.LoadConfig(*environment)

	blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
	repository := repositories.NewPostgres(config.Database)
	numberOfBlocksCreated := history.PopulateBlocks(blockchain, repository, int64(*startingBlockNumber))
	fmt.Printf("Populated %d blocks", numberOfBlocksCreated)
}
