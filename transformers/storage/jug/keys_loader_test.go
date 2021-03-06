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

package jug_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
)

var _ = Describe("jug storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = jug.NewKeysLoader(storageRepository)
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[jug.VatKey]).To(Equal(jug.VatMetadata))
		Expect(mappings[jug.VowKey]).To(Equal(jug.VowMetadata))
		Expect(mappings[jug.BaseKey]).To(Equal(jug.BaseMetadata))
	})

	Describe("when getting ilks fails", func() {
		It("returns error", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := storageKeysLoader.LoadMappings()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("when getting ilks succeeds", func() {
		It("returns value metadata for tax", func() {
			storageRepository.Ilks = []string{test_helpers.FakeIlk}
			ilkTaxKey := common.BytesToHash(crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + jug.IlkMappingIndex)))
			expectedMetadata := utils.StorageValueMetadata{
				Name: jug.IlkDuty,
				Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
				Type: utils.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()

			Expect(err).NotTo(HaveOccurred())
			Expect(mappings[ilkTaxKey]).To(Equal(expectedMetadata))
		})

		It("returns value metadata for rho", func() {
			storageRepository.Ilks = []string{test_helpers.FakeIlk}
			ilkTaxKeyBytes := crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + jug.IlkMappingIndex))
			ilkTaxAsInt := big.NewInt(0).SetBytes(ilkTaxKeyBytes)
			incrementedIlkTax := big.NewInt(0).Add(ilkTaxAsInt, big.NewInt(1))
			ilkRhoKey := common.BytesToHash(incrementedIlkTax.Bytes())
			expectedMetadata := utils.StorageValueMetadata{
				Name: jug.IlkRho,
				Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
				Type: utils.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()

			Expect(err).NotTo(HaveOccurred())
			Expect(mappings[ilkRhoKey]).To(Equal(expectedMetadata))
		})
	})
})
