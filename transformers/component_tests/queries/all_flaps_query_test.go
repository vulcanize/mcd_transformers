package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("All flaps query", func() {
	var (
		db              *postgres.DB
		flapKickRepo    flap_kick.FlapKickRepository
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

		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the most recent flap for every bid id", func() {
		fakeBidIDOne := rand.Int()
		fakeBidIDTwo := fakeBidIDOne + 1

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(blockOneTimestamp, int64(blockOne))
		headerOneID, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(blockTwoTimestamp, int64(blockTwo))
		blockTwoHeader.Hash = "blockTwoHeader"
		headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		contextErr := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidIDOne,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlapKickRepo:     flapKickRepo,
			FlapKickHeaderID: headerOneID,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidIDOne)
		test_helpers.CreateFlap(db, blockOneHeader, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIDOne)), contractAddress)

		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidIDOne)
		test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIDOne)), contractAddress)

		contextErr = test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidID:           fakeBidIDTwo,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlapKickRepo:     flapKickRepo,
			FlapKickHeaderID: headerTwoID,
		})
		Expect(contextErr).NotTo(HaveOccurred())
		flapStorageValuesThree := test_helpers.GetFlapStorageValues(3, fakeBidIDTwo)
		test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesThree, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIDTwo)), contractAddress)

		var actualBids []test_helpers.FlapBid
		queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flaps()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedBidOne := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIDOne), "false", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flapStorageValuesTwo)
		expectedBidTwo := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIDTwo), "false", blockTwoHeader.Timestamp, blockTwoHeader.Timestamp, flapStorageValuesThree)

		Expect(len(actualBids)).To(Equal(2))
		Expect(actualBids).To(ConsistOf([]test_helpers.FlapBid{
			expectedBidOne,
			expectedBidTwo,
		}))
	})

	Describe("result pagination", func() {
		var (
			fakeBidIDOne, fakeBidIDTwo                 int
			blockOneHeader, blockTwoHeader             core.Header
			flapStorageValuesOne, flapStorageValuesTwo map[string]interface{}
		)

		BeforeEach(func() {
			fakeBidIDOne = rand.Int()
			fakeBidIDTwo = fakeBidIDOne + 1

			blockOneHeader = fakes.GetFakeHeaderWithTimestamp(blockOneTimestamp, int64(blockOne))
			headerOneID, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerOneErr).NotTo(HaveOccurred())

			contextErr := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidID:           fakeBidIDOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlapKickRepo:     flapKickRepo,
				FlapKickHeaderID: headerOneID,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			flapStorageValuesOne = test_helpers.GetFlapStorageValues(2, fakeBidIDOne)
			test_helpers.CreateFlap(db, blockOneHeader, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIDOne)), contractAddress)

			blockTwoHeader = fakes.GetFakeHeaderWithTimestamp(blockTwoTimestamp, int64(blockTwo))
			blockTwoHeader.Hash = "blockTwoHeader"
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			contextErr = test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidID:           fakeBidIDTwo,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlapKickRepo:     flapKickRepo,
				FlapKickHeaderID: headerTwoID,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			flapStorageValuesTwo = test_helpers.GetFlapStorageValues(3, fakeBidIDTwo)
			test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIDTwo)), contractAddress)
		})

		It("limits results if max_results argument is provided", func() {
			maxResults := 1
			var actualBids []test_helpers.FlapBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flaps($1)`, maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIDTwo), "false",
				blockTwoHeader.Timestamp, blockTwoHeader.Timestamp, flapStorageValuesTwo)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})

		It("offsets results if offset is provided", func() {
			maxResults := 1
			resultOffset := 1
			var actualBids []test_helpers.FlapBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flaps($1, $2)`,
				maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIDOne), "false",
				blockOneHeader.Timestamp, blockOneHeader.Timestamp, flapStorageValuesOne)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})
	})
})
