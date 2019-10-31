package flip

import (
	"fmt"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
)

const (
	insertFlipVatQuery   = `INSERT INTO maker.flip_vat (block_number, block_hash, address_id, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipIlkQuery   = `INSERT INTO maker.flip_ilk (block_number, block_hash, address_id, ilk_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipBegQuery   = `INSERT INTO maker.flip_beg (block_number, block_hash, address_id, beg) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipTTLQuery   = `INSERT INTO maker.flip_ttl (block_number, block_hash, address_id, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipTauQuery   = `INSERT INTO maker.flip_tau (block_number, block_hash, address_id, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipKicksQuery = `INSERT INTO maker.flip_kicks (block_number, block_hash, address_id, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	InsertFlipBidBidQuery = `INSERT INTO maker.flip_bid_bid (block_number, block_hash, address_id, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidLotQuery = `INSERT INTO maker.flip_bid_lot (block_number, block_hash, address_id, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidGuyQuery = `INSERT INTO maker.flip_bid_guy (block_number, block_hash, address_id, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidTicQuery = `INSERT INTO maker.flip_bid_tic (block_number, block_hash, address_id, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidEndQuery = `INSERT INTO maker.flip_bid_end (block_number, block_hash, address_id, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidUsrQuery = `INSERT INTO maker.flip_bid_usr (block_number, block_hash, address_id, bid_id, usr) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidGalQuery = `INSERT INTO maker.flip_bid_gal (block_number, block_hash, address_id, bid_id, gal) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidTabQuery = `INSERT INTO maker.flip_bid_tab (block_number, block_hash, address_id, bid_id, tab) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type FlipStorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repository *FlipStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(blockNumber, blockHash, value.(string))
	case storage.Ilk:
		return repository.insertIlk(blockNumber, blockHash, value.(string))
	case storage.Beg:
		return repository.insertBeg(blockNumber, blockHash, value.(string))
	case storage.Kicks:
		return repository.insertKicks(blockNumber, blockHash, value.(string))
	case storage.BidBid:
		return repository.insertBidBid(blockNumber, blockHash, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(blockNumber, blockHash, metadata, value.(string))
	case storage.BidUsr:
		return repository.insertBidUsr(blockNumber, blockHash, metadata, value.(string))
	case storage.BidGal:
		return repository.insertBidGal(blockNumber, blockHash, metadata, value.(string))
	case storage.BidTab:
		return repository.insertBidTab(blockNumber, blockHash, metadata, value.(string))
	case storage.Packed:
		return repository.insertPackedValueRecord(blockNumber, blockHash, metadata, value.(map[int]string))
	default:
		panic(fmt.Sprintf("unrecognized flip contract storage name: %s", metadata.Name))
	}
}

func (repository *FlipStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlipStorageRepository) insertVat(blockNumber int, blockHash string, vat string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlipVatQuery, vat)
}

func (repository *FlipStorageRepository) insertIlk(blockNumber int, blockHash string, ilk string) error {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, repository.db)
	if ilkErr != nil {
		return ilkErr
	}

	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlipIlkQuery, strconv.FormatInt(ilkID, 10))
}

func (repository *FlipStorageRepository) insertBeg(blockNumber int, blockHash string, beg string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlipBegQuery, beg)
}

func (repository *FlipStorageRepository) insertTTL(blockNumber int, blockHash string, ttl string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlipTTLQuery, ttl)
}

func (repository *FlipStorageRepository) insertTau(blockNumber int, blockHash string, tau string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlipTauQuery, tau)
}

func (repository *FlipStorageRepository) insertKicks(blockNumber int, blockHash, kicks string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, InsertFlipKicksQuery, kicks)
}

func (repository *FlipStorageRepository) insertBidBid(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, bid string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, InsertFlipBidBidQuery, bidID, bid)
}

func (repository *FlipStorageRepository) insertBidLot(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, lot string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, InsertFlipBidLotQuery, bidID, lot)
}

func (repository *FlipStorageRepository) insertBidGuy(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, guy string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, InsertFlipBidGuyQuery, bidID, guy)
}

func (repository *FlipStorageRepository) insertBidTic(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, tic string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, InsertFlipBidTicQuery, bidID, tic)
}

func (repository *FlipStorageRepository) insertBidEnd(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, end string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, InsertFlipBidEndQuery, bidID, end)
}

func (repository *FlipStorageRepository) insertBidUsr(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, usr string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, InsertFlipBidUsrQuery, bidID, usr)
}

func (repository *FlipStorageRepository) insertBidGal(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, gal string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, InsertFlipBidGalQuery, bidID, gal)
}

func (repository *FlipStorageRepository) insertBidTab(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, tab string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, InsertFlipBidTabQuery, bidID, tab)
}

func (repository *FlipStorageRepository) insertPackedValueRecord(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case storage.TTL:
			insertErr = repository.insertTTL(blockNumber, blockHash, value)
		case storage.Tau:
			insertErr = repository.insertTau(blockNumber, blockHash, value)
		case storage.BidGuy:
			insertErr = repository.insertBidGuy(blockNumber, blockHash, metadata, value)
		case storage.BidTic:
			insertErr = repository.insertBidTic(blockNumber, blockHash, metadata, value)
		case storage.BidEnd:
			insertErr = repository.insertBidEnd(blockNumber, blockHash, metadata, value)
		default:
			panic(fmt.Sprintf("unrecognized flip contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

func (repository *FlipStorageRepository) insertRecordWithAddress(blockNumber int, blockHash, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flip address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, blockNumber, blockHash, addressID, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flip field with address", insertErr.Error())
		}
		return insertErr
	}

	return tx.Commit()
}

func (repository *FlipStorageRepository) insertRecordWithAddressAndBidID(blockNumber int, blockHash, query, bidID, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	addressID, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flip address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, blockNumber, blockHash, addressID, bidID, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			errorString := fmt.Sprintf("flip field with address for bid id %s", bidID)
			return shared.FormatRollbackError(errorString, insertErr.Error())
		}
		return insertErr
	}
	return tx.Commit()
}

func getBidID(keys map[utils.Key]string) (string, error) {
	bidID, ok := keys[constants.BidID]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.BidID}
	}
	return bidID, nil
}
