// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package vow

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Converter struct {
	db *postgres.DB
}

const (
	logDataRequired                    = true
	numTopicsRequired                  = 4
	What              event.ColumnName = "what"
	Data              event.ColumnName = "data"
)

func (Converter) ToModels(_ string, logs []core.HeaderSyncLog) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if err != nil {
			return nil, err
		}

		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		data := common.BytesToAddress(log.Log.Topics[3].Bytes()).String()

		result := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "cat_file_vow",
			OrderedColumns: []event.ColumnName{
				constants.HeaderFK, What, Data, event.LogFK,
			},
			ColumnValues: event.ColumnValues{
				constants.HeaderFK: log.HeaderID,
				What:               what,
				Data:               data,
				event.LogFK:        log.ID,
			},
		}
		results = append(results, result)
	}
	return results, nil
}

func (converter *Converter) SetDB(db *postgres.DB) {
	converter.db = db
}
