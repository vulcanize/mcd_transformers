package queries

import (
	"github.com/vulcanize/mcd_transformers/transformers/events/tick"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/dent"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/yank"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flop bid events query", func() {
	var (
		db              *postgres.DB
		flopKickRepo    flop_kick.FlopKickRepository
		dentRepo        dent.DentRepository
		dealRepo        deal.DealRepository
		yankRepo        yank.YankRepository
		tickRepo        tick.TickRepository
		headerRepo      repositories.HeaderRepository
		blockOne        int64
		headerOne       core.Header
		headerOneID     int64
		headerOneErr    error
		contractAddress string
		fakeBidID       int
		flopKickEvent   shared.InsertionModel
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)
		dentRepo = dent.DentRepository{}
		dentRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		yankRepo = yank.YankRepository{}
		yankRepo.SetDB(db)
		tickRepo = tick.TickRepository{}
		tickRepo.SetDB(db)

		fakeBidID = rand.Int()
		contractAddress = "0x763ztv6x68exwqrgtl325e7hrcvavid4e3fcb4g"

		blockOne = 1
		headerOne = fakes.GetFakeHeader(blockOne)
		headerOneID, headerOneErr = headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(headerOneErr).NotTo(HaveOccurred())
		flopKickLog := test_data.CreateTestLog(headerOneID, db)

		flopKickEvent = test_data.FlopKickModel()
		flopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
		flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidID)
		flopKickEvent.ColumnValues[constants.HeaderFK] = headerOneID
		flopKickEvent.ColumnValues[constants.LogFK] = flopKickLog.ID
		flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})
		Expect(flopKickErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_flop_bid_events", func() {
		It("returns all flop bid events", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			flopDentLog := test_data.CreateTestLog(headerOneID, db)
			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				ContractAddress: contractAddress,
				BidID:           fakeBidID,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				DentRepo:        dentRepo,
				DentHeaderID:    headerOneID,
				DentLogID:       flopDentLog.ID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flopDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerTwoID,
			})
			Expect(flopDealErr).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidID)
			test_helpers.CreateFlop(db, headerTwo, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "dent"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: flopStorageValues[storage.BidBid].(string), Lot: flopStorageValues[storage.BidLot].(string), Act: "deal"},
			))
		})

		It("returns bid events from floppers that have different bid ids", func() {
			lotOne := rand.Int()
			bidAmountOne := rand.Int()

			bidIDTwo := rand.Int()
			lotTwo := rand.Int()
			bidAmountTwo := rand.Int()

			flopKickEventTwoLog := test_data.CreateTestLog(headerOneID, db)

			flopKickEventTwo := test_data.FlopKickModel()
			flopKickEventTwo.ForeignKeyValues[constants.AddressFK] = contractAddress
			flopKickEventTwo.ColumnValues["bid_id"] = strconv.Itoa(bidIDTwo)
			flopKickEventTwo.ColumnValues[constants.HeaderFK] = headerOneID
			flopKickEventTwo.ColumnValues[constants.LogFK] = flopKickEventTwoLog.ID
			flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEventTwo})

			Expect(flopKickErr).NotTo(HaveOccurred())

			flopDentLog := test_data.CreateTestLog(headerOneID, db)
			flopDentOneErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				Lot:             lotOne,
				BidAmount:       bidAmountOne,
				DentRepo:        dentRepo,
				DentHeaderID:    headerOneID,
				DentLogID:       flopDentLog.ID,
			})
			Expect(flopDentOneErr).NotTo(HaveOccurred())

			flopDentTwoLog := test_data.CreateTestLog(headerOneID, db)
			flopDentTwoErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidID:           bidIDTwo,
				ContractAddress: contractAddress,
				Lot:             lotTwo,
				BidAmount:       bidAmountTwo,
				DentRepo:        dentRepo,
				DentHeaderID:    headerOneID,
				DentLogID:       flopDentTwoLog.ID,
			})
			Expect(flopDentTwoErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidID: flopKickEventTwo.ColumnValues["bid_id"].(string), BidAmount: flopKickEventTwo.ColumnValues["bid"].(string), Lot: flopKickEventTwo.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: strconv.Itoa(bidAmountOne), Lot: strconv.Itoa(lotOne), Act: "dent"},
				test_helpers.BidEvent{BidID: strconv.Itoa(bidIDTwo), BidAmount: strconv.Itoa(bidAmountTwo), Lot: strconv.Itoa(lotTwo), Act: "dent"},
			))
		})

		It("ignores bid events from flaps", func() {
			flapKickLog := test_data.CreateTestLog(headerOneID, db)
			flapKickRepo := flap_kick.FlapKickRepository{}
			flapKickRepo.SetDB(db)

			flapKickEvent := test_data.FlapKickModel()
			flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidID)
			flapKickEvent.ColumnValues[constants.HeaderFK] = headerOneID
			flapKickEvent.ColumnValues[constants.LogFK] = flapKickLog.ID
			flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			flapKickBidEvent := test_helpers.BidEvent{
				BidID:     flapKickEvent.ColumnValues["bid_id"].(string),
				BidAmount: flapKickEvent.ColumnValues["bid"].(string),
				Lot:       flapKickEvent.ColumnValues["lot"].(string),
				Act:       "kick"}

			var flopBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&flopBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(flopBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
			))
			Expect(flopBidEvents).NotTo(ContainElement(flapKickBidEvent))
		})
	})

	Describe("dent", func() {
		It("returns flop dent bid events from multiple blocks", func() {
			lot := rand.Int()
			bidAmount := rand.Int()
			updatedLot := lot + 100
			updatedBidAmount := bidAmount + 100

			flopDentHeaderOneLog := test_data.CreateTestLog(headerOneID, db)
			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				DentRepo:        dentRepo,
				DentHeaderID:    headerOneID,
				DentLogID:       flopDentHeaderOneLog.ID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flopDentHeaderTwoLog := test_data.CreateTestLog(headerTwoID, db)

			flopDentHeaderTwoErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				DentRepo:        dentRepo,
				DentHeaderID:    headerTwoID,
				DentLogID:       flopDentHeaderTwoLog.ID,
			})
			Expect(flopDentHeaderTwoErr).NotTo(HaveOccurred())

			headerThree := fakes.GetFakeHeaderWithTimestamp(int64(333333333), 3)
			headerThreeID, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())
			flapDentLog := test_data.CreateTestLog(headerThreeID, db)

			// create irrelevant flap dent
			flapDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidID:           fakeBidID,
				ContractAddress: "flap contract address",
				Lot:             lot,
				BidAmount:       bidAmount,
				DentRepo:        dentRepo,
				DentHeaderID:    headerThreeID,
				DentLogID:       flapDentLog.ID,
			})
			Expect(flapDentErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: strconv.Itoa(bidAmount), Lot: strconv.Itoa(lot), Act: "dent"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: strconv.Itoa(updatedBidAmount), Lot: strconv.Itoa(updatedLot), Act: "dent"},
			))
		})
	})

	Describe("Deal", func() {
		It("returns bid events with lot and bid amount values from the block where the deal occurred", func() {
			fakeBidID := rand.Int()
			blockOne := rand.Int()
			blockTwo := blockOne + 1
			blockThree := blockTwo + 1

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidID)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			headerTwo := fakes.GetFakeHeader(int64(blockTwo))
			_, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			updatedFlopStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidID)
			test_helpers.CreateFlop(db, headerTwo, updatedFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			headerThree := fakes.GetFakeHeader(int64(blockThree))
			headerThreeID, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())

			flopDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerThreeID,
			})
			Expect(flopDealErr).NotTo(HaveOccurred())

			dealBlockFlopStorageValues := test_helpers.GetFlopStorageValues(0, fakeBidID)
			test_helpers.CreateFlop(db, headerThree, dealBlockFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: dealBlockFlopStorageValues[storage.BidBid].(string), Lot: dealBlockFlopStorageValues[storage.BidLot].(string), Act: "deal"},
				test_helpers.BidEvent{BidID: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"}))
		})
	})

	Describe("Yank event", func() {
		It("includes yank in all flop bid events", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			dentLog := test_data.CreateTestLog(headerOneID, db)
			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				DentRepo:        dentRepo,
				DentHeaderID:    headerOneID,
				DentLogID:       dentLog.ID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidID)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flopYankLog := test_data.CreateTestLog(headerOneID, db)

			flopYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				YankRepo:        yankRepo,
				YankHeaderID:    headerTwoID,
				YankLogID:       flopYankLog.ID,
			})
			Expect(flopYankErr).NotTo(HaveOccurred())

			updatedFlopStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidID)
			test_helpers.CreateFlop(db, headerTwo, updatedFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "dent"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: updatedFlopStorageValues[storage.BidBid].(string), Lot: updatedFlopStorageValues[storage.BidLot].(string), Act: "yank"},
			))
		})

		It("ignores flap yank events", func() {
			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidID)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flapYankLog := test_data.CreateTestLog(headerTwoID, db)

			// irrelevant flap yank
			flapYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidID:           fakeBidID,
				ContractAddress: "flap contract address",
				YankRepo:        yankRepo,
				YankHeaderID:    headerTwoID,
				YankLogID:       flapYankLog.ID,
			})
			Expect(flapYankErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
			))
		})
	})

	Describe("tick event", func() {
		It("ignores tick events from non flop contracts", func() {
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
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
			))
		})

		It("includes flop tick bid events", func() {
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
			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidID)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidID: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidID: strconv.Itoa(fakeBidID), BidAmount: flopStorageValues[storage.BidBid].(string), Lot: flopStorageValues[storage.BidLot].(string), Act: "tick"},
			))
		})
	})

	Describe("result pagination", func() {
		var (
			updatedBidAmount, updatedLot int
			flopKickBlockOne             shared.InsertionModel
		)

		BeforeEach(func() {
			lot := rand.Int()
			bidAmount := rand.Int()
			updatedLot = lot + 100
			updatedBidAmount = bidAmount + 100

			logID := test_data.CreateTestLog(headerOneID, db).ID
			flopKickBlockOne = test_data.FlopKickModel()
			flopKickBlockOne.ColumnValues["bid_id"] = strconv.Itoa(fakeBidID)
			flopKickBlockOne.ForeignKeyValues[constants.AddressFK] = contractAddress
			flopKickBlockOne.ColumnValues[constants.HeaderFK] = headerOneID
			flopKickBlockOne.ColumnValues[constants.LogFK] = logID
			flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickBlockOne})
			Expect(flopKickErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			logTwoID := test_data.CreateTestLog(headerTwoID, db).ID

			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidID:           fakeBidID,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				DentRepo:        dentRepo,
				DentHeaderID:    headerTwoID,
				DentLogID:       logTwoID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())
		})

		It("limits result to latest blocks if max_results argument is provided", func() {
			maxResults := 1
			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events($1)`,
				maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:     strconv.Itoa(fakeBidID),
					BidAmount: strconv.Itoa(updatedBidAmount),
					Lot:       strconv.Itoa(updatedLot),
					Act:       "dent",
				},
			))
		})

		It("offsets results if offset is provided", func() {
			maxResults := 1
			resultOffset := 1
			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events($1, $2)`,
				maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidID:     flopKickBlockOne.ColumnValues["bid_id"].(string),
					BidAmount: flopKickBlockOne.ColumnValues["bid"].(string),
					Lot:       flopKickBlockOne.ColumnValues["lot"].(string),
					Act:       "kick",
				},
			))
		})
	})
})
