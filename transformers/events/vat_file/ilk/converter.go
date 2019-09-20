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

package ilk

import (
	"errors"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	constants2 "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type VatFileIlkConverter struct{}

func (VatFileIlkConverter) ToModels(logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	//NOTE: the vat contract defines its own custom Note event, rather than relying on DS-Note
	var models []shared.InsertionModel
	for _, log := range logs {
		err := verifyLog(log.Log)
		if err != nil {
			return nil, err
		}
		ilk := log.Log.Topics[1].Hex()
		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		data := log.Log.Topics[3].Big().String()

		model := shared.InsertionModel{
			TableName: "vat_file_ilk",
			OrderedColumns: []string{
				"header_id", string(constants2.IlkFK), "what", "data", "log_id",
			},
			ColumnValues: shared.ColumnValues{
				"what":      what,
				"data":      data,
				"header_id": log.HeaderID,
				"log_id":    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants2.IlkFK: ilk,
			},
		}
		models = append(models, model)
	}
	return models, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < constants.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
