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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_flux"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = Describe("VatFlux LogNoteTransformer", func() {
	vatFluxConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.VatFluxLabel,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       mcdConstants.VatABI(),
		Topic:             mcdConstants.VatFluxSignature(),
	}

	It("transforms VatFlux log events", func() {
		blockNumber := int64(13297132)
		vatFluxConfig.StartingBlockNumber = blockNumber
		vatFluxConfig.EndingBlockNumber = blockNumber

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
			transformer.HexStringsToAddresses(vatFluxConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatFluxConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		initializer := shared.LogNoteTransformer{
			Config:     vatFluxConfig,
			Converter:  &vat_flux.VatFluxConverter{},
			Repository: &vat_flux.VatFluxRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vatFluxModel
		err = db.Select(&dbResult, `SELECT ilk_id, src, dst, wad from maker.vat_flux`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		ilkID, err := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].Src).To(Equal("0x764B9b6326141C5912eBb6948b2b3d51B408d3E6"))
		Expect(dbResult[0].Dst).To(Equal("0x6bCc9f143D9C799E2C79DB9C921095130d371A16"))
		Expect(dbResult[0].Wad).To(Equal("1000000000000000000"))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(0)))
	})
})

type vatFluxModel struct {
	Ilk              string `db:"ilk_id"`
	Src              string
	Dst              string
	Wad              string
	TransactionIndex uint   `db:"tx_idx"`
	LogIndex         uint   `db:"log_idx"`
	Raw              []byte `db:"raw_log"`
}
