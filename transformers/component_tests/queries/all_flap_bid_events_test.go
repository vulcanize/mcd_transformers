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

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/tend"
	"github.com/vulcanize/mcd_transformers/transformers/events/tick"
	"github.com/vulcanize/mcd_transformers/transformers/events/yank"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flap bid events query", func() {
	var (
		db                     *postgres.DB
		flapKickRepo           flap_kick.FlapKickRepository
		tendRepo               tend.TendRepository
		tickRepo               tick.TickRepository
		dealRepo               deal.DealRepository
		yankRepo               yank.YankRepository
		headerRepo             repositories.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		anotherContractAddress = common.HexToAddress("0xabcdef123456789").Hex()
		blockOne               int64
		headerOne              core.Header
		headerOneID            int64
		headerOneErr           error
		fakeBidID              int
		flapKickEvent          shared.InsertionModel
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)
		tendRepo = tend.TendRepository{}
		tendRepo.SetDB(db)
		tickRepo = tick.TickRepository{}
		tickRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		yankRepo = yank.YankRepository{}
		yankRepo.SetDB(db)
		fakeBidID = rand.Int()

		blockOne = 1
		headerOne = fakes.GetFakeHeader(blockOne)
		headerOneID, headerOneErr = headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(headerOneErr).NotTo(HaveOccurred())
		flapKickLog := test_data.CreateTestLog(headerOneID, db)

		flapKickEvent = test_data.FlapKickModel()
		flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
		flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidID)
		flapKickEvent.ColumnValues[constants.HeaderFK] = headerOneID
		flapKickEvent.ColumnValues[constants.LogFK] = flapKickLog.ID
		flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
		Expect(flapKickErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_flap_bid_events", func() {
		It("returns all flap bid events (same block)", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()
			tendLog := test_data.CreateTestLog(headerOneID, db)

			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				TendRepo:        tendRepo,
				TendHeaderID:    headerOneID,
				TendLogID:       tendLog.ID,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			tickLog := test_data.CreateTestLog(headerOneID, db)
			flapTickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				TickRepo:        tickRepo,
				TickHeaderID:    headerOneID,
				TickLogID:       tickLog.ID,
			})
			Expect(flapTickErr).NotTo(HaveOccurred())

			flapDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerOneID,
			})
			Expect(flapDealErr).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidID)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "tend"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: flapStorageValues[storage.BidBid].(string), Lot: flapStorageValues[storage.BidLot].(string), Act: "tick"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: flapStorageValues[storage.BidBid].(string), Lot: flapStorageValues[storage.BidLot].(string), Act: "deal"},
			))
		})

		It("returns all flap bid events across all blocks", func() {
			fakeBidIDTwo := fakeBidID + 1

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flapKickEventTwoLog := test_data.CreateTestLog(headerTwoID, db)
			flapKickEventTwo := test_data.FlapKickModel()
			flapKickEventTwo.ColumnValues["bid"] = strconv.Itoa(rand.Int())
			flapKickEventTwo.ColumnValues["lot"] = strconv.Itoa(rand.Int())
			flapKickEventTwo.ColumnValues["bid_id"] = strconv.Itoa(fakeBidIDTwo)
			flapKickEventTwo.ColumnValues[constants.HeaderFK] = headerTwoID
			flapKickEventTwo.ColumnValues[constants.LogFK] = flapKickEventTwoLog.ID
			flapKickEventTwo.ForeignKeyValues[constants.AddressFK] = contractAddress
			flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEventTwo})
			Expect(flapKickErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidID:     flapKickEventTwo.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEventTwo.ColumnValues["bid"].(string),
					Lot:       flapKickEventTwo.ColumnValues["lot"].(string),
					Act:       "kick"},
			))
		})

		It("returns bid events for multiple bid ids", func() {
			bidIDOne := fakeBidID
			bidIDTwo := rand.Int()
			lotOne := rand.Int()
			bidAmountOne := rand.Int()

			flapKickEventTwoLog := test_data.CreateTestLog(headerOneID, db)
			flapKickEventTwo := test_data.FlapKickModel()
			flapKickEventTwo.ColumnValues["bid_id"] = strconv.Itoa(bidIDTwo)
			flapKickEventTwo.ColumnValues[constants.HeaderFK] = headerOneID
			flapKickEventTwo.ColumnValues[constants.LogFK] = flapKickEventTwoLog.ID
			flapKickEventTwo.ForeignKeyValues[constants.AddressFK] = contractAddress
			flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEventTwo})
			Expect(flapKickErr).NotTo(HaveOccurred())

			flapTendOneLog := test_data.CreateTestLog(headerOneID, db)
			flapTendOneErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidID:           bidIDOne,
				ContractAddress: contractAddress,
				Lot:             lotOne,
				BidAmount:       bidAmountOne,
				TendRepo:        tendRepo,
				TendHeaderID:    headerOneID,
				TendLogID:       flapTendOneLog.ID,
			})
			Expect(flapTendOneErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidID:     flapKickEventTwo.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEventTwo.ColumnValues["bid"].(string),
					Lot:       flapKickEventTwo.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidID:     strconv.Itoa(bidIDOne),
					BidAmount: strconv.Itoa(bidAmountOne),
					Lot:       strconv.Itoa(lotOne),
					Act:       "tend"},
			))
		})

		Describe("result pagination", func() {
			var (
				bidAmount, lotAmount int
			)

			BeforeEach(func() {
				lotAmount = rand.Int()
				bidAmount = rand.Int()

				headerTwo := fakes.GetFakeHeader(2)
				headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
				Expect(headerTwoErr).NotTo(HaveOccurred())
				tendLogID := test_data.CreateTestLog(headerTwoID, db).ID

				flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					BidID:           fakeBidID,
					ContractAddress: contractAddress,
					Lot:             lotAmount,
					BidAmount:       bidAmount,
					TendRepo:        tendRepo,
					TendHeaderID:    headerTwoID,
					TendLogID:       tendLogID,
				})
				Expect(flapTendErr).NotTo(HaveOccurred())
			})

			It("limits results to latest blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events($1)`,
					maxResults)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{
						BidID:     strconv.Itoa(fakeBidID),
						BidAmount: strconv.Itoa(bidAmount),
						Lot:       strconv.Itoa(lotAmount),
						Act:       "tend",
					},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events($1, $2)`,
					maxResults, resultOffset)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{
						BidID:     flapKickEvent.ColumnValues["bid_id"].(string),
						BidAmount: flapKickEvent.ColumnValues["bid"].(string),
						Lot:       flapKickEvent.ColumnValues["lot"].(string),
						Act:       "kick",
					},
				))
			})
		})

		It("ignores bid events from flops", func() {
			flopKickLog := test_data.CreateTestLog(headerOneID, db)
			flopKickRepo := flop_kick.FlopKickRepository{}
			flopKickRepo.SetDB(db)

			flopKickEvent := test_data.FlopKickModel()
			flopKickEvent.ForeignKeyValues[constants.AddressFK] = "flop"
			flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidID)
			flopKickEvent.ColumnValues[constants.HeaderFK] = headerOneID
			flopKickEvent.ColumnValues[constants.LogFK] = flopKickLog.ID
			flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())
			flopKickBidEvent := test_helpers.BidEvent{
				BidID:           flopKickEvent.ColumnValues["bid_id"].(string),
				BidAmount:       flopKickEvent.ColumnValues["bid"].(string),
				Lot:             flopKickEvent.ColumnValues["lot"].(string),
				Act:             "kick",
				ContractAddress: flopKickEvent.ForeignKeyValues[constants.AddressFK]}

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act, contract_address FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:           flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount:       flapKickEvent.ColumnValues["bid"].(string),
					Lot:             flapKickEvent.ColumnValues["lot"].(string),
					Act:             "kick",
					ContractAddress: flapKickEvent.ForeignKeyValues[constants.AddressFK]},
			))
			Expect(actualBidEvents).NotTo(ContainElement(flopKickBidEvent))
		})
	})

	Describe("tend", func() {
		It("returns flap tend bid events from multiple blocks", func() {
			lot := rand.Int()
			bidAmount := rand.Int()
			updatedLot := lot + 100
			updatedBidAmount := bidAmount + 100
			flapKickBlockOne := flapKickEvent
			flapTendOneLog := test_data.CreateTestLog(headerOneID, db)

			flapTendOneErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				TendRepo:        tendRepo,
				TendHeaderID:    headerOneID,
				TendLogID:       flapTendOneLog.ID,
			})
			Expect(flapTendOneErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flapTendTwoLog := test_data.CreateTestLog(headerTwoID, db)

			flapTendTwoErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				TendRepo:        tendRepo,
				TendHeaderID:    headerTwoID,
				TendLogID:       flapTendTwoLog.ID,
			})
			Expect(flapTendTwoErr).NotTo(HaveOccurred())

			headerThree := fakes.GetFakeHeaderWithTimestamp(int64(333333333), 3)
			headerThreeID, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())
			tendLog := test_data.CreateTestLog(headerThreeID, db)

			// create irrelevant flop tend
			flopTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidID:           fakeBidID,
				ContractAddress: anotherContractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				TendRepo:        tendRepo,
				TendHeaderID:    headerThreeID,
				TendLogID:       tendLog.ID,
			})
			Expect(flopTendErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: strconv.Itoa(bidAmount), Lot: strconv.Itoa(lot), Act: "tend"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: strconv.Itoa(updatedBidAmount), Lot: strconv.Itoa(updatedLot), Act: "tend"},
				test_helpers.BidEvent{
					BidID:     flapKickBlockOne.ColumnValues["bid_id"].(string),
					BidAmount: flapKickBlockOne.ColumnValues["bid"].(string),
					Lot:       flapKickBlockOne.ColumnValues["lot"].(string),
					Act:       "kick"},
			))
		})
	})

	Describe("tick event", func() {
		It("ignores tick events from non flap contracts", func() {
			fakeBidID := rand.Int()
			tickLog := test_data.CreateTestLog(headerOneID, db)

			// irrelevant tick event
			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				BidID:           fakeBidID,
				ContractAddress: "flip",
				TickRepo:        tickRepo,
				TickHeaderID:    headerOneID,
				TickLogID:       tickLog.ID,
			})
			Expect(tickErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick",
				},
			))
		})

		It("includes flap tick bid events", func() {
			fakeBidID := rand.Int()
			tickLog := test_data.CreateTestLog(headerOneID, db)

			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				TickRepo:        tickRepo,
				TickHeaderID:    headerOneID,
				TickLogID:       tickLog.ID,
			})
			Expect(tickErr).NotTo(HaveOccurred())
			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidID)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidID:     strconv.Itoa(fakeBidID),
					BidAmount: flapStorageValues[storage.BidBid].(string),
					Lot:       flapStorageValues[storage.BidLot].(string),
					Act:       "tick"},
			))
		})
	})

	Describe("Deal", func() {
		It("returns bid events with lot and bid amount values from the block where the deal occurred", func() {
			blockTwo := blockOne + 1
			blockThree := blockTwo + 1

			flapKickBlockOne := flapKickEvent

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidID)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			headerTwo := fakes.GetFakeHeader(int64(blockTwo))
			_, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			updatedFlapStorageValues := test_helpers.GetFlapStorageValues(2, fakeBidID)
			test_helpers.CreateFlap(db, headerTwo, updatedFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			headerThree := fakes.GetFakeHeader(int64(blockThree))
			headerThreeID, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())

			flapDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerThreeID,
			})
			Expect(flapDealErr).NotTo(HaveOccurred())

			dealBlockFlapStorageValues := test_helpers.GetFlapStorageValues(0, fakeBidID)
			test_helpers.CreateFlap(db, headerThree, dealBlockFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: dealBlockFlapStorageValues[storage.BidBid].(string), Lot: dealBlockFlapStorageValues[storage.BidLot].(string), Act: "deal"},
				test_helpers.BidEvent{
					BidID:     flapKickBlockOne.ColumnValues["bid_id"].(string),
					BidAmount: flapKickBlockOne.ColumnValues["bid"].(string),
					Lot:       flapKickBlockOne.ColumnValues["lot"].(string),
					Act:       "kick"}))
		})
	})

	Describe("Yank event", func() {
		It("includes yank in all flap bid events", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			tendLog := test_data.CreateTestLog(headerOneID, db)
			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				TendRepo:        tendRepo,
				TendHeaderID:    headerOneID,
				TendLogID:       tendLog.ID,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidID)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flapYankLog := test_data.CreateTestLog(headerTwoID, db)

			flapYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				YankRepo:        yankRepo,
				YankHeaderID:    headerTwoID,
				YankLogID:       flapYankLog.ID,
			})
			Expect(flapYankErr).NotTo(HaveOccurred())

			updatedFlapStorageValues := test_helpers.GetFlapStorageValues(2, fakeBidID)
			test_helpers.CreateFlap(db, headerTwo, updatedFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "tend"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: updatedFlapStorageValues[storage.BidBid].(string), Lot: updatedFlapStorageValues[storage.BidLot].(string), Act: "yank"},
			))
		})

		It("ignores flop yank events", func() {
			fakeBidID := rand.Int()

			flopStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidID)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			yankLog := test_data.CreateTestLog(headerTwoID, db)

			// irrelevant flop yank
			flopYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidID:           fakeBidID,
				ContractAddress: anotherContractAddress,
				YankRepo:        yankRepo,
				YankHeaderID:    headerTwoID,
				YankLogID:       yankLog.ID,
			})
			Expect(flopYankErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
			))
		})
	})
})
