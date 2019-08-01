package flap

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	vdbStorage "github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var (
	BidsIndex = vdbStorage.IndexOne

	VatStorageKey = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata   = utils.GetStorageValueMetadata(storage.Vat, nil, utils.Address)

	GemStorageKey = common.HexToHash(vdbStorage.IndexThree)
	GemMetadata   = utils.GetStorageValueMetadata(storage.Gem, nil, utils.Address)

	BegStorageKey = common.HexToHash(vdbStorage.IndexFour)
	BegMetadata   = utils.GetStorageValueMetadata(storage.Beg, nil, utils.Uint256)

	TtlAndTauStorageKey = common.HexToHash(vdbStorage.IndexFive)
	packedTypes         = map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48}
	packedNames         = map[int]string{0: storage.Ttl, 1: storage.Tau}
	TtlAndTauMetadata   = utils.GetStorageValueMetadataForPackedSlot(storage.Packed, nil, utils.PackedSlot, packedNames, packedTypes)

	KicksStorageKey = common.HexToHash(vdbStorage.IndexSix)
	KicksMetadata   = utils.GetStorageValueMetadata(storage.Kicks, nil, utils.Uint256)

	LiveStorageKey = common.HexToHash(vdbStorage.IndexSeven)
	LiveMetadata   = utils.GetStorageValueMetadata(storage.Live, nil, utils.Uint256)
)

type StorageKeysLookup struct {
	StorageRepository storage.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
	ContractAddress   string
}

func (mapping *StorageKeysLookup) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
	metadata, ok := mapping.mappings[key]
	if !ok {
		loadErr := mapping.loadMapping()
		if loadErr != nil {
			return utils.StorageValueMetadata{}, loadErr
		}

		metadata, ok = mapping.mappings[key]
		if !ok {
			return metadata, utils.ErrStorageKeyNotFound{Key: key.Hex()}
		}
	}

	return metadata, nil
}

func (mapping *StorageKeysLookup) SetDB(db *postgres.DB) {
	mapping.StorageRepository.SetDB(db)
}

func (mapping *StorageKeysLookup) loadMapping() error {
	mapping.loadStaticKeys()
	return mapping.loadBidKeys()
}

func (mapping *StorageKeysLookup) loadBidKeys() error {
	bidIds, getBidIdsErr := mapping.StorageRepository.GetFlapBidIds(mapping.ContractAddress)
	for _, bidId := range bidIds {
		hexBidId, convertErr := shared.ConvertIntStringToHex(bidId)
		if convertErr != nil {
			return convertErr
		}

		mapping.mappings[getBidBidKey(hexBidId)] = getBidBidMetadata(bidId)
		mapping.mappings[getBidLotKey(hexBidId)] = getBidLotMetadata(bidId)
		mapping.mappings[getBidGuyKey(hexBidId)] = getBidGuyMetadata(bidId)
		mapping.mappings[getBidTicKey(hexBidId)] = getBidTicMetadata(bidId)
		mapping.mappings[getBidEndKey(hexBidId)] = getBidEndMetadata(bidId)
		mapping.mappings[getBidGalKey(hexBidId)] = getBidGalMetadata(bidId)
	}

	return getBidIdsErr
}

func getBidBidKey(bidId string) common.Hash {
	return vdbStorage.GetMapping(BidsIndex, bidId)
}

func getBidBidMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: storage.BidBid,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Uint256,
	}
}

func getBidLotKey(bidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(bidId), 1) //should this be renamed GetMappingKey?
}

func getBidLotMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: storage.BidLot,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Uint256,
	}
}

func getBidGuyKey(bidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(bidId), 2)
}

func getBidGuyMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: storage.BidGuy,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Address,
	}
}

func getBidTicKey(bidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(bidId), 3)
}

func getBidTicMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: storage.BidTic,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Uint48,
	}
}

func getBidEndKey(bidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(bidId), 4)
}

func getBidEndMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: storage.BidEnd,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Uint48,
	}
}

func getBidGalKey(bidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(bidId), 5)
}

func getBidGalMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: storage.BidGal,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Address,
	}
}

func (mapping *StorageKeysLookup) loadStaticKeys() {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatStorageKey] = VatMetadata
	mappings[GemStorageKey] = GemMetadata
	mappings[BegStorageKey] = BegMetadata
	mappings[TtlAndTauStorageKey] = TtlAndTauMetadata
	mappings[KicksStorageKey] = KicksMetadata
	mappings[LiveStorageKey] = LiveMetadata
	mapping.mappings = mappings
}
