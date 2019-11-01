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

package storage

import (
	"errors"
	"strconv"

	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Urn struct {
	Ilk        string
	Identifier string
}

var ErrNoFlips = errors.New("no flips exist in db")

type IMakerStorageRepository interface {
	GetDaiKeys() ([]string, error)
	GetFlapBidIDs(string) ([]string, error)
	GetGemKeys() ([]Urn, error)
	GetIlks() ([]string, error)
	GetVatSinKeys() ([]string, error)
	GetVowSinKeys() ([]string, error)
	GetUrns() ([]Urn, error)
	GetCDPIs() ([]string, error)
	GetOwners() ([]string, error)
	GetFlipBidIDs(contractAddress string) ([]string, error)
	GetFlopBidIDs(contractAddress string) ([]string, error)
	SetDB(db *postgres.DB)
}

type MakerStorageRepository struct {
	db *postgres.DB
}

func (repo *MakerStorageRepository) GetFlapBidIDs(contractAddress string) ([]string, error) {
	var bidIDs []string
	addressID, addressErr := repo.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return []string{}, addressErr
	}
	err := repo.db.Select(&bidIDs, `
		SELECT bid_id FROM maker.flap_kick WHERE address_id = $1
		UNION
		SELECT kicks FROM maker.flap_kicks WHERE address_id = $1
		UNION
		SELECT bid_id from maker.tend WHERE address_id = $1
		UNION
		SELECT bid_id from maker.deal WHERE address_id = $1
		UNION
		SELECT bid_id from maker.yank WHERE address_id = $1`, addressID)
	return bidIDs, err
}

func (repo *MakerStorageRepository) GetDaiKeys() ([]string, error) {
	var daiKeys []string
	err := repo.db.Select(&daiKeys, `
		SELECT DISTINCT src FROM maker.vat_move
		UNION
		SELECT DISTINCT dst FROM maker.vat_move
		UNION
		SELECT DISTINCT w FROM maker.vat_frob
		UNION
		SELECT DISTINCT v FROM maker.vat_suck
		UNION
		SELECT DISTINCT tx_from FROM public.header_sync_transactions AS transactions
			LEFT JOIN maker.vat_heal ON vat_heal.header_id = transactions.header_id
			LEFT JOIN public.header_sync_logs ON header_sync_logs.id = vat_heal.log_id
			WHERE header_sync_logs.tx_index = transactions.tx_index
		UNION
		SELECT DISTINCT urns.identifier FROM maker.vat_fold
			INNER JOIN maker.urns on urns.id = maker.vat_fold.urn_id
	`)
	return daiKeys, err
}

func (repo *MakerStorageRepository) GetGemKeys() ([]Urn, error) {
	var gems []Urn
	err := repo.db.Select(&gems, `
		SELECT DISTINCT ilks.ilk, slip.usr AS identifier
		FROM maker.vat_slip slip
		INNER JOIN maker.ilks ilks ON ilks.id = slip.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, flux.src AS identifier
		FROM maker.vat_flux flux
		INNER JOIN maker.ilks ilks ON ilks.id = flux.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, flux.dst AS identifier
		FROM maker.vat_flux flux
		INNER JOIN maker.ilks ilks ON ilks.id = flux.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, frob.v AS identifier
		FROM maker.vat_frob frob
		INNER JOIN maker.urns on urns.id = frob.urn_id
		INNER JOIN maker.ilks ilks ON ilks.id = urns.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, grab.v AS identifier
		FROM maker.vat_grab grab
		INNER JOIN maker.urns on urns.id = grab.urn_id
		INNER JOIN maker.ilks ilks ON ilks.id = urns.ilk_id
	`)
	return gems, err
}

func (repo MakerStorageRepository) GetIlks() ([]string, error) {
	var ilks []string
	err := repo.db.Select(&ilks, `SELECT DISTINCT ilk FROM maker.ilks`)
	return ilks, err
}

func (repo *MakerStorageRepository) GetVatSinKeys() ([]string, error) {
	var sinKeys []string
	err := repo.db.Select(&sinKeys, `
		SELECT DISTINCT w FROM maker.vat_grab
		UNION
		SELECT DISTINCT u FROM maker.vat_suck
		UNION
		SELECT DISTINCT tx_from FROM public.header_sync_transactions AS transactions
			LEFT JOIN maker.vat_heal ON vat_heal.header_id = transactions.header_id
			LEFT JOIN public.header_sync_logs ON header_sync_logs.id = vat_heal.log_id
			WHERE header_sync_logs.tx_index = transactions.tx_index`)
	return sinKeys, err
}

func (repo *MakerStorageRepository) GetVowSinKeys() ([]string, error) {
	var sinKeys []string
	err := repo.db.Select(&sinKeys, `
		SELECT DISTINCT era FROM maker.vow_flog
		UNION
		SELECT DISTINCT headers.block_timestamp
		FROM maker.vow_fess
		JOIN headers ON maker.vow_fess.header_id = headers.id`)
	return sinKeys, err
}

func (repo *MakerStorageRepository) GetUrns() ([]Urn, error) {
	var urns []Urn
	err := repo.db.Select(&urns, `
		SELECT DISTINCT ilks.ilk, urns.identifier
		FROM maker.urns
		JOIN maker.ilks on maker.ilks.id = maker.urns.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, fork.src AS identifier
		FROM maker.vat_fork fork
		INNER JOIN maker.ilks ilks ON ilks.id = fork.ilk_id
		UNION
		SELECT DISTINCT ilks.ilk, fork.dst AS identifier
		FROM maker.vat_fork fork
		INNER JOIN maker.ilks ilks ON ilks.id = fork.ilk_id`)
	return urns, err
}

func (repo *MakerStorageRepository) GetCDPIs() ([]string, error) {
	nullValue := 0
	var maxCDPI int
	readErr := repo.db.Get(&maxCDPI, `
		SELECT COALESCE(MAX(cdpi), $1)
		FROM maker.cdp_manager_cdpi`, nullValue)
	if readErr != nil {
		return nil, readErr
	}
	if maxCDPI == nullValue {
		return []string{}, nil
	}
	return rangeIntsAsStrings(1, maxCDPI), readErr
}

func (repo *MakerStorageRepository) GetOwners() ([]string, error) {
	var owners []string
	err := repo.db.Select(&owners, `
		SELECT DISTINCT owner
		FROM maker.cdp_manager_owns`)
	return owners, err
}

func (repo *MakerStorageRepository) GetFlipBidIDs(contractAddress string) ([]string, error) {
	var bidIds []string
	addressID, addressErr := repo.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return []string{}, addressErr
	}
	err := repo.db.Select(&bidIds, `
   		SELECT DISTINCT bid_id FROM maker.tick
		WHERE address_id = $1
		UNION
   		SELECT DISTINCT bid_id FROM maker.flip_kick
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.tend
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.dent
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.deal
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.yank
		WHERE address_id = $1
		UNION
		SELECT DISTINCT kicks FROM maker.flip_kicks
		WHERE address_id = $1`, addressID)
	return bidIds, err
}

func (repo *MakerStorageRepository) GetFlopBidIDs(contractAddress string) ([]string, error) {
	var bidIDs []string
	addressID, addressErr := repo.GetOrCreateAddress(contractAddress)
	if addressErr != nil {
		return []string{}, addressErr
	}
	err := repo.db.Select(&bidIDs, `
		SELECT bid_id FROM maker.flop_kick
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.dent
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.deal
		WHERE address_id = $1
		UNION
		SELECT DISTINCT bid_id FROM maker.yank
		WHERE address_id = $1
		UNION
		SELECT DISTINCT kicks FROM maker.flop_kicks
		WHERE address_id = $1`, addressID)
	return bidIDs, err
}

func (repo *MakerStorageRepository) GetOrCreateAddress(contractAddress string) (int64, error) {
	return repository.GetOrCreateAddress(repo.db, contractAddress)
}

func (repo *MakerStorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func rangeIntsAsStrings(start, end int) []string {
	var strSlice []string
	for i := start; i <= end; i++ {
		strSlice = append(strSlice, strconv.Itoa(i))
	}
	return strSlice
}
