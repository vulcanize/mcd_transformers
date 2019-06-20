package probabilistic_data_generator

import (
	"fmt"
	"github.com/vulcanize/mcd_transformers/data_generators/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
)


type ProbabilisticDataGeneratorState struct{
	shared.GeneratorState
}

func NewProbabilisticDataGeneratorState(db *postgres.DB) ProbabilisticDataGeneratorState {
	generatorState := ProbabilisticDataGeneratorState{}
	generatorState.DB = db
	return generatorState
}

// Runs probabilistic generator for random ilk/urn interaction.
func (state *ProbabilisticDataGeneratorState) Run(steps int) error {
	pgTx, txErr := state.DB.Beginx()
	if txErr != nil {
		return txErr
	}

	state.PgTx = pgTx
	initErr := state.doInitialSetup()
	if initErr != nil {
		return initErr
	}

	var p float32
	var err error

	for i := 1; i <= steps; i++ {
		state.CurrentHeader = fakes.GetFakeHeaderWithTimestamp(int64(i), int64(i))
		state.CurrentHeader.Hash = test_data.AlreadySeededRandomString(10)
		headerErr := state.InsertCurrentHeader()
		if headerErr != nil {
			return fmt.Errorf("error inserting current header: %v", headerErr)
		}

		p = rand.Float32()
		if p < 0.2 { // Interact with Ilks
			err = state.TouchIlks()
			if err != nil {
				return fmt.Errorf("error touching ilks: %v", err)
			}
		} else { // Interact with Urns
			err = state.TouchUrns()
			if err != nil {
				return fmt.Errorf("error touching urns: %v", err)
			}
		}
	}
	return state.PgTx.Commit()
}

func (state *ProbabilisticDataGeneratorState) doInitialSetup() error {
	// This may or may not have been initialised, needed for a FK constraint
	_, nodeErr := state.InsertEthNode()
	if nodeErr != nil {
		return nodeErr
	}

	state.CurrentHeader = fakes.GetFakeHeaderWithTimestamp(0, 0)
	state.CurrentHeader.Hash = test_data.AlreadySeededRandomString(10)
	headerErr := state.InsertCurrentHeader()
	if headerErr != nil {
		return fmt.Errorf("could not insert initial header: %v", headerErr)
	}

	ilkErr := state.CreateIlk()
	if ilkErr != nil {
		return fmt.Errorf("could not create initial ilk: %v", ilkErr)
	}
	urnErr := state.CreateUrn()
	if urnErr != nil {
		return fmt.Errorf("could not create initial urn: %v", urnErr)
	}
	return nil
}
