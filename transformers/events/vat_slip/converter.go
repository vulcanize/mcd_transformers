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

package vat_slip

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type VatSlipConverter struct{}

func (VatSlipConverter) ToModels(_ string, ethLogs []types.Log) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := ethLog.Topics[1].Hex()
		usr := common.BytesToAddress(ethLog.Topics[2].Bytes()).String()
		wad := shared.ConvertInt256HexToBigInt(ethLog.Topics[3].Hex())

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "vat_slip",
			OrderedColumns: []string{
				"header_id", string(constants.IlkFK), "usr", "wad", "tx_idx", "log_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"usr":     usr,
				"wad":     wad.String(),
				"tx_idx":  ethLog.TxIndex,
				"log_idx": ethLog.Index,
				"raw_log": raw,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.IlkFK: ilk,
			},
		}
		models = append(models, model)
	}
	return models, nil
}

func verifyLog(log types.Log) error {
	numTopicInValidLog := 4
	if len(log.Topics) < numTopicInValidLog {
		return errors.New("log missing topics")
	}
	return nil
}
