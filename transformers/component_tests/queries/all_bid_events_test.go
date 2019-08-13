package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/dent"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/tend"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
	"time"
)

var _ = Describe("All bid events query", func() {
	var (
		db                  *postgres.DB
		flipKickRepo        flip_kick.FlipKickRepository
		flapKickRepo        flap_kick.FlapKickRepository
		flopKickRepo        flop_kick.FlopKickRepository
		tendRepo            tend.TendRepository
		dentRepo            dent.DentRepository
		headerRepo          repositories.HeaderRepository
		bidId               int
		anotherBidId        int
		flipContractAddress = "flip"
		flapContractAddress = "flap"
		flopContractAddress = "flop"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flipKickRepo = flip_kick.FlipKickRepository{}
		flipKickRepo.SetDB(db)
		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)
		tendRepo = tend.TendRepository{}
		tendRepo.SetDB(db)
		dentRepo = dent.DentRepository{}
		dentRepo.SetDB(db)
		bidId = rand.Intn(100)
		anotherBidId = rand.Intn(100)
		rand.Seed(time.Now().UnixNano())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_bid_events", func() {
		It("returns all bid events for flip, flap, and flop (multiple bid ids, blocks)", func() {
			flipTendLot := rand.Int()
			flipTendBid := rand.Int()
			flipDentLot := rand.Int()
			flipDentBid := rand.Int()

			flapTendLot := rand.Int()
			flapTendBid := rand.Int()
			flapTendLotTwo := rand.Int()
			flapTendBidTwo := rand.Int()

			flopDentLot := rand.Int()
			flopDentBid := rand.Int()

			headerOne := fakes.GetFakeHeader(1)
			headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(headerOneErr).NotTo(HaveOccurred())

			flipKickEvent := test_data.FlipKickModel
			flipKickEvent.ContractAddress = flipContractAddress
			flipKickEvent.BidId = strconv.Itoa(bidId)
			flipKickErr := flipKickRepo.Create(headerOneId, []interface{}{flipKickEvent})
			Expect(flipKickErr).NotTo(HaveOccurred())

			flipTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           bidId,
				ContractAddress: flipContractAddress,
				Lot:             flipTendLot,
				BidAmount:       flipTendBid,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
				TxIndex:         1,
				LogIndex:        2,
			})
			Expect(flipTendErr).NotTo(HaveOccurred())

			flipDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           bidId,
				ContractAddress: flipContractAddress,
				Lot:             flipDentLot,
				BidAmount:       flipDentBid,
				DentRepo:        dentRepo,
				DentHeaderId:    headerOneId,
			})
			Expect(flipDentErr).NotTo(HaveOccurred())

			flapKickEvent := test_data.FlapKickModel
			flapKickEvent.ContractAddress = flapContractAddress
			flapKickEvent.BidId = strconv.Itoa(bidId)
			flapKickErr := flapKickRepo.Create(headerOneId, []interface{}{flapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           bidId,
				ContractAddress: flapContractAddress,
				Lot:             flapTendLot,
				BidAmount:       flapTendBid,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
				TxIndex:         3,
				LogIndex:        4,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			anotherFlipKickEvent := test_data.FlipKickModel
			anotherFlipKickEvent.ContractAddress = flipContractAddress
			anotherFlipKickEvent.BidId = strconv.Itoa(anotherBidId)
			flipKickErr = flipKickRepo.Create(headerTwoId, []interface{}{anotherFlipKickEvent})
			Expect(flipKickErr).NotTo(HaveOccurred())

			flapTendBlockTwoErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           bidId,
				ContractAddress: flapContractAddress,
				Lot:             flapTendLotTwo,
				BidAmount:       flapTendBidTwo,
				TendRepo:        tendRepo,
				TendHeaderId:    headerTwoId,
				TxIndex:         5,
				LogIndex:        6,
			})
			Expect(flapTendBlockTwoErr).NotTo(HaveOccurred())

			flopKickEvent := test_data.FlopKickModel
			flopKickEvent.ContractAddress = flopContractAddress
			flopKickEvent.BidId = strconv.Itoa(bidId)
			flopKickErr := flopKickRepo.Create(headerTwoId, []interface{}{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())

			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           bidId,
				ContractAddress: flopContractAddress,
				Lot:             flopDentLot,
				BidAmount:       flopDentBid,
				DentRepo:        dentRepo,
				DentHeaderId:    headerTwoId,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act, contract_address FROM api.all_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flipKickEvent.BidId, BidAmount: flipKickEvent.Bid, Lot: flipKickEvent.Lot, Act: "kick", ContractAddress: flipContractAddress},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(flipTendBid), Lot: strconv.Itoa(flipTendLot), Act: "tend", ContractAddress: flipContractAddress},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(flipDentBid), Lot: strconv.Itoa(flipDentLot), Act: "dent", ContractAddress: flipContractAddress},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flapKickEvent.Bid, Lot: flapKickEvent.Lot, Act: "kick", ContractAddress: flapContractAddress},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(flapTendBid), Lot: strconv.Itoa(flapTendLot), Act: "tend", ContractAddress: flapContractAddress},
				test_helpers.BidEvent{BidId: anotherFlipKickEvent.BidId, BidAmount: anotherFlipKickEvent.Bid, Lot: anotherFlipKickEvent.Lot, Act: "kick", ContractAddress: flipContractAddress},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(flapTendBidTwo), Lot: strconv.Itoa(flapTendLotTwo), Act: "tend", ContractAddress: flapContractAddress},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick", ContractAddress: flopContractAddress},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(flopDentBid), Lot: strconv.Itoa(flopDentLot), Act: "dent", ContractAddress: flopContractAddress},
			))
		})
	})
})
