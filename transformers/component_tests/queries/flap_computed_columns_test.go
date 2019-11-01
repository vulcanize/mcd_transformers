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
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/tend"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flap computed columns", func() {
	var (
		db              *postgres.DB
		flapKickRepo    flap_kick.FlapKickRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = fakes.FakeAddress.Hex()

		fakeBidID      = rand.Int()
		blockOne       = rand.Int()
		blockOneHeader = fakes.GetFakeHeader(int64(blockOne))
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("flap_bid_events", func() {
		It("returns the bid events for a flap", func() {
			headerID, headerErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerErr).NotTo(HaveOccurred())
			flapKickLog := test_data.CreateTestLog(headerID, db)

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidID)
			test_helpers.CreateFlap(db, blockOneHeader, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			flapKickEvent := test_data.FlapKickModel()
			flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidID)
			flapKickEvent.ColumnValues[constants.HeaderFK] = headerID
			flapKickEvent.ColumnValues[constants.LogFK] = flapKickLog.ID
			flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			expectedBidEvents := test_helpers.BidEvent{
				BidID:           strconv.Itoa(fakeBidID),
				Lot:             flapKickEvent.ColumnValues["lot"].(string),
				BidAmount:       flapKickEvent.ColumnValues["bid"].(string),
				Act:             "kick",
				ContractAddress: contractAddress,
			}
			var actualBidEvents test_helpers.BidEvent
			queryErr := db.Get(&actualBidEvents,
				`SELECT bid_id, bid_amount, lot, act, contract_address FROM api.flap_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flap_state
    					FROM api.all_flaps()))`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(Equal(expectedBidEvents))
		})

		It("does not include bid events for a different flap", func() {
			headerID, headerErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerErr).NotTo(HaveOccurred())
			flapKickLog := test_data.CreateTestLog(headerID, db)

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidID)
			test_helpers.CreateFlap(db, blockOneHeader, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			flapKickEvent := test_data.FlapKickModel()
			flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidID)
			flapKickEvent.ColumnValues[constants.HeaderFK] = headerID
			flapKickEvent.ColumnValues[constants.LogFK] = flapKickLog.ID
			flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			blockTwo := blockOne + 1
			blockTwoHeader := fakes.GetFakeHeader(int64(blockTwo))
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			irrelevantFlipKickLog := test_data.CreateTestLog(headerTwoID, db)

			irrelevantBidID := fakeBidID + 9999999999999
			irrelevantFlapStorageValues := test_helpers.GetFlapStorageValues(2, irrelevantBidID)
			test_helpers.CreateFlap(db, blockTwoHeader, irrelevantFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(irrelevantBidID)), contractAddress)

			irrelevantFlapKickEvent := test_data.FlapKickModel()
			irrelevantFlapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			irrelevantFlapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(irrelevantBidID)
			irrelevantFlapKickEvent.ColumnValues[constants.HeaderFK] = headerTwoID
			irrelevantFlapKickEvent.ColumnValues[constants.LogFK] = irrelevantFlipKickLog.ID

			flapKickErr = flapKickRepo.Create([]shared.InsertionModel{irrelevantFlapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			expectedBidEvents := test_helpers.BidEvent{
				BidID:           strconv.Itoa(fakeBidID),
				Lot:             flapKickEvent.ColumnValues["lot"].(string),
				BidAmount:       flapKickEvent.ColumnValues["bid"].(string),
				Act:             "kick",
				ContractAddress: contractAddress,
			}

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents,
				`SELECT bid_id, bid_amount, lot, act, contract_address FROM api.flap_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flap_state
    					FROM api.all_flaps() WHERE bid_id = $1))`, fakeBidID)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(ConsistOf(expectedBidEvents))
		})

		Describe("result pagination", func() {
			var (
				tendBid, tendLot int
				flapKickEvent    shared.InsertionModel
			)

			BeforeEach(func() {
				headerID, headerErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
				Expect(headerErr).NotTo(HaveOccurred())
				logID := test_data.CreateTestLog(headerID, db).ID

				flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidID)
				test_helpers.CreateFlap(db, blockOneHeader, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

				flapKickEvent = test_data.FlapKickModel()
				flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
				flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidID)
				flapKickEvent.ColumnValues[constants.HeaderFK] = headerID
				flapKickEvent.ColumnValues[constants.LogFK] = logID
				flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
				Expect(flapKickErr).NotTo(HaveOccurred())

				blockTwo := blockOne + 1
				blockTwoHeader := fakes.GetFakeHeader(int64(blockTwo))
				headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
				Expect(headerTwoErr).NotTo(HaveOccurred())
				logTwoID := test_data.CreateTestLog(headerTwoID, db).ID

				tendBid = rand.Int()
				tendLot = rand.Int()
				tendRepo := tend.TendRepository{}
				tendRepo.SetDB(db)
				flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					BidID:           fakeBidID,
					ContractAddress: contractAddress,
					Lot:             tendLot,
					BidAmount:       tendBid,
					TendRepo:        tendRepo,
					TendHeaderID:    headerTwoID,
					TendLogID:       logTwoID,
				})
				Expect(flapTendErr).NotTo(HaveOccurred())
			})

			It("limits result to most recent block if max_results argument is provided", func() {
				expectedBidEvent := test_helpers.BidEvent{
					BidID:           strconv.Itoa(fakeBidID),
					Lot:             strconv.Itoa(tendLot),
					BidAmount:       strconv.Itoa(tendBid),
					Act:             "tend",
					ContractAddress: contractAddress,
				}

				maxResults := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents,
					`SELECT bid_id, bid_amount, lot, act, contract_address FROM api.flap_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flap_state
    					FROM api.all_flaps() WHERE bid_id = $1), $2)`, fakeBidID, maxResults)

				Expect(queryErr).NotTo(HaveOccurred())
				Expect(actualBidEvents).To(ConsistOf(expectedBidEvent))
			})

			It("offsets results if offset is provided", func() {
				expectedBidEvent := test_helpers.BidEvent{
					BidID:           strconv.Itoa(fakeBidID),
					Lot:             flapKickEvent.ColumnValues["lot"].(string),
					BidAmount:       flapKickEvent.ColumnValues["bid"].(string),
					Act:             "kick",
					ContractAddress: contractAddress,
				}

				maxResults := 1
				resultOffset := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents,
					`SELECT bid_id, bid_amount, lot, act, contract_address FROM api.flap_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flap_state
    					FROM api.all_flaps() WHERE bid_id = $1), $2, $3)`, fakeBidID, maxResults, resultOffset)

				Expect(queryErr).NotTo(HaveOccurred())
				Expect(actualBidEvents).To(ConsistOf(expectedBidEvent))
			})
		})
	})
})
