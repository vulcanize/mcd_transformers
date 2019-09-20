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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_suck"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

var _ = Describe("VatSuck Transformer", func() {
	vatSuckConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatSuckLabel,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatSuckSignature(),
	}

	It("transforms VatSuck log events", func() {
		blockNumber := int64(13424194)
		vatSuckConfig.StartingBlockNumber = blockNumber
		vatSuckConfig.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatSuckConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatSuckConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := shared.LogNoteTransformer{
			Config:     vatSuckConfig,
			Converter:  &vat_suck.VatSuckConverter{},
			Repository: &vat_suck.VatSuckRepository{},
		}.NewLogNoteTransformer(db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatSuckModel
		err = db.Select(&dbResults, `SELECT u, v, rad from maker.vat_suck`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.U).To(Equal("0x022688b43Bf76a9E6f4d3a96350ffDe90a752d25"))
		Expect(dbResult.V).To(Equal("0xE8fC4FC4D5Ab7Fa20BE296277EF157A8B0ec20Ce"))
		Expect(dbResult.Rad).To(Equal("117003578721766336311574734269352563186228"))
	})
})

type vatSuckModel struct {
	U   string
	V   string
	Rad string
}
