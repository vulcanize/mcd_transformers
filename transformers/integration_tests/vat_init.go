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

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_init"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("VatInit LogNoteTransformer", func() {
	vatInitConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.VatInitLabel,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       mcdConstants.VatABI(),
		Topic:             mcdConstants.VatInitSignature(),
	}

	It("transforms vat init log events", func() {
		blockNumber := int64(13171646)
		vatInitConfig.StartingBlockNumber = blockNumber
		vatInitConfig.EndingBlockNumber = blockNumber

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
			transformer.HexStringsToAddresses(vatInitConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatInitConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := shared.LogNoteTransformer{
			Config:     vatInitConfig,
			Converter:  &vat_init.VatInitConverter{},
			Repository: &vat_init.VatInitRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vatInitModel
		err = db.Select(&dbResults, `SELECT ilk_id from maker.vat_init`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
	})
})

type vatInitModel struct {
	Ilk              string `db:"ilk_id"`
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
