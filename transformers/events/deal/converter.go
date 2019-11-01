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

package deal

import (
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
	logDataRequired   = true
	numTopicsRequired = 3
)

func (converter Converter) ToModels(_ string, logs []core.HeaderSyncLog) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		validationErr := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if validationErr != nil {
			return nil, validationErr
		}

		bidId := log.Log.Topics[2].Big()

		result := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "deal",
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				constants.BidColumn,
				constants.AddressColumn,
				event.LogFK,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:          log.HeaderID,
				constants.BidColumn:     bidId.String(),
				constants.AddressColumn: log.Log.Address.String(),
				event.LogFK:             log.ID,
			},
		}
		results = append(results, result)
	}

	return results, nil
}

func (converter *Converter) SetDB(db *postgres.DB) {
	converter.db = db
}
