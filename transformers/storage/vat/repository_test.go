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
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"database/sql"
)

var _ = Describe("Vat storage repository", func() {
	var (
		db              *postgres.DB
		repo            vat.VatStorageRepository
		fakeBlockNumber = 123
		fakeBlockHash   = "expected_block_hash"
		fakeIlk         = "fake_ilk"
		fakeGuy         = "fake_urn"
		fakeUint256     = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = vat.VatStorageRepository{}
		repo.SetDB(db)
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

		It("returns error if metadata missing guy", func() {
			malformedDaiMetadata := utils.GetStorageValueMetadata(vat.Dai, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedDaiMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("gem", func() {
		It("writes row", func() {
			gemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: fakeIlk, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, gemMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key_one, guy AS key_two, gem AS value FROM maker.vat_gem`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeGuy, fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedGemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedGemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("ilk Art", func() {
		It("writes row", func() {
			ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, art AS value FROM maker.vat_ilk_art`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})
	})

	Describe("ilk dust", func() {
		It("writes row", func() {
			ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDustMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, dust AS value FROM maker.vat_ilk_dust`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkDustMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})
	})

	Describe("ilk line", func() {
		It("writes row", func() {
			ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, line AS value FROM maker.vat_ilk_line`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkLineMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})
	})

	Describe("ilk rate", func() {
		It("writes row", func() {
			ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, rate AS value FROM maker.vat_ilk_rate`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkRateMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})
	})

	Describe("ilk spot", func() {
		It("writes row", func() {
			ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, spot AS value FROM maker.vat_ilk_spot`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkSpotMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
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

		It("returns error if metadata missing guy", func() {
			malformedSinMetadata := utils.GetStorageValueMetadata(vat.Sin, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedSinMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("urn art", func() {
		It("writes row", func() {
			urnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: fakeIlk, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, urnArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `
				SELECT block_number, block_hash, ilks.id AS key_one, urns.guy AS key_two, art AS value
				FROM maker.vat_urn_art
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_art.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id
			`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeGuy, fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("urn ink", func() {
		It("writes row", func() {
			urnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: fakeIlk, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, urnInkMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `
				SELECT block_number, block_hash, ilks.id AS key_one, urns.guy AS key_two, ink AS value
				FROM maker.vat_urn_ink
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_ink.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id
			`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeGuy, fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("vat debt", func() {
		It("persists conflicting header", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			header := fakes.GetFakeHeader(int64(fakeBlockNumber))
			headerId, err := headerRepository.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(fakeBlockNumber, fakeBlockHash, vat.DebtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, debt AS key, conflicting_header_id AS value FROM maker.vat_debt`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, strconv.Itoa(int(headerId)))
		})

		It("persists corresponding header", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			header := fakes.GetFakeHeader(int64(fakeBlockNumber))
			headerId, err := headerRepository.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(fakeBlockNumber, header.Hash, vat.DebtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, debt AS key, confirmed_header_id AS value FROM maker.vat_debt`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, header.Hash, fakeUint256, strconv.Itoa(int(headerId)))
		})

		It("persist when there's neither a corresponding nor conflicting header", func() {
			err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.DebtMetadata, fakeUint256)
			Expect(err).NotTo(HaveOccurred())

			var debtResult VariableRes
			err = db.Get(&debtResult, `SELECT block_number, block_hash, debt AS value FROM maker.vat_debt`)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(debtResult, fakeBlockNumber, fakeBlockHash, fakeUint256)

			type Res struct {
				BlockMetadata
				ConfirmedHeaderId sql.NullInt64 `db:"confirmed_header_id"`
				ConflictingHeaderId sql.NullInt64 `db:"conflicting_header_id"`
			}
			var result Res
			err = db.Get(&result, `SELECT block_number, block_hash, confirmed_header_id, conflicting_header_id FROM maker.vat_debt`)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.ConfirmedHeaderId).To(Equal(sql.NullInt64{}))
			Expect(result.ConflictingHeaderId).To(Equal(sql.NullInt64{}))
		})

		It("handles a reorg for a confirmed header", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			header := fakes.GetFakeHeader(int64(fakeBlockNumber))
			headerId, err := headerRepository.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(fakeBlockNumber, header.Hash, vat.DebtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, debt AS key, confirmed_header_id AS value FROM maker.vat_debt`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, header.Hash, fakeUint256, strconv.Itoa(int(headerId)))

			_, err = db.Exec(`DELETE from headers where id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			type Res struct {
				BlockMetadata
				ConfirmedHeaderId sql.NullInt64 `db:"confirmed_header_id"`
				ConflictingHeaderId sql.NullInt64 `db:"conflicting_header_id"`
			}

			var res Res
			err = db.Get(&res, `SELECT block_number, block_hash, confirmed_header_id, conflicting_header_id FROM maker.vat_debt`)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.ConfirmedHeaderId).To(Equal(sql.NullInt64{}))
			Expect(res.ConflictingHeaderId).To(Equal(sql.NullInt64{}))
		})

		It("handles a reorg for a conflicting header", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			header := fakes.GetFakeHeader(int64(fakeBlockNumber))
			headerId, err := headerRepository.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(fakeBlockNumber, fakeBlockHash, vat.DebtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, debt AS key, conflicting_header_id AS value FROM maker.vat_debt`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, strconv.Itoa(int(headerId)))

			_, err = db.Exec(`DELETE from headers where id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			type Res struct {
				BlockMetadata
				ConfirmedHeaderId sql.NullInt64 `db:"confirmed_header_id"`
				ConflictingHeaderId sql.NullInt64 `db:"conflicting_header_id"`
			}

			var res Res
			err = db.Get(&res, `SELECT block_number, block_hash, confirmed_header_id, conflicting_header_id FROM maker.vat_debt`)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.ConfirmedHeaderId).To(Equal(sql.NullInt64{}))
			Expect(res.ConflictingHeaderId).To(Equal(sql.NullInt64{}))
		})
	})


	FDescribe("vat vice", func() {
		It("persists vat vice", func() {
			err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.ViceMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result VariableRes
			err = db.Get(&result, `SELECT block_number, block_hash, vice AS value FROM maker.vat_vice`)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)

			type JoinResult struct {
				ConfirmedHeaderId sql.NullInt64 `db:"confirmed_header_id"`
				ConflictingHeaderId sql.NullInt64 `db:"conflicting_header_id"`
				VatViceId sql.NullInt64 `db:"vat_vice_id"`
			}

			var joinResult JoinResult
			err = db.Get(&joinResult, `SELECT confirmed_header_id, conflicting_header_id, vat_vice_id FROM maker.vat_vice_header`)
			Expect(err).NotTo(HaveOccurred())
			Expect(joinResult.ConfirmedHeaderId).To(Equal(sql.NullInt64{}))
			Expect(joinResult.ConflictingHeaderId).To(Equal(sql.NullInt64{}))
		})

		It("persists conflicting header", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			header := fakes.GetFakeHeader(int64(fakeBlockNumber))
			headerId, err := headerRepository.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(fakeBlockNumber, fakeBlockHash, vat.ViceMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result VariableRes
			err = db.Get(&result, `SELECT block_number, block_hash, vice AS value FROM maker.vat_vice`)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)

			type JoinResult struct {
				ConfirmedHeaderId sql.NullInt64 `db:"confirmed_header_id"`
				ConflictingHeaderId sql.NullInt64 `db:"conflicting_header_id"`
				VatViceId sql.NullInt64 `db:"vat_vice_id"`
			}

			var joinResult JoinResult
			err = db.Get(&joinResult, `SELECT confirmed_header_id, conflicting_header_id, vat_vice_id FROM maker.vat_vice_header`)
			Expect(err).NotTo(HaveOccurred())
			Expect(joinResult.ConflictingHeaderId).To(Equal(sql.NullInt64{Valid: true, Int64: headerId}))
			Expect(joinResult.ConfirmedHeaderId).To(Equal(sql.NullInt64{}))

		})

		FIt("persists confirmed header", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			header := fakes.GetFakeHeader(int64(fakeBlockNumber))
			headerId, err := headerRepository.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(fakeBlockNumber, header.Hash, vat.ViceMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result VariableRes
			err = db.Get(&result, `SELECT block_number, block_hash, vice AS value FROM maker.vat_vice`)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(result, fakeBlockNumber, header.Hash, fakeUint256)

			type JoinResult struct {
				ConfirmedHeaderId sql.NullInt64 `db:"confirmed_header_id"`
				ConflictingHeaderId sql.NullInt64 `db:"conflicting_header_id"`
				VatViceId sql.NullInt64 `db:"vat_vice_id"`
			}

			var joinResult JoinResult
			err = db.Get(&joinResult, `SELECT confirmed_header_id, conflicting_header_id, vat_vice_id FROM maker.vat_vice_header`)
			Expect(err).NotTo(HaveOccurred())
			Expect(joinResult.ConfirmedHeaderId).To(Equal(sql.NullInt64{Valid: true, Int64: headerId}))
			Expect(joinResult.ConflictingHeaderId).To(Equal(sql.NullInt64{}))
		})


		It("handles a reorg for a confirmed header", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			header := fakes.GetFakeHeader(int64(fakeBlockNumber))
			headerId, err := headerRepository.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(fakeBlockNumber, header.Hash, vat.ViceMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result VariableRes
			err = db.Get(&result, `SELECT block_number, block_hash, vice AS value FROM maker.vat_vice`)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(result, fakeBlockNumber, header.Hash, fakeUint256)

			type JoinResult struct {
				ConfirmedHeaderId sql.NullInt64 `db:"confirmed_header_id"`
				ConflictingHeaderId sql.NullInt64 `db:"conflicting_header_id"`
				VatViceId sql.NullInt64 `db:"vat_vice_id"`
			}

			var joinResult JoinResult
			err = db.Get(&joinResult, `SELECT confirmed_header_id, conflicting_header_id, vat_vice_id FROM maker.vat_vice_header`)
			Expect(err).NotTo(HaveOccurred())
			Expect(joinResult.ConfirmedHeaderId).To(Equal(sql.NullInt64{Valid: true, Int64: headerId}))
			Expect(joinResult.ConflictingHeaderId).To(Equal(sql.NullInt64{}))

			_, err = db.Exec(`DELETE from headers where id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())

			type Res struct {
				BlockMetadata
				ConfirmedHeaderId sql.NullInt64 `db:"confirmed_header_id"`
				ConflictingHeaderId sql.NullInt64 `db:"conflicting_header_id"`
			}

			var res Res
			err = db.Get(&res, `SELECT confirmed_header_id, conflicting_header_id FROM maker.vat_vice_header`)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.ConfirmedHeaderId).To(Equal(sql.NullInt64{}))
			Expect(res.ConflictingHeaderId).To(Equal(sql.NullInt64{}))
		})

	})

	It("persists vat Line", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LineMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, line AS value FROM maker.vat_line`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("persists vat live", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LiveMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, live AS value FROM maker.vat_live`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})
})
