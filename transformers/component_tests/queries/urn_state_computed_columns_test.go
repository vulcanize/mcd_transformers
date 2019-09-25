// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package queries

import (
	"math/big"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Urn state computed columns", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeGuy          = "fakeAddress"
		fakeHeader       core.Header
		headerId, logId  int64
		vatRepository    vat.VatStorageRepository
		catRepository    cat.CatStorageRepository
		jugRepository    jug.JugStorageRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		fakeHeaderSyncLog := test_data.CreateTestLog(headerId, db)
		logId = fakeHeaderSyncLog.ID

		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("urn_state_ilk", func() {
		It("returns the ilk for an urn", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			fakeGuy := "fakeAddress"
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			ilkRate, convertRateErr := strconv.Atoi(ilkValues[vat.IlkRate])
			Expect(convertRateErr).NotTo(HaveOccurred())
			urnSetupData.Rate = ilkRate
			ilkSpot, convertSpotErr := strconv.Atoi(ilkValues[vat.IlkSpot])
			Expect(convertSpotErr).NotTo(HaveOccurred())
			urnSetupData.Spot = ilkSpot
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			getIlkErr := db.Get(&result,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
					FROM api.urn_state_ilk(
					(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
					FROM api.get_urn($1, $2, $3)))`, test_helpers.FakeIlk.Identifier, fakeGuy, fakeHeader.BlockNumber)

			Expect(getIlkErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("urn_state_frobs", func() {
		It("returns frobs for an urn_state", func() {
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			frobRepo := vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			frobEvent := test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobEvent.ForeignKeyValues[constants.UrnFK] = fakeGuy
			frobEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			frobEvent.ColumnValues[constants.HeaderFK] = headerId
			frobEvent.ColumnValues[constants.LogFK] = logId
			insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobEvent})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs test_helpers.FrobEvent
			getFrobsErr := db.Get(&actualFrobs,
				`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_state_frobs(
                        (SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
                         FROM api.all_urns($1))
                    )`, fakeBlock)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			expectedFrobs := test_helpers.FrobEvent{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: fakeGuy,
				Dink:          frobEvent.ColumnValues["dink"].(string),
				Dart:          frobEvent.ColumnValues["dart"].(string),
			}

			Expect(actualFrobs).To(Equal(expectedFrobs))
		})

		Describe("result pagination", func() {
			var frobEventOne, frobEventTwo shared.InsertionModel

			BeforeEach(func() {
				urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
				urnSetupData.Header.Hash = fakeHeader.Hash
				urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
				test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

				frobRepo := vat_frob.VatFrobRepository{}
				frobRepo.SetDB(db)

				frobEventOne = test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
				frobEventOne.ForeignKeyValues[constants.UrnFK] = fakeGuy
				frobEventOne.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				frobEventOne.ColumnValues[constants.HeaderFK] = headerId
				frobEventOne.ColumnValues[constants.LogFK] = logId
				insertFrobErrOne := frobRepo.Create([]shared.InsertionModel{frobEventOne})
				Expect(insertFrobErrOne).NotTo(HaveOccurred())

				// insert more recent frob for same urn
				laterBlock := fakeBlock + 1
				fakeHeaderTwo := fakes.GetFakeHeader(int64(laterBlock))
				headerTwoId, insertHeaderTwoErr := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
				Expect(insertHeaderTwoErr).NotTo(HaveOccurred())
				logTwoId := test_data.CreateTestLog(headerTwoId, db).ID

				frobEventTwo = test_data.CopyModel(test_data.VatFrobModelWithNegativeDink)
				frobEventTwo.ForeignKeyValues[constants.UrnFK] = fakeGuy
				frobEventTwo.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				frobEventTwo.ColumnValues[constants.HeaderFK] = headerTwoId
				frobEventTwo.ColumnValues[constants.LogFK] = logTwoId
				insertFrobErrTwo := frobRepo.Create([]shared.InsertionModel{frobEventTwo})
				Expect(insertFrobErrTwo).NotTo(HaveOccurred())
			})

			It("limits results to latest block number if max_results argument is provided", func() {
				maxResults := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs,
					`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_state_frobs(
						(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
						 FROM api.get_urn($1, $2)), $3)`, test_helpers.FakeIlk.Identifier, fakeGuy, maxResults)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				expectedFrob := test_helpers.FrobEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Dink:          frobEventTwo.ColumnValues["dink"].(string),
					Dart:          frobEventTwo.ColumnValues["dart"].(string),
				}
				Expect(actualFrobs).To(ConsistOf(expectedFrob))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs,
					`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_state_frobs(
						(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
						 FROM api.get_urn($1, $2)), $3, $4)`,
					test_helpers.FakeIlk.Identifier, fakeGuy, maxResults, resultOffset)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				expectedFrobs := test_helpers.FrobEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Dink:          frobEventOne.ColumnValues["dink"].(string),
					Dart:          frobEventOne.ColumnValues["dart"].(string),
				}
				Expect(actualFrobs).To(ConsistOf(expectedFrobs))
			})
		})
	})

	Describe("urn_state_bites", func() {
		It("returns bites for an urn_state", func() {
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			biteRepo := bite.BiteRepository{}
			biteRepo.SetDB(db)
			biteEvent := randomizeBite(test_data.BiteModel)
			biteEvent.Urn = fakeGuy
			biteEvent.Ilk = test_helpers.FakeIlk.Hex
			biteEvent.HeaderID = headerId
			biteEvent.LogID = logId
			insertBiteErr := biteRepo.Create([]interface{}{biteEvent})
			Expect(insertBiteErr).NotTo(HaveOccurred())

			var actualBites test_helpers.BiteEvent
			getBitesErr := db.Get(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_state_bites(
				    (SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
				    FROM api.all_urns($1)))`,
				fakeBlock)
			Expect(getBitesErr).NotTo(HaveOccurred())

			expectedBites := test_helpers.BiteEvent{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: fakeGuy,
				Ink:           biteEvent.Ink,
				Art:           biteEvent.Art,
				Tab:           biteEvent.Tab,
			}

			Expect(actualBites).To(Equal(expectedBites))
		})

		Describe("result pagination", func() {
			var biteEventOne, biteEventTwo bite.BiteModel

			BeforeEach(func() {
				urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
				urnSetupData.Header.Hash = fakeHeader.Hash
				urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
				test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

				biteRepo := bite.BiteRepository{}
				biteRepo.SetDB(db)

				biteEventOne = randomizeBite(test_data.BiteModel)
				biteEventOne.Urn = fakeGuy
				biteEventOne.Ilk = test_helpers.FakeIlk.Hex
				biteEventOne.HeaderID = headerId
				biteEventOne.LogID = logId
				insertBiteOneErr := biteRepo.Create([]interface{}{biteEventOne})
				Expect(insertBiteOneErr).NotTo(HaveOccurred())

				// insert more recent bite for same urn
				laterBlock := fakeBlock + 1
				fakeHeaderTwo := fakes.GetFakeHeader(int64(laterBlock))
				headerTwoId, insertHeaderTwoErr := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
				Expect(insertHeaderTwoErr).NotTo(HaveOccurred())
				logTwoId := test_data.CreateTestLog(headerTwoId, db).ID

				biteEventTwo = randomizeBite(test_data.BiteModel)
				biteEventTwo.Urn = fakeGuy
				biteEventTwo.Ilk = test_helpers.FakeIlk.Hex
				biteEventTwo.HeaderID = headerTwoId
				biteEventTwo.LogID = logTwoId
				insertBiteTwoErr := biteRepo.Create([]interface{}{biteEventTwo})
				Expect(insertBiteTwoErr).NotTo(HaveOccurred())
			})

			It("limits results to latest block number if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				getBitesErr := db.Select(&actualBites, `
					SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_state_bites(
						(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
						 FROM api.get_urn($1, $2)), $3)`,
					test_helpers.FakeIlk.Identifier, fakeGuy, maxResults)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Ink:           biteEventTwo.Ink,
					Art:           biteEventTwo.Art,
					Tab:           biteEventTwo.Tab,
				}
				Expect(actualBites).To(ConsistOf(expectedBite))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBites []test_helpers.BiteEvent
				getBitesErr := db.Select(&actualBites, `
					SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_state_bites(
						(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, ratio, safe, created, updated)::api.urn_state
						 FROM api.get_urn($1, $2)), $3, $4)`,
					test_helpers.FakeIlk.Identifier, fakeGuy, maxResults, resultOffset)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Ink:           biteEventOne.Ink,
					Art:           biteEventOne.Art,
					Tab:           biteEventOne.Tab,
				}
				Expect(actualBites).To(ConsistOf(expectedBite))
			})
		})
	})
})

func randomizeBite(bite bite.BiteModel) bite.BiteModel {
	bite.Ink = big.NewInt(rand.Int63()).String()
	bite.Art = big.NewInt(rand.Int63()).String()
	bite.Tab = big.NewInt(rand.Int63()).String()
	return bite
}
