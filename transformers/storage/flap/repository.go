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

package flap

import (
	"fmt"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery    = `INSERT INTO maker.flap_vat (block_number, block_hash, address_id, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertGemQuery    = `INSERT INTO maker.flap_gem (block_number, block_hash, address_id, gem) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBegQuery    = `INSERT INTO maker.flap_beg (block_number, block_hash, address_id, beg) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertTTLQuery    = `INSERT INTO maker.flap_ttl (block_number, block_hash, address_id, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertTauQuery    = `INSERT INTO maker.flap_tau (block_number, block_hash, address_id, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertKicksQuery  = `INSERT INTO maker.flap_kicks (block_number, block_hash, address_id, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertLiveQuery   = `INSERT INTO maker.flap_live (block_number, block_hash, address_id, live) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBidBidQuery = `INSERT INTO maker.flap_bid_bid (block_number, block_hash, address_id, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidLotQuery = `INSERT INTO maker.flap_bid_lot (block_number, block_hash, address_id, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidGuyQuery = `INSERT INTO maker.flap_bid_guy (block_number, block_hash, address_id, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidTicQuery = `INSERT INTO maker.flap_bid_tic (block_number, block_hash, address_id, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidEndQuery = `INSERT INTO maker.flap_bid_end (block_number, block_hash, address_id, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type FlapStorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository *FlapStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(blockNumber, blockHash, value.(string))
	case storage.Gem:
		return repository.insertGem(blockNumber, blockHash, value.(string))
	case storage.Beg:
		return repository.insertBeg(blockNumber, blockHash, value.(string))
	case storage.Kicks:
		return repository.insertKicks(blockNumber, blockHash, value.(string))
	case storage.Live:
		return repository.insertLive(blockNumber, blockHash, value.(string))
	case storage.BidBid:
		return repository.insertBidBid(blockNumber, blockHash, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(blockNumber, blockHash, metadata, value.(string))
	case storage.Packed:
		return repository.insertPackedValueRecord(blockNumber, blockHash, metadata, value.(map[int]string))
	default:
		panic(fmt.Sprintf("unrecognized flap contract storage name: %s", metadata.Name))
	}
}

func (repository *FlapStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlapStorageRepository) insertVat(blockNumber int, blockHash string, vat string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertVatQuery, vat)
}

func (repository *FlapStorageRepository) insertGem(blockNumber int, blockHash string, gem string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertGemQuery, gem)
}

func (repository *FlapStorageRepository) insertBeg(blockNumber int, blockHash string, beg string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertBegQuery, beg)
}

func (repository *FlapStorageRepository) insertTTL(blockNumber int, blockHash string, ttl string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertTTLQuery, ttl)
}

func (repository *FlapStorageRepository) insertTau(blockNumber int, blockHash string, tau string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertTauQuery, tau)
}

func (repository *FlapStorageRepository) insertKicks(blockNumber int, blockHash string, kicks string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, InsertKicksQuery, kicks)
}

func (repository *FlapStorageRepository) insertLive(blockNumber int, blockHash string, live string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertLiveQuery, live)
}

func (repository *FlapStorageRepository) insertBidBid(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, bid string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, insertBidBidQuery, bidID, bid)
}

func (repository *FlapStorageRepository) insertBidLot(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, lot string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, insertBidLotQuery, bidID, lot)
}

func (repository *FlapStorageRepository) insertBidGuy(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, guy string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, insertBidGuyQuery, bidID, guy)
}

func (repository *FlapStorageRepository) insertBidTic(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, tic string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, insertBidTicQuery, bidID, tic)
}

func (repository *FlapStorageRepository) insertBidEnd(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, end string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidID(blockNumber, blockHash, insertBidEndQuery, bidID, end)
}

func (repository *FlapStorageRepository) insertPackedValueRecord(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, packedValues map[int]string) error {
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
			panic(fmt.Sprintf("unrecognized flap contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

func getBidID(keys map[utils.Key]string) (string, error) {
	bidID, ok := keys[constants.BidID]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.BidID}
	}
	return bidID, nil
}

func (repository *FlapStorageRepository) insertRecordWithAddress(blockNumber int, blockHash, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flap address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, blockNumber, blockHash, addressID, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flap field with address", insertErr.Error())
		}
		return insertErr
	}

	return tx.Commit()
}

func (repository *FlapStorageRepository) insertRecordWithAddressAndBidID(blockNumber int, blockHash, query, bidID, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	addressID, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flap address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, blockNumber, blockHash, addressID, bidID, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			errorString := fmt.Sprintf("flap field with address for bid id %s", bidID)
			return shared.FormatRollbackError(errorString, insertErr.Error())
		}
		return insertErr
	}
	return tx.Commit()
}
