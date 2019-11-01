// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package test_helpers

import . "github.com/onsi/gomega"

type BlockMetadata struct {
	BlockNumber int    `db:"block_number"`
	BlockHash   string `db:"block_hash"`
}

type VariableRes struct {
	BlockMetadata
	Value string
}

type AuctionVariableRes struct {
	VariableRes
	ContractAddress string `db:"contract_address"`
}

type MappingRes struct {
	BlockMetadata
	Key   string
	Value string
}

type DoubleMappingRes struct {
	BlockMetadata
	KeyOne string `db:"key_one"`
	KeyTwo string `db:"key_two"`
	Value  string
}

type FlapRes struct {
	BlockMetadata
	ContractAddress string `db:"contract_address"`
	ID              string
	BidID           string `db:"bid_id"`
	Guy             string
	Tic             string
	End             string
	Lot             string
	Bid             string
}

type FlopRes struct {
	BlockMetadata
	ContractAddress string `db:"contract_address"`
	ID              string
	BidID           string `db:"bid_id"`
	Guy             string
	Tic             string
	End             string
	Lot             string
	Bid             string
}

func AssertVariable(res VariableRes, blockNumber int, blockHash, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.Value).To(Equal(value))
}

func AssertMapping(res MappingRes, blockNumber int, blockHash, key, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.Key).To(Equal(key))
	Expect(res.Value).To(Equal(value))
}

func AssertDoubleMapping(res DoubleMappingRes, blockNumber int, blockHash, keyOne, keyTwo, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.KeyOne).To(Equal(keyOne))
	Expect(res.KeyTwo).To(Equal(keyTwo))
	Expect(res.Value).To(Equal(value))
}
