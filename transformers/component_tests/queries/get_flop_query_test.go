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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("get flop query", func() {
	var (
		db              *postgres.DB
		flopKickRepo    flop_kick.FlopKickRepository
		dealRepo        deal.DealRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = "contract address"

		fakeBidID      = rand.Int()
		blockOne       = rand.Int()
		timestampOne   = int(rand.Int31())
		hashOne        = "hashOne"
		blockOneHeader = fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))

		blockTwo       = blockOne + 1
		timestampTwo   = timestampOne + 1000
		hashTwo        = "hashTwo"
		blockTwoHeader = fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))

		flopStorageValuesOne = test_helpers.GetFlopStorageValues(1, fakeBidID)
		flopStorageValuesTwo = test_helpers.GetFlopStorageValues(2, fakeBidID)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		blockOneHeader.Hash = hashOne
		blockTwoHeader.Hash = hashTwo
		rand.Seed(time.Now().UnixNano())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the specified flop", func() {
		headerID, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		err := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerID,
			},
			Dealt:            true,
			FlopKickRepo:     flopKickRepo,
			FlopKickHeaderID: headerID,
		})
		Expect(err).NotTo(HaveOccurred())

		test_helpers.CreateFlop(db, blockOneHeader, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		test_helpers.CreateFlop(db, blockTwoHeader, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

		expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidID), "true", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flopStorageValuesOne)

		var actualBid test_helpers.FlopBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidID, blockOne)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	It("gets created and updated blocks", func() {
		headerOneID, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		err := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerTwoID,
			},
			Dealt:            true,
			FlopKickRepo:     flopKickRepo,
			FlopKickHeaderID: headerOneID,
		})
		Expect(err).NotTo(HaveOccurred())

		test_helpers.CreateFlop(db, blockOneHeader, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)
		test_helpers.CreateFlop(db, blockTwoHeader, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

		blockThree := blockTwo + 1
		timestampThree := timestampTwo + 1000
		blockThreeHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampThree), int64(blockThree))
		flopStorageValuesThree := test_helpers.GetFlopStorageValues(3, fakeBidID)
		test_helpers.CreateFlop(db, blockThreeHeader, flopStorageValuesThree, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

		expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidID), "true", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flopStorageValuesTwo)

		var actualBid test_helpers.FlopBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidID, blockTwo)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	Describe("dealt", func() {
		It("is false if no deal events", func() {
			blockNumber := rand.Int()
			timestamp := int(rand.Int31())

			header := fakes.GetFakeHeaderWithTimestamp(int64(timestamp), int64(blockNumber))
			headerID, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			err := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidID:           fakeBidID,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickRepo:     flopKickRepo,
				FlopKickHeaderID: headerID,
			})
			Expect(err).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidID)
			test_helpers.CreateFlop(db, header, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidID), "false", header.Timestamp, header.Timestamp, flopStorageValues)

			var actualBid test_helpers.FlopBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidID, blockNumber)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})

		It("is false if deal event in later block", func() {
			headerID, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerOneErr).NotTo(HaveOccurred())

			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			err := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidID:           fakeBidID,
					ContractAddress: contractAddress,
					DealRepo:        dealRepo,
					DealHeaderID:    headerTwoID,
				},
				Dealt:            true,
				FlopKickRepo:     flopKickRepo,
				FlopKickHeaderID: headerID,
			})
			Expect(err).NotTo(HaveOccurred())

			test_helpers.CreateFlop(db, blockOneHeader, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)
			test_helpers.CreateFlop(db, blockTwoHeader, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			expectedBid := test_helpers.FlopBidFromValues(
				strconv.Itoa(fakeBidID), "false", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flopStorageValuesOne)

			var actualBid test_helpers.FlopBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidID, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})
	})
})
