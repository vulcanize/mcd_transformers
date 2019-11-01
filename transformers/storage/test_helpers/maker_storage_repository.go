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

package test_helpers

import (
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockMakerStorageRepository struct {
	Cdpis                   []string
	DaiKeys                 []string
	FlapBidIds              []string
	FlipBidIds              []string
	FlopBidIds              []string
	GemKeys                 []storage.Urn
	Ilks                    []string
	Owners                  []string
	SinKeys                 []string
	Urns                    []storage.Urn
	GetCdpisCalled          bool
	GetCdpisError           error
	GetDaiKeysCalled        bool
	GetDaiKeysError         error
	GetGemKeysCalled        bool
	GetGemKeysError         error
	GetFlapBidIdsCalled     bool
	GetFlapBidIdsError      error
	GetFlipBidIdsCalledWith string
	GetFlipBidIdsError      error
	GetFlopBidIdsCalledWith string
	GetFlopBidIdsError      error
	GetIlksCalled           bool
	GetIlksError            error
	GetOwnersCalled         bool
	GetOwnersError          error
	GetVatSinKeysCalled     bool
	GetVatSinKeysError      error
	GetVowSinKeysCalled     bool
	GetVowSinKeysError      error
	GetUrnsCalled           bool
	GetUrnsError            error
}

func (repository *MockMakerStorageRepository) GetFlapBidIDs(string) ([]string, error) {
	repository.GetFlapBidIdsCalled = true
	return repository.FlapBidIds, repository.GetFlapBidIdsError
}

func (repository *MockMakerStorageRepository) GetDaiKeys() ([]string, error) {
	repository.GetDaiKeysCalled = true
	return repository.DaiKeys, repository.GetDaiKeysError
}

func (repository *MockMakerStorageRepository) GetGemKeys() ([]storage.Urn, error) {
	repository.GetGemKeysCalled = true
	return repository.GemKeys, repository.GetGemKeysError
}

func (repository *MockMakerStorageRepository) GetIlks() ([]string, error) {
	repository.GetIlksCalled = true
	return repository.Ilks, repository.GetIlksError
}

func (repository *MockMakerStorageRepository) GetVatSinKeys() ([]string, error) {
	repository.GetVatSinKeysCalled = true
	return repository.SinKeys, repository.GetVatSinKeysError
}

func (repository *MockMakerStorageRepository) GetVowSinKeys() ([]string, error) {
	repository.GetVowSinKeysCalled = true
	return repository.SinKeys, repository.GetVowSinKeysError
}

func (repository *MockMakerStorageRepository) GetUrns() ([]storage.Urn, error) {
	repository.GetUrnsCalled = true
	return repository.Urns, repository.GetUrnsError
}

func (repository *MockMakerStorageRepository) GetFlipBidIDs(contractAddress string) ([]string, error) {
	repository.GetFlipBidIdsCalledWith = contractAddress
	return repository.FlipBidIds, repository.GetFlipBidIdsError
}

func (repository *MockMakerStorageRepository) GetFlopBidIDs(contractAddress string) ([]string, error) {
	repository.GetFlopBidIdsCalledWith = contractAddress
	return repository.FlopBidIds, repository.GetFlopBidIdsError
}

func (repository *MockMakerStorageRepository) GetCDPIs() ([]string, error) {
	repository.GetCdpisCalled = true
	return repository.Cdpis, repository.GetCdpisError
}

func (repository *MockMakerStorageRepository) GetOwners() ([]string, error) {
	repository.GetOwnersCalled = true
	return repository.Owners, repository.GetOwnersError
}

func (repository *MockMakerStorageRepository) SetDB(db *postgres.DB) {}
