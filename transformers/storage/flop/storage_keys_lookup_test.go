package flop_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	vdbStorage "github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Flop storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		mappings          flop.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		mappings = flop.StorageKeysLookup{StorageRepository: storageRepository, ContractAddress: "0x668001c75a9c02d6b10c7a17dbd8aa4afff95037"}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(mappings.Lookup(flop.VatKey)).To(Equal(flop.VatMetadata))
			Expect(mappings.Lookup(flop.GemKey)).To(Equal(flop.GemMetadata))
			Expect(mappings.Lookup(flop.BegKey)).To(Equal(flop.BegMetadata))
			Expect(mappings.Lookup(flop.TtlAndTauKey)).To(Equal(flop.TtlAndTauMetadata))
			Expect(mappings.Lookup(flop.KicksKey)).To(Equal(flop.KicksMetadata))
			Expect(mappings.Lookup(flop.LiveKey)).To(Equal(flop.LiveMetadata))
		})

		It("returns value metadata if keccak hash of key exists", func() {
			Expect(mappings.Lookup(crypto.Keccak256Hash(flop.VatKey[:]))).To(Equal(flop.VatMetadata))
			Expect(mappings.Lookup(crypto.Keccak256Hash(flop.GemKey[:]))).To(Equal(flop.GemMetadata))
			Expect(mappings.Lookup(crypto.Keccak256Hash(flop.BegKey[:]))).To(Equal(flop.BegMetadata))
			Expect(mappings.Lookup(crypto.Keccak256Hash(flop.TtlAndTauKey[:]))).To(Equal(flop.TtlAndTauMetadata))
			Expect(mappings.Lookup(crypto.Keccak256Hash(flop.KicksKey[:]))).To(Equal(flop.KicksMetadata))
			Expect(mappings.Lookup(crypto.Keccak256Hash(flop.LiveKey[:]))).To(Equal(flop.LiveMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			_, _ = mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetFlopBidIdsCalledWith).To(Equal(mappings.ContractAddress))
		})

		It("returns error if bid ID lookup fails", func() {
			storageRepository.GetFlopBidIdsError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("bid", func() {
		var fakeBidId string
		var bidBidKey common.Hash

		BeforeEach(func() {
			fakeBidId = "42"
			fakeHexBidId, conversionErr := shared.ConvertIntStringToHex(fakeBidId)

			Expect(conversionErr).NotTo(HaveOccurred())

			bidBidKey = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeHexBidId + flop.BidsIndex)))
			storageRepository.FlopBidIds = []string{fakeBidId}
		})

		It("returns value metadata for bid bid", func() {
			expectedMetadata := utils.StorageValueMetadata{
				Name: storage.BidBid,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			Expect(mappings.Lookup(bidBidKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for bid lot", func() {
			bidLotKey := vdbStorage.GetIncrementedKey(bidBidKey, 1)
			expectedMetadata := utils.StorageValueMetadata{
				Name: storage.BidLot,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			Expect(mappings.Lookup(bidLotKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for bid guy + tic + end packed slot", func() {
			bidGuyKey := vdbStorage.GetIncrementedKey(bidBidKey, 2)
			expectedMetadata := utils.StorageValueMetadata{
				Name:        storage.Packed,
				Keys:        map[utils.Key]string{constants.BidId: fakeBidId},
				Type:        utils.PackedSlot,
				PackedTypes: map[int]utils.ValueType{0: utils.Address, 1: utils.Uint48, 2: utils.Uint48},
				PackedNames: map[int]string{0: storage.BidGuy, 1: storage.BidTic, 2: storage.BidEnd},
			}
			Expect(mappings.Lookup(bidGuyKey)).To(Equal(expectedMetadata))
		})
	})
})
