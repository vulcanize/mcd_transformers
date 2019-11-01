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

package queries

import (
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("All flops query", func() {
	var (
		db              *postgres.DB
		flopKickRepo    flop_kick.FlopKickRepository
		flopRepo        flop.FlopStorageRepository
		dealRepo        deal.DealRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = "contract address"

		blockOne          = rand.Int()
		blockOneTimestamp = int64(111111111)

		blockTwo          = blockOne + 1
		blockTwoTimestamp = int64(222222222)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flopRepo = flop.FlopStorageRepository{}
		flopRepo.SetDB(db)
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the most recent flop for every bid id", func() {
		fakeBidIDOne := rand.Int()
		fakeBidIDTwo := fakeBidIDOne + 1

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(blockOneTimestamp, int64(blockOne))
		headerOneID, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(blockTwoTimestamp, int64(blockTwo))
		blockTwoHeader.Hash = "blockTwoHeader"
		headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidIDOne,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlopKickRepo:     flopKickRepo,
			FlopKickHeaderID: headerOneID,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		initialFlopOneStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidIDOne)
		test_helpers.CreateFlop(db, blockOneHeader, initialFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIDOne)), contractAddress)

		updatedFlopOneStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidIDOne)
		test_helpers.CreateFlop(db, blockTwoHeader, updatedFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIDOne)), contractAddress)

		flopStorageValuesTwo := test_helpers.GetFlopStorageValues(3, fakeBidIDTwo)
		test_helpers.CreateFlop(db, blockTwoHeader, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIDTwo)), contractAddress)

		contextErr = test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidIDTwo,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlopKickRepo:     flopKickRepo,
			FlopKickHeaderID: headerTwoID,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		var actualBids []test_helpers.FlopBid
		queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flops()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedBidOne := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIDOne), "false", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, updatedFlopOneStorageValues)
		expectedBidTwo := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIDTwo), "false", blockTwoHeader.Timestamp, blockTwoHeader.Timestamp, flopStorageValuesTwo)

		Expect(len(actualBids)).To(Equal(2))
		Expect(actualBids).To(ConsistOf([]test_helpers.FlopBid{
			expectedBidOne,
			expectedBidTwo,
		}))
	})

	Describe("result pagination", func() {
		var (
			headerID                                   int64
			header                                     core.Header
			fakeBidIDOne, fakeBidIDTwo                 int
			flopStorageValuesOne, flopStorageValuesTwo map[string]interface{}
		)

		BeforeEach(func() {
			fakeBidIDOne = rand.Int()
			fakeBidIDTwo = fakeBidIDOne + 1

			header = fakes.GetFakeHeaderWithTimestamp(blockOneTimestamp, int64(blockOne))
			var headerErr error
			headerID, headerErr = headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			flopStorageValuesOne = test_helpers.GetFlopStorageValues(1, fakeBidIDOne)
			test_helpers.CreateFlop(db, header, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIDOne)), contractAddress)

			flopStorageValuesTwo = test_helpers.GetFlopStorageValues(2, fakeBidIDTwo)
			test_helpers.CreateFlop(db, header, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIDTwo)), contractAddress)
		})

		It("limits results if max_results argument is provided", func() {
			contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidID:           fakeBidIDTwo,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickRepo:     flopKickRepo,
				FlopKickHeaderID: headerID,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			maxResults := 1
			var actualBids []test_helpers.FlopBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flops($1)`,
				maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIDTwo), "false", header.Timestamp,
				header.Timestamp, flopStorageValuesTwo)
			Expect(actualBids).To(Equal([]test_helpers.FlopBid{expectedBid}))
		})

		It("offsets results if offset is provided", func() {
			contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidID:           fakeBidIDOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickRepo:     flopKickRepo,
				FlopKickHeaderID: headerID,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			maxResults := 1
			resultOffset := 1
			var actualBids []test_helpers.FlopBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flops($1, $2)`,
				maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIDOne), "false", header.Timestamp,
				header.Timestamp, flopStorageValuesOne)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})
	})
})
