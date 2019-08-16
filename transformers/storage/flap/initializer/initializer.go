package initializer

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	storage2 "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flap"
)

var StorageTransformerInitializer transformer.StorageTransformerInitializer = storage.Transformer{
	HashedAddress: utils.HexToKeccak256Hash(constants.GetContractAddress("MCD_FLAP")),
	Mappings: &flap.StorageKeysLookup{
		StorageRepository: &storage2.MakerStorageRepository{},
		ContractAddress:   constants.GetContractAddress("MCD_FLAP")},
	Repository: &flap.FlapStorageRepository{ContractAddress: constants.GetContractAddress("MCD_FLAP")},
}.NewTransformer
