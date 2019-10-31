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

package flop

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var (
	BidsIndex = utils.IndexOne

	VatKey      = common.HexToHash(utils.IndexTwo)
	VatMetadata = utils.GetStorageValueMetadata(mcdStorage.Vat, nil, utils.Address)

	GemKey      = common.HexToHash(utils.IndexThree)
	GemMetadata = utils.GetStorageValueMetadata(mcdStorage.Gem, nil, utils.Address)

	BegKey      = common.HexToHash(utils.IndexFour)
	BegMetadata = utils.GetStorageValueMetadata(mcdStorage.Beg, nil, utils.Uint256)

	PadKey      = common.HexToHash(utils.IndexFive)
	PadMetadata = utils.GetStorageValueMetadata(mcdStorage.Pad, nil, utils.Uint256)

	TTLAndTauKey      = common.HexToHash(utils.IndexSix)
	ttlAndTauTypes    = map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48}
	ttlAndTauNames    = map[int]string{0: mcdStorage.TTL, 1: mcdStorage.Tau}
	TTLAndTauMetadata = utils.GetStorageValueMetadataForPackedSlot(mcdStorage.Packed, nil, utils.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksKey      = common.HexToHash(utils.IndexSeven)
	KicksMetadata = utils.GetStorageValueMetadata(mcdStorage.Kicks, nil, utils.Uint256)

	LiveKey      = common.HexToHash(utils.IndexEight)
	LiveMetadata = utils.GetStorageValueMetadata(mcdStorage.Live, nil, utils.Uint256)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{
		storageRepository: storageRepository,
		contractAddress:   contractAddress,
	}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]utils.StorageValueMetadata, error) {
	mappings := loadStaticMappings()
	return loader.loadBidKeys(mappings)
}

func (loader *keysLoader) loadBidKeys(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
	bidIDs, getBidsErr := loader.storageRepository.GetFlopBidIDs(loader.contractAddress)
	if getBidsErr != nil {
		return nil, getBidsErr
	}
	for _, bidID := range bidIDs {
		hexBidID, convertErr := shared.ConvertIntStringToHex(bidID)
		if convertErr != nil {
			return nil, convertErr
		}
		mappings[getBidBidKey(hexBidID)] = getBidBidMetadata(bidID)
		mappings[getBidLotKey(hexBidID)] = getBidLotMetadata(bidID)
		mappings[getBidGuyTicEndKey(hexBidID)] = getBidGuyTicEndMetadata(bidID)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[GemKey] = GemMetadata
	mappings[BegKey] = BegMetadata
	mappings[PadKey] = PadMetadata
	mappings[TTLAndTauKey] = TTLAndTauMetadata
	mappings[KicksKey] = KicksMetadata
	mappings[LiveKey] = LiveMetadata
	return mappings
}

func getBidBidKey(hexBidID string) common.Hash {
	return utils.GetStorageKeyForMapping(BidsIndex, hexBidID)
}

func getBidBidMetadata(bidID string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidID: bidID}
	return utils.GetStorageValueMetadata(mcdStorage.BidBid, keys, utils.Uint256)
}

func getBidLotKey(hexBidID string) common.Hash {
	return utils.GetIncrementedStorageKey(getBidBidKey(hexBidID), 1)
}

func getBidLotMetadata(bidID string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidID: bidID}
	return utils.GetStorageValueMetadata(mcdStorage.BidLot, keys, utils.Uint256)
}

func getBidGuyTicEndKey(hexBidID string) common.Hash {
	return utils.GetIncrementedStorageKey(getBidBidKey(hexBidID), 2)
}

func getBidGuyTicEndMetadata(bidID string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidID: bidID}
	packedTypes := map[int]utils.ValueType{0: utils.Address, 1: utils.Uint48, 2: utils.Uint48}
	packedNames := map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd}
	return utils.GetStorageValueMetadataForPackedSlot(mcdStorage.Packed, keys, utils.PackedSlot, packedNames, packedTypes)
}
