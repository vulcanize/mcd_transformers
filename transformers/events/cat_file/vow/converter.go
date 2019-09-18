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
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type CatFileVowConverter struct{}

func (CatFileVowConverter) ToModels(ethLogs []types.Log) ([]shared.InsertionModel, error) {
	var results []shared.InsertionModel
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}

		what := shared.DecodeHexToText(ethLog.Topics[2].Hex())
		data := common.BytesToAddress(ethLog.Topics[3].Bytes()).String()

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}

		result := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "cat_file_vow",
			OrderedColumns: []string{
				"header_id", "what", "data", "tx_idx", "log_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"what":    what,
				"data":    data,
				"tx_idx":  ethLog.TxIndex,
				"log_idx": ethLog.Index,
				"raw_log": raw,
			},
			ForeignKeyValues: shared.ForeignKeyValues{},
		}
		results = append(results, result)
	}
	return results, nil
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
