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

package test_data

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthJugDripLog = types.Log{
	Address: common.HexToAddress(JugAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.JugDripSignature()),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000002444e2a5a866616b6520696c6b000000000000000000000000000000000000000000000000"),
	BlockNumber: 62,
	TxHash:      common.HexToHash("0xa34fd5cfcb125ebfc81d33633495701b531753669712092bdb8aa6159a240040"),
	TxIndex:     10,
	BlockHash:   fakes.FakeHash,
	Index:       11,
	Removed:     false,
}

var rawJugDripLog, _ = json.Marshal(EthJugDripLog)
var JugDripModel = shared.InsertionModel{
	SchemaName: "maker",
	TableName:  "jug_drip",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"log_idx": EthJugDripLog.Index,
		"tx_idx":  EthJugDripLog.TxIndex,
		"raw_log": rawJugDripLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x66616b6520696c6b000000000000000000000000000000000000000000000000",
	},
}
