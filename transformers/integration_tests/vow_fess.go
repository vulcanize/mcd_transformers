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

package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_fess"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = XDescribe("VowFess LogNoteTransformer", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)
	})

	vowFessConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VowFessLabel,
		ContractAddresses: []string{test_data.VowAddress()},
		ContractAbi:       constants.VowABI(),
		Topic:             constants.VowFessSignature(),
	}

	// TODO: replace block number when there is a fess event on the updated Vow
	It("transforms VowFess log events", func() {
		blockNumber := int64(9377319)
		vowFessConfig.StartingBlockNumber = blockNumber
		vowFessConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vowFessConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vowFessConfig.Topic)},
			header)
		Expect(len(logs)).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := shared.LogNoteTransformer{
			Config:     vowFessConfig,
			Converter:  &vow_fess.VowFessConverter{},
			Repository: &vow_fess.VowFessRepository{},
		}.NewLogNoteTransformer(db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vowFessModel
		err = db.Select(&dbResult, `SELECT tab, log_id from maker.vow_fess`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult[0].Tab).To(Equal("11000000000000000000000"))
	})
})

type vowFessModel struct {
	Tab   string
	LogID uint `db:"log_id"`
}
