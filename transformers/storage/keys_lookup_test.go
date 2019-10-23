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

package storage_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Storage keys lookup", func() {
	var (
		fakeMetadata = utils.GetStorageValueMetadata("name", map[utils.Key]string{}, utils.Uint256)
		lookup       storage.Mappings
		loader       *test_helpers.MockKeysLoader
	)

	BeforeEach(func() {
		loader = &test_helpers.MockKeysLoader{}
		lookup = mcdStorage.NewKeysLookup(loader)
	})

	Describe("when key not found", func() {
		It("refreshes keys", func() {
			loader.StorageKeyMappings = map[common.Hash]utils.StorageValueMetadata{fakes.FakeHash: fakeMetadata}
			_, err := lookup.Lookup(fakes.FakeHash)

			Expect(err).NotTo(HaveOccurred())
			Expect(loader.LoadMappingsCallCount).To(Equal(1))
		})

		It("returns error if refreshing keys fails", func() {
			loader.LoadMappingsError = fakes.FakeError

			_, err := lookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("when key found", func() {
		BeforeEach(func() {
			loader.StorageKeyMappings = map[common.Hash]utils.StorageValueMetadata{fakes.FakeHash: fakeMetadata}
			_, err := lookup.Lookup(fakes.FakeHash)
			Expect(err).NotTo(HaveOccurred())
			Expect(loader.LoadMappingsCallCount).To(Equal(1))
		})

		It("does not refresh keys", func() {
			_, err := lookup.Lookup(fakes.FakeHash)

			Expect(err).NotTo(HaveOccurred())
			Expect(loader.LoadMappingsCallCount).To(Equal(1))
		})
	})

	It("returns metadata for loaded static key", func() {
		loader.StorageKeyMappings = map[common.Hash]utils.StorageValueMetadata{fakes.FakeHash: fakeMetadata}

		metadata, err := lookup.Lookup(fakes.FakeHash)

		Expect(err).NotTo(HaveOccurred())
		Expect(metadata).To(Equal(fakeMetadata))
	})

	It("returns metadata for hashed version of key (accommodates keys emitted from Geth)", func() {
		loader.StorageKeyMappings = map[common.Hash]utils.StorageValueMetadata{fakes.FakeHash: fakeMetadata}

		hashedKey := common.BytesToHash(crypto.Keccak256(fakes.FakeHash.Bytes()))
		metadata, err := lookup.Lookup(hashedKey)

		Expect(err).NotTo(HaveOccurred())
		Expect(metadata).To(Equal(fakeMetadata))
	})

	It("returns key not found error if key not found", func() {
		_, err := lookup.Lookup(fakes.FakeHash)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
	})
})
