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

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Single flip view", func() {
	var (
		db              *postgres.DB
		flipKickRepo    flip_kick.FlipKickRepository
		dealRepo        deal.DealRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = "flip"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flipKickRepo = flip_kick.FlipKickRepository{}
		flipKickRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		rand.Seed(time.Now().UnixNano())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets only the specified flip", func() {
		fakeBidID := rand.Int()
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1000

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
		headerID, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidID)
		test_helpers.CreateFlip(db, blockOneHeader, flipStorageValuesOne,
			test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

		ilkID, urnID, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerID,
			},
			Dealt:            true,
			IlkHex:           test_helpers.FakeIlk.Hex,
			UrnGuy:           test_data.FlipKickModel().ColumnValues["usr"].(string),
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderID: headerID,
		})
		Expect(err).NotTo(HaveOccurred())

		expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidID), strconv.FormatInt(ilkID, 10),
			strconv.FormatInt(urnID, 10), "true", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flipStorageValuesOne)

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
		blockTwoHeader.Hash = common.BytesToHash([]byte{5, 4, 3, 2, 1}).String()
		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidID)
		test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValuesTwo,
			test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

		var actualBid test_helpers.FlipBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip($1, $2, $3)`,
			fakeBidID, test_helpers.FakeIlk.Identifier, blockOne)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	Describe("dealt", func() {
		It("is false if no deal events", func() {
			fakeBidID := rand.Int()
			blockNumber := rand.Int()
			timestamp := int(rand.Int31())

			header := fakes.GetFakeHeaderWithTimestamp(int64(timestamp), int64(blockNumber))
			headerID, headerOneErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerOneErr).NotTo(HaveOccurred())

			flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidID)
			test_helpers.CreateFlip(db, header, flipStorageValues,
				test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			ilkID, urnID, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidID:           fakeBidID,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel().ColumnValues["usr"].(string),
				FlipKickRepo:     flipKickRepo,
				FlipKickHeaderID: headerID,
			})
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidID), strconv.FormatInt(ilkID, 10),
				strconv.FormatInt(urnID, 10), "false", header.Timestamp, header.Timestamp, flipStorageValues)

			var actualBid test_helpers.FlipBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip($1, $2, $3)`,
				fakeBidID, test_helpers.FakeIlk.Identifier, blockNumber)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})

		It("is false if deal event in later block", func() {
			fakeBidID := rand.Int()
			blockOne := rand.Int()
			timestampOne := int(rand.Int31())
			blockTwo := blockOne + 1
			timestampTwo := timestampOne + 1000

			blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
			headerOneID, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerOneErr).NotTo(HaveOccurred())

			flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidID)
			test_helpers.CreateFlip(db, blockOneHeader, flipStorageValuesOne,
				test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
			blockTwoHeader.Hash = common.BytesToHash([]byte{5, 4, 3, 2, 1}).String()
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidID)
			test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValuesTwo,
				test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			ilkID, urnID, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidID:           fakeBidID,
					ContractAddress: contractAddress,
					DealRepo:        dealRepo,
					DealHeaderID:    headerTwoID,
				},
				Dealt:            true,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel().ColumnValues["usr"].(string),
				FlipKickRepo:     flipKickRepo,
				FlipKickHeaderID: headerOneID,
			})
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(
				strconv.Itoa(fakeBidID), strconv.FormatInt(ilkID, 10), strconv.FormatInt(urnID, 10), "false",
				blockOneHeader.Timestamp, blockOneHeader.Timestamp, flipStorageValuesOne)

			var actualBid test_helpers.FlipBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip($1, $2, $3)`,
				fakeBidID, test_helpers.FakeIlk.Identifier, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})
	})

	It("gets created and updated blocks", func() {
		fakeBidID := rand.Int()
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1000

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
		headerID, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidID)
		test_helpers.CreateFlip(db, blockOneHeader, flipStorageValuesOne,
			test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

		ilkID, urnID, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerID,
			},
			Dealt:            true,
			IlkHex:           test_helpers.FakeIlk.Hex,
			UrnGuy:           test_data.FlipKickModel().ColumnValues["usr"].(string),
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderID: headerID,
		})
		Expect(err).NotTo(HaveOccurred())

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
		blockTwoHeader.Hash = common.BytesToHash([]byte{5, 4, 3, 2, 1}).String()
		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidID)
		test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValuesTwo,
			test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

		expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidID), strconv.FormatInt(ilkID, 10),
			strconv.FormatInt(urnID, 10), "true", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flipStorageValuesOne)

		var actualBid test_helpers.FlipBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip($1, $2, $3)`,
			fakeBidID, test_helpers.FakeIlk.Identifier, blockTwo)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid.Created).To(Equal(actualBid.Created))
		Expect(expectedBid.Updated).To(Equal(actualBid.Updated))
	})
})
