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
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_fold"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("VatFold Transformer", func() {
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

	vatFoldConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.VatFoldLabel,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       mcdConstants.VatABI(),
		Topic:             mcdConstants.VatFoldSignature(),
	}

	It("transforms VatFold log events", func() {
		blockNumber := int64(13424126)
		vatFoldConfig.StartingBlockNumber = blockNumber
		vatFoldConfig.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatFoldConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFoldConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := shared.LogNoteTransformer{
			Config:     vatFoldConfig,
			Converter:  &vat_fold.VatFoldConverter{},
			Repository: &vat_fold.VatFoldRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatFoldModel
		err = db.Select(&dbResults, `SELECT urn_id, rate from maker.vat_fold`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		urnID, err := shared.GetOrCreateUrn("0x022688b43Bf76a9E6f4d3a96350ffDe90a752d25",
			"0x4447442d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult.Urn).To(Equal(strconv.FormatInt(urnID, 10)))
		Expect(dbResult.Rate).To(Equal("909758435446422415095"))
	})
})

type vatFoldModel struct {
	Ilk   string
	Urn   string `db:"urn_id"`
	Rate  string
	LogID uint `db:"log_id"`
}
