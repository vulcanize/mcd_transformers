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

package vat_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
)

var _ = Describe("Vat storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		mappings          vat.VatMappings
		fakeDiff          utils.StorageDiff
		keccakOfAddress = crypto.Keccak256Hash(fakes.FakeAddress[:])
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		mappings = vat.VatMappings{StorageRepository: storageRepository}
		fakeDiff = utils.StorageDiff{
			Contract:                fakes.FakeAddress,
			StorageKey:              fakes.FakeHash,
		}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			debtDiff := utils.StorageDiff{
				Contract:                fakes.FakeAddress,
				StorageKey:              vat.DebtKey,
			}
			viceDiff := utils.StorageDiff{
				Contract:                fakes.FakeAddress,
				StorageKey:              vat.ViceKey,
			}
			lineDiff := utils.StorageDiff{
				Contract:                fakes.FakeAddress,
				StorageKey:              vat.LineKey,
			}
			liveDiff := utils.StorageDiff{
				Contract:                fakes.FakeAddress,
				StorageKey:              vat.LiveKey,
			}
			Expect(mappings.Lookup(debtDiff)).To(Equal(vat.DebtMetadata))
			Expect(mappings.Lookup(viceDiff)).To(Equal(vat.ViceMetadata))
			Expect(mappings.Lookup(lineDiff)).To(Equal(vat.LineMetadata))
			Expect(mappings.Lookup(liveDiff)).To(Equal(vat.LiveMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := mappings.Lookup(fakeDiff)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			mappings.Lookup(fakeDiff)

			Expect(storageRepository.GetDaiKeysCalled).To(BeTrue())
			Expect(storageRepository.GetGemKeysCalled).To(BeTrue())
			Expect(storageRepository.GetIlksCalled).To(BeTrue())
			Expect(storageRepository.GetVatSinKeysCalled).To(BeTrue())
			Expect(storageRepository.GetUrnsCalled).To(BeTrue())
		})

		It("returns error if dai keys lookup fails", func() {
			storageRepository.GetDaiKeysError = fakes.FakeError

			_, err := mappings.Lookup(fakeDiff)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if gem keys lookup fails", func() {
			storageRepository.GetGemKeysError = fakes.FakeError

			_, err := mappings.Lookup(fakeDiff)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if ilks lookup fails", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := mappings.Lookup(fakeDiff)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if sin keys lookup fails", func() {
			storageRepository.GetVatSinKeysError = fakes.FakeError

			_, err := mappings.Lookup(fakeDiff)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if urns lookup fails", func() {
			storageRepository.GetUrnsError = fakes.FakeError

			_, err := mappings.Lookup(fakeDiff)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if lookups return addresses not of length 42", func() {
			storageRepository.DaiKeys = []string{"0xshortAddress"}

			_, err := mappings.Lookup(fakeDiff)

			Expect(err).To(HaveOccurred())
		})

		Describe("ilk", func() {
			var (
				expectedIlkArtMetadata = utils.StorageValueMetadata{
					Name: vat.IlkArt,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}
				expectedIlkRateMetadata = utils.StorageValueMetadata{
					Name: vat.IlkRate,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}
				expectedIlkSpotMetadata = utils.StorageValueMetadata{
					Name: vat.IlkSpot,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}
				expectedIlkLineMetadata = utils.StorageValueMetadata{
					Name: vat.IlkLine,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}
				expectedIlkDustMetadata = utils.StorageValueMetadata{
					Name: vat.IlkDust,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}
				ilkArtKey common.Hash
				ilkArtAsInt *big.Int
				ilkRateKey common.Hash
				ilkSpotKey common.Hash
				ilkLineKey common.Hash
				ilkDustKey common.Hash
			)

			BeforeEach(func() {
				storageRepository.Ilks = []string{test_helpers.FakeIlk}
				ilkArtKey = common.BytesToHash(crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + vat.IlksMappingIndex)))
				ilkArtAsInt = big.NewInt(0).SetBytes(ilkArtKey.Bytes())

				incrementedIlkArtBy1 := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(1))
				ilkRateKey = common.BytesToHash(incrementedIlkArtBy1.Bytes())

				incrementedIlkArtBy2 := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(2))
				ilkSpotKey = common.BytesToHash(incrementedIlkArtBy2.Bytes())

				incrementedIlkArtBy3 := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(3))
				ilkLineKey = common.BytesToHash(incrementedIlkArtBy3.Bytes())

				incrementedIlkArtBy4 := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(4))
				ilkDustKey = common.BytesToHash(incrementedIlkArtBy4.Bytes())
			})

			It("returns value metadata for ilk Art for csv diff", func() {
				ilkArtDiff := utils.StorageDiff{
					Contract:                fakes.FakeAddress,
					StorageKey:              ilkArtKey,
				}

				Expect(mappings.Lookup(ilkArtDiff)).To(Equal(expectedIlkArtMetadata))
			})

			It("returns value metadata for ilk Art for geth diff", func() {
				keccakOfIlkArtKey := crypto.Keccak256Hash(ilkArtKey[:])
				ilkArtDiff := utils.StorageDiff{
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakOfIlkArtKey,
				}

				Expect(mappings.Lookup(ilkArtDiff)).To(Equal(expectedIlkArtMetadata))
			})

			It("returns value metadata for ilk rate for csv", func() {
				ilkRateDiff := utils.StorageDiff{
					Contract:                fakes.FakeAddress,
					StorageKey:              ilkRateKey,
				}

				Expect(mappings.Lookup(ilkRateDiff)).To(Equal(expectedIlkRateMetadata))
			})

			It("returns value metadata for ilk rate for geth diff ", func() {
				keccakOfIlkRateKey := crypto.Keccak256Hash(ilkRateKey[:])
				ilkRateDiff := utils.StorageDiff{
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakOfIlkRateKey,
				}

				Expect(mappings.Lookup(ilkRateDiff)).To(Equal(expectedIlkRateMetadata))
			})

			It("returns value metadata for ilk spot for csv diff", func() {
				ilkSpotDiff := utils.StorageDiff{
					Contract:                fakes.FakeAddress,
					StorageKey:              ilkSpotKey,
				}

				Expect(mappings.Lookup(ilkSpotDiff)).To(Equal(expectedIlkSpotMetadata))
			})

			It("returns value metadata for ilk spot for geth diff", func() {
				keccakIlkSpotKey := crypto.Keccak256Hash(ilkSpotKey[:])
				ilkSpotDiff := utils.StorageDiff{
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakIlkSpotKey,
				}

				Expect(mappings.Lookup(ilkSpotDiff)).To(Equal(expectedIlkSpotMetadata))
			})

			It("returns value metadata for ilk line for csv diff", func() {
				ilkLineDiff := utils.StorageDiff{
					Contract:                fakes.FakeAddress,
					StorageKey:              ilkLineKey,
				}

				Expect(mappings.Lookup(ilkLineDiff)).To(Equal(expectedIlkLineMetadata))
			})

			It("returns value metadata for ilk line for geth diff", func() {
			    keccakOfilkLineKey := crypto.Keccak256Hash(ilkLineKey[:])
				ilkLineDiff := utils.StorageDiff{
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakOfilkLineKey,
				}

				Expect(mappings.Lookup(ilkLineDiff)).To(Equal(expectedIlkLineMetadata))
			})

			It("returns value metadata for ilk dust for csv diff", func() {
				ilkDustDiff := utils.StorageDiff{
					Contract:                fakes.FakeAddress,
					StorageKey:              ilkDustKey,
				}

				Expect(mappings.Lookup(ilkDustDiff)).To(Equal(expectedIlkDustMetadata))
			})

			It("returns value metadata for ilk dust for geth diff", func() {
				keccakOfIlkDustKey := crypto.Keccak256Hash(ilkDustKey[:])
				ilkDustDiff := utils.StorageDiff{
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakOfIlkDustKey,
				}

				Expect(mappings.Lookup(ilkDustDiff)).To(Equal(expectedIlkDustMetadata))
			})
		})

		Describe("urn", func() {
			var (
				expectedUrnInkMetadata = utils.StorageValueMetadata{
					Name: vat.UrnInk,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk, constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}
				expectedUrnArtMetadata = utils.StorageValueMetadata{
					Name: vat.UrnArt,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk, constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}
				encodedSecondaryMapIndex []byte
				urnInkKey common.Hash
				urnArtKey common.Hash
			)

			BeforeEach(func() {
				storageRepository.Urns = []storage.Urn{{Ilk: test_helpers.FakeIlk, Identifier: test_helpers.FakeAddress}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + vat.UrnsMappingIndex))
				paddedUrnGuy := common.FromHex("0x000000000000000000000000" + test_helpers.FakeAddress[2:])
				encodedSecondaryMapIndex = crypto.Keccak256(paddedUrnGuy, encodedPrimaryMapIndex)
				urnInkKey = common.BytesToHash(encodedSecondaryMapIndex)
				urnInkAsInt := big.NewInt(0).SetBytes(encodedSecondaryMapIndex)
				incrementedUrnInk := big.NewInt(0).Add(urnInkAsInt, big.NewInt(1))
				urnArtKey = common.BytesToHash(incrementedUrnInk.Bytes())
			})
			It("returns value metadata for urn ink for a csv diff", func() {
				urnInkDiff := utils.StorageDiff{
					Contract:                fakes.FakeAddress,
					StorageKey:              urnInkKey,
				}

				Expect(mappings.Lookup(urnInkDiff)).To(Equal(expectedUrnInkMetadata))
			})

			It("returns value metadata for urn ink for a geth diff", func() {
				keccakOfUrnInkKey := crypto.Keccak256Hash(urnInkKey[:])
				urnInkDiff := utils.StorageDiff{
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakOfUrnInkKey,
				}

				Expect(mappings.Lookup(urnInkDiff)).To(Equal(expectedUrnInkMetadata))
			})

			It("returns value metadata for urn art for a csv diff", func() {
				urnArtDiff := utils.StorageDiff{
					Contract:                fakes.FakeAddress,
					StorageKey:              urnArtKey,
				}

				Expect(mappings.Lookup(urnArtDiff)).To(Equal(expectedUrnArtMetadata))
			})

			It("returns value metadata for urn art for a geth diff", func() {
				keccakOfUrnArtKey := crypto.Keccak256Hash(urnArtKey[:])
				urnArtDiff := utils.StorageDiff{
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakOfUrnArtKey,
				}

				Expect(mappings.Lookup(urnArtDiff)).To(Equal(expectedUrnArtMetadata))
			})
		})

		Describe("gem", func() {
			var (
				expectedMetadata = utils.StorageValueMetadata{
					Name: vat.Gem,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk, constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}
				gemKey common.Hash
			)

			BeforeEach(func() {
				storageRepository.GemKeys = []storage.Urn{{Ilk: test_helpers.FakeIlk, Identifier: test_helpers.FakeAddress}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + vat.GemsMappingIndex))
				paddedGemAddress := common.FromHex("0x000000000000000000000000" + test_helpers.FakeAddress[2:])
				encodedSecondaryMapIndex := crypto.Keccak256(paddedGemAddress, encodedPrimaryMapIndex)
				gemKey = common.BytesToHash(encodedSecondaryMapIndex)
			})

			It("returns value metadata for gem", func() {
				gemDiff := utils.StorageDiff{
					Contract:                fakes.FakeAddress,
					StorageKey:              gemKey,
				}

				Expect(mappings.Lookup(gemDiff)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for gem", func() {
				keccakOfGemKey := crypto.Keccak256Hash(gemKey[:])
				gemDiff := utils.StorageDiff{
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakOfGemKey,
				}

				Expect(mappings.Lookup(gemDiff)).To(Equal(expectedMetadata))
			})
		})

		Describe("dai", func() {
			var (
				expectedMetadata = utils.StorageValueMetadata{
					Name: vat.Dai,
					Keys: map[utils.Key]string{constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}
				daiKey common.Hash
			)

			BeforeEach(func() {
				storageRepository.DaiKeys = []string{test_helpers.FakeAddress}
				paddedDaiAddress := "0x000000000000000000000000" + test_helpers.FakeAddress[2:]
				daiKey = common.BytesToHash(crypto.Keccak256(common.FromHex(paddedDaiAddress + vat.DaiMappingIndex)))
			})
			It("returns value metadata for dai for a csv diff", func() {
				daiDiff := utils.StorageDiff{
					Contract:                fakes.FakeAddress,
					StorageKey:              daiKey,
				}

				Expect(mappings.Lookup(daiDiff)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for dai for a geth diff", func() {
				keccakOfDaiKey := crypto.Keccak256Hash(daiKey[:])
				daiDiff := utils.StorageDiff{
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakOfDaiKey,
				}

				Expect(mappings.Lookup(daiDiff)).To(Equal(expectedMetadata))
			})
		})

		Describe("sin", func() {
			var (
				expectedMetadata = utils.StorageValueMetadata{
					Name: vat.Sin,
					Keys: map[utils.Key]string{constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}
				sinKey common.Hash
			)
			BeforeEach(func() {
				storageRepository.SinKeys = []string{test_helpers.FakeAddress}
				paddedSinAddress := "0x000000000000000000000000" + test_helpers.FakeAddress[2:]
				sinKey = common.BytesToHash(crypto.Keccak256(common.FromHex(paddedSinAddress + vat.SinMappingIndex)))
			})

			It("returns value metadata for sin for a csv diff", func() {
				sinDiff := utils.StorageDiff{
					Contract:   fakes.FakeAddress,
					StorageKey: sinKey,
				}

				Expect(mappings.Lookup(sinDiff)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for sin for a geth diff", func() {
				keccakOfSinKey := crypto.Keccak256Hash(sinKey[:])
				sinDiff := utils.StorageDiff{
					//this is how it should be, but in the geth patch, this is being converted to an address
					// TODO: fix this in the geth patch
					KeccakOfContractAddress: keccakOfAddress,
					StorageKey:              keccakOfSinKey,
				}

				Expect(mappings.Lookup(sinDiff)).To(Equal(expectedMetadata))
			})
		})
	})
})
