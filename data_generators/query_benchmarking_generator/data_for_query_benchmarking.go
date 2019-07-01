//  VulcanizeDB
//  Copyright © 2019 Vulcanize
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package query_benchmarking_generator

import (
	"fmt"
	"github.com/vulcanize/mcd_transformers/data_generators/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

type BenchmarkingDataGeneratorState struct {
	shared.GeneratorState
}

func NewBenchmarkingDataGeneratorState(db *postgres.DB) BenchmarkingDataGeneratorState {
	generatorState := BenchmarkingDataGeneratorState{}
	generatorState.DB = db
	return generatorState
}

func (state *BenchmarkingDataGeneratorState) GenerateDataForQueryTesting(numberOfIlks, numberOfUrns, numberOfAdditionalBlocks int) error {
	pgTx, txErr := state.DB.Beginx()
	if txErr != nil {
		return txErr
	}

	state.PgTx = pgTx

	_, nodeErr := state.InsertEthNode()
	if nodeErr != nil {
		return state.handleErrorWithRollback(nodeErr)
	}

	ilkErr := state.generateIlks(numberOfIlks)
	if ilkErr != nil {
		return state.handleErrorWithRollback(ilkErr)
	}

	urnErr := state.generateUrns(numberOfUrns)
	if urnErr != nil {
		return state.handleErrorWithRollback(urnErr)
	}

	for i := 1; i <= numberOfAdditionalBlocks; i++ {
		state.CurrentHeader = fakes.GetFakeHeaderWithTimestamp(int64(i), int64(i))
		state.CurrentHeader.Hash = test_data.AlreadySeededRandomString(10)
		headerErr := state.InsertCurrentHeader()
		if headerErr != nil {
			return fmt.Errorf("error inserting current header: %v", headerErr)
		}

		storageErr := state.generateStorageRecordsForIlks()
		if storageErr != nil {
			return state.handleErrorWithRollback(storageErr)
		}

		urnStorageErr := state.generateStorageRecordsForUrns()
		if urnStorageErr != nil {
			return state.handleErrorWithRollback(urnStorageErr)
		}

	}

	return state.PgTx.Commit()
}

// creates n ilks
func (state *BenchmarkingDataGeneratorState) generateIlks(numberOfIlks int) error {
	for i := 1; i <= numberOfIlks; i++ {
		err := state.CreateIlk()
		if err != nil {
			return err
		}
	}
	return nil
}

// creates n urns per ilk
func (state *BenchmarkingDataGeneratorState) generateUrns(numberOfUrns int) error {
	for i := 1; i <= numberOfUrns; i++ {
		state.CurrentHeader = fakes.GetFakeHeaderWithTimestamp(int64(i), int64(i))
		state.CurrentHeader.Hash = test_data.AlreadySeededRandomString(10)
		headerErr := state.InsertCurrentHeader()
		if headerErr != nil {
			return fmt.Errorf("error inserting current header: %v", headerErr)
		}
		for _, ilkId := range state.Ilks {
			createUrnErr := state.CreateUrn(ilkId)
			if createUrnErr != nil {
				return createUrnErr
			}
		}
	}

	return nil
}

//generates storage records for each ilk
func (state *BenchmarkingDataGeneratorState) generateStorageRecordsForIlks() error {
	for _, ilkID := range state.Ilks {
		err := state.InsertInitialIlkData(ilkID)
		if err != nil {
			return err
		}
	}
	return nil
}

//generates storage records for each urn
func (state *BenchmarkingDataGeneratorState) generateStorageRecordsForUrns() error {
	for _, urnId := range state.Urns {
		guy := shared.GetRandomAddress()
		err := state.InsertInitialUrnData(urnId, guy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (state *BenchmarkingDataGeneratorState) handleErrorWithRollback(err error) error {
	rollbackErr := state.PgTx.Rollback()
	if rollbackErr != nil {
		return rollbackErr
	}
	return err
}
