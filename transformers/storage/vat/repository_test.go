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

package vat_test

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"math/rand"
	"strconv"
)

var _ = Describe("Vat storage repository", func() {
	var (
		db              *postgres.DB
		repo            vat.VatStorageRepository
		fakeGuy         string
		fakeBlockNumber = rand.Int()
		fakeBlockHash   = "expected_block_hash"
		fakeUint256     = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = vat.VatStorageRepository{}
		repo.SetDB(db)
		fakeGuy = test_data.RandomString(10)
	})

	Describe("dai", func() {
		It("writes a row", func() {
			daiMetadata := utils.GetStorageValueMetadata(vat.Dai, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, daiMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, guy AS key, dai AS value FROM maker.vat_dai`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			daiMetadata := utils.GetStorageValueMetadata(vat.Dai, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, daiMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, daiMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_dai`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing guy", func() {
			malformedDaiMetadata := utils.GetStorageValueMetadata(vat.Dai, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedDaiMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("gem", func() {
		It("writes row", func() {
			gemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, gemMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key_one, guy AS key_two, gem AS value FROM maker.vat_gem`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			gemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, gemMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, gemMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_gem`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedGemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedGemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("ilk Art", func() {
		It("writes row", func() {
			ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, art AS value FROM maker.vat_ilk_art`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_art`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Art",
			PropertyValue: strconv.Itoa(rand.Int()),
		})
	})

	Describe("ilk dust", func() {
		It("writes row", func() {
			ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDustMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, dust AS value FROM maker.vat_ilk_dust`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDustMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDustMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_dust`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkDustMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Dust",
			PropertyValue: strconv.Itoa(rand.Int()),
		})
	})

	Describe("ilk line", func() {
		It("writes row", func() {
			ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, line AS value FROM maker.vat_ilk_line`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_line`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkLineMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Line",
			PropertyValue: strconv.Itoa(rand.Int()),
		})
	})

	Describe("ilk rate", func() {
		It("writes row", func() {
			ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, rate AS value FROM maker.vat_ilk_rate`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_rate`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkRateMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Rate",
			PropertyValue: strconv.Itoa(rand.Int()),
		})
	})

	Describe("ilk spot", func() {
		It("writes row", func() {
			ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, spot AS value FROM maker.vat_ilk_spot`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_spot`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkSpotMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Spot",
			PropertyValue: strconv.Itoa(rand.Int()),
		})
	})

	Describe("sin", func() {
		It("writes a row", func() {
			sinMetadata := utils.GetStorageValueMetadata(vat.Sin, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, sinMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, guy AS key, sin AS value FROM maker.vat_sin`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			sinMetadata := utils.GetStorageValueMetadata(vat.Sin, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, sinMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, sinMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_sin`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing guy", func() {
			malformedSinMetadata := utils.GetStorageValueMetadata(vat.Sin, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedSinMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("urn art", func() {
		It("writes row", func() {
			urnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, urnArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `
				SELECT block_number, block_hash, ilks.id AS key_one, urns.identifier AS key_two, art AS value
				FROM maker.vat_urn_art
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_art.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id
			`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			urnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, urnArtMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, urnArtMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_urn_art`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})

		Describe("updating current_urn_state trigger table", func() {
			var (
				urnId, rawTimestamp int64
			)

			BeforeEach(func() {
				rawTimestamp = int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				var urnErr error
				urnId, urnErr = shared.GetOrCreateUrn(fakeGuy, test_helpers.FakeIlk.Hex, db)
				Expect(urnErr).NotTo(HaveOccurred())
			})

			It("inserts art, ilk values, and timestamp if urn is not yet in table", func() {
				_, urnArtErr := db.Exec(`INSERT INTO maker.vat_urn_art (block_number, urn_id, art) VALUES ($1, $2, $3)`,
					fakeBlockNumber, urnId, fakeUint256)
				Expect(urnArtErr).NotTo(HaveOccurred())

				expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}
				var urnState currentUrnState
				queryErr := db.Get(&urnState, `SELECT urn_identifier, ilk_identifier, art, created, updated FROM api.current_urn_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(urnState.UrnIdentifier).To(Equal(fakeGuy))
				Expect(urnState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(urnState.Art).To(Equal(fakeUint256))
				Expect(urnState.Created).To(Equal(expectedTime))
				Expect(urnState.Updated).To(Equal(expectedTime))
			})

			It("sets art and time updated if new diff is from later block", func() {
				// set up existing row from earlier block
				earlierTimestamp := FormatTimestamp(rawTimestamp - 1)
				_, urnSetupErr := db.Exec(
					`INSERT INTO api.current_urn_state (urn_identifier, ilk_identifier, art, created, updated) VALUES ($1, $2, $3, $4::TIMESTAMP, $4::TIMESTAMP)`,
					fakeGuy, test_helpers.FakeIlk.Identifier, rand.Int(), earlierTimestamp)
				Expect(urnSetupErr).NotTo(HaveOccurred())

				// trigger update to row from later block
				_, urnArtErr := db.Exec(`INSERT INTO maker.vat_urn_art (block_number, urn_id, art) VALUES ($1, $2, $3)`,
					fakeBlockNumber, urnId, fakeUint256)
				Expect(urnArtErr).NotTo(HaveOccurred())

				expectedTimeCreated := sql.NullString{String: earlierTimestamp, Valid: true}
				expectedTimeUpdated := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}
				var urnState currentUrnState
				queryErr := db.Get(&urnState, `SELECT urn_identifier, ilk_identifier, art, created, updated FROM api.current_urn_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(urnState.UrnIdentifier).To(Equal(fakeGuy))
				Expect(urnState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(urnState.Art).To(Equal(fakeUint256))
				Expect(urnState.Created).To(Equal(expectedTimeCreated))
				Expect(urnState.Updated).To(Equal(expectedTimeUpdated))
			})
		})
	})

	Describe("urn ink", func() {
		It("writes row", func() {
			urnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, urnInkMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `
				SELECT block_number, block_hash, ilks.id AS key_one, urns.identifier AS key_two, ink AS value
				FROM maker.vat_urn_ink
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_ink.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id
			`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			urnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, urnInkMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, urnInkMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_urn_ink`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})

		Describe("updating current_urn_state trigger table", func() {
			var (
				urnId, rawTimestamp int64
			)

			BeforeEach(func() {
				rawTimestamp = int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				var urnErr error
				urnId, urnErr = shared.GetOrCreateUrn(fakeGuy, test_helpers.FakeIlk.Hex, db)
				Expect(urnErr).NotTo(HaveOccurred())
			})

			It("inserts ink and timestamp if urn_id is not yet in table", func() {
				_, urnInkErr := db.Exec(`INSERT INTO maker.vat_urn_ink (block_number, urn_id, ink) VALUES ($1, $2, $3)`,
					fakeBlockNumber, urnId, fakeUint256)
				Expect(urnInkErr).NotTo(HaveOccurred())

				expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}
				var urnState currentUrnState
				queryErr := db.Get(&urnState, `SELECT urn_identifier, ilk_identifier, ink, created, updated FROM api.current_urn_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(urnState.UrnIdentifier).To(Equal(fakeGuy))
				Expect(urnState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(urnState.Ink).To(Equal(fakeUint256))
				Expect(urnState.Created).To(Equal(expectedTime))
				Expect(urnState.Updated).To(Equal(expectedTime))
			})

			It("sets ink and time updated if new diff is from later block", func() {
				// set up existing row from earlier block
				earlierTimestamp := FormatTimestamp(rawTimestamp - 1)
				_, urnSetupErr := db.Exec(
					`INSERT INTO api.current_urn_state (urn_identifier, ilk_identifier, ink, created, updated) VALUES ($1, $2, $3, $4::TIMESTAMP, $4::TIMESTAMP)`,
					fakeGuy, test_helpers.FakeIlk.Identifier, rand.Int(), earlierTimestamp)
				Expect(urnSetupErr).NotTo(HaveOccurred())

				// trigger update to row from later block
				_, urnInkErr := db.Exec(`INSERT INTO maker.vat_urn_ink (block_number, urn_id, ink) VALUES ($1, $2, $3)`,
					fakeBlockNumber, urnId, fakeUint256)
				Expect(urnInkErr).NotTo(HaveOccurred())

				expectedTimeCreated := sql.NullString{String: earlierTimestamp, Valid: true}
				expectedTimeUpdated := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}
				var urnState currentUrnState
				queryErr := db.Get(&urnState, `SELECT urn_identifier, ilk_identifier, ink, created, updated FROM api.current_urn_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(urnState.UrnIdentifier).To(Equal(fakeGuy))
				Expect(urnState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(urnState.Ink).To(Equal(fakeUint256))
				Expect(urnState.Created).To(Equal(expectedTimeCreated))
				Expect(urnState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("sets time created if new diff is from earlier block", func() {
				// set up existing row from later block
				laterTimestamp := FormatTimestamp(rawTimestamp + 1)
				_, urnSetupErr := db.Exec(`INSERT INTO api.current_urn_state (urn_identifier, ilk_identifier, ink, created, updated) VALUES ($1, $2, $3, $4::TIMESTAMP, $4::TIMESTAMP)`,
					fakeGuy, test_helpers.FakeIlk.Identifier, fakeUint256, laterTimestamp)
				Expect(urnSetupErr).NotTo(HaveOccurred())

				// trigger update to row from earlier block
				_, urnInkErr := db.Exec(`INSERT INTO maker.vat_urn_ink (block_number, urn_id, ink) VALUES ($1, $2, $3)`,
					fakeBlockNumber, urnId, rand.Int())
				Expect(urnInkErr).NotTo(HaveOccurred())

				expectedTimeCreated := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}
				expectedTimeUpdated := sql.NullString{String: laterTimestamp, Valid: true}
				var urnState currentUrnState
				queryErr := db.Get(&urnState, `SELECT urn_identifier, ilk_identifier, ink, created, updated FROM api.current_urn_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(urnState.UrnIdentifier).To(Equal(fakeGuy))
				Expect(urnState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(urnState.Ink).To(Equal(fakeUint256))
				Expect(urnState.Created).To(Equal(expectedTimeCreated))
				Expect(urnState.Updated).To(Equal(expectedTimeUpdated))
			})
		})
	})

	It("persists vat debt", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.DebtMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, debt AS value FROM maker.vat_debt`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vat debt", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.DebtMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.DebtMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_debt`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat vice", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.ViceMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vice AS value FROM maker.vat_vice`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vat vice", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.ViceMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.ViceMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_vice`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat Line", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LineMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, line AS value FROM maker.vat_line`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vat Line", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LineMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LineMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_line`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat live", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LiveMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, live AS value FROM maker.vat_live`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vat live", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LiveMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LiveMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_live`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})

type currentUrnState struct {
	UrnIdentifier string `db:"urn_identifier"`
	IlkIdentifier string `db:"ilk_identifier"`
	Ink           string
	Art           string
	Created       sql.NullString
	Updated       sql.NullString
}
