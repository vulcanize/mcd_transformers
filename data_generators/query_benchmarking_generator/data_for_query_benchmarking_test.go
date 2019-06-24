package query_benchmarking_generator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"math/rand"
)

var _ = Describe("data generator for query benchmarking", func() {
	var (
		db    *postgres.DB
		state BenchmarkingDataGeneratorState
		seed  int64
	)

	BeforeEach(func() {
		seed = rand.Int63()
		rand.Seed(int64(seed))
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		state = NewBenchmarkingDataGeneratorState(db)
	})

	Describe("Ilks", func() {
		It("generates the specified number of ilks with initial storage records", func() {
			ilkErr := state.GenerateDataForQueryTesting(10, 0, 0)
			Expect(ilkErr).NotTo(HaveOccurred())

			Expect(len(state.Ilks)).To(Equal(10))

			var ilkCount int
			ilkCountErr := db.Get(&ilkCount, `SELECT count(*) FROM maker.ilks`)
			Expect(ilkCountErr).NotTo(HaveOccurred())
			Expect(ilkCount).To(Equal(10))

			//also generates the ilk initial storage records
			assertIlkStorageRecordCount(db, 10)
		})

		It("generates the additional storage records associated with ilks for each additional block", func() {
			ilkErr := state.GenerateDataForQueryTesting(1, 0, 10)
			Expect(ilkErr).NotTo(HaveOccurred())
			Expect(len(state.Ilks)).To(Equal(1))

			assertIlkStorageRecordCount(db, 11)
		})
	})

	Describe("Urns", func() {
		It("generates the n urns per ilk per block with initial storage records", func() {
			ilkErr := state.GenerateDataForQueryTesting(1, 10, 0)
			Expect(ilkErr).NotTo(HaveOccurred())
			Expect(len(state.Ilks)).To(Equal(1))
			assertIlkStorageRecordCount(db, 1)

			Expect(len(state.Urns)).To(Equal(10))
			var urnCount int
			getErr := db.Get(&urnCount, `SELECT count(*) FROM maker.urns`)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(urnCount).To(Equal(10))
			assertUrnStorageRecordCount(db, 10) //also generates the initial storage records for each block for each urn
		})

		It("generates additional storage records associated with urns for each additional block", func() {
			ilkErr := state.GenerateDataForQueryTesting(1, 1, 10)
			Expect(ilkErr).NotTo(HaveOccurred())
			Expect(len(state.Ilks)).To(Equal(1))

			Expect(len(state.Urns)).To(Equal(1))
			assertUrnStorageRecordCount(db, 11)
		})
	})
})

func assertIlkStorageRecordCount(db *postgres.DB, count int) {
	var recordCount int

	vatIlkRateErr := db.Get(&recordCount, `SELECT count(*) FROM maker.vat_ilk_rate`)
	Expect(vatIlkRateErr).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	vatIlkArtErr := db.Get(&recordCount, `SELECT count(*) FROM maker.vat_ilk_art`)
	Expect(vatIlkArtErr).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	vatIlkSpot := db.Get(&recordCount, `SELECT count(*) FROM maker.vat_ilk_spot`)
	Expect(vatIlkSpot).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	vatIlkLine := db.Get(&recordCount, `SELECT count(*) FROM maker.vat_ilk_line`)
	Expect(vatIlkLine).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	vatIlkDust := db.Get(&recordCount, `SELECT count(*) FROM maker.vat_ilk_dust`)
	Expect(vatIlkDust).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	catIlkLump := db.Get(&recordCount, `SELECT count(*) FROM maker.cat_ilk_lump`)
	Expect(catIlkLump).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	catIlkFlip := db.Get(&recordCount, `SELECT count(*) FROM maker.cat_ilk_flip`)
	Expect(catIlkFlip).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	jugIlkRho := db.Get(&recordCount, `SELECT count(*) FROM maker.jug_ilk_rho`)
	Expect(jugIlkRho).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	jugIlkDuty := db.Get(&recordCount, `SELECT count(*) FROM maker.jug_ilk_duty`)
	Expect(jugIlkDuty).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	spotIlkPip := db.Get(&recordCount, `SELECT count(*) FROM maker.spot_ilk_pip`)
	Expect(spotIlkPip).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	spotIlkMat := db.Get(&recordCount, `SELECT count(*) FROM maker.spot_ilk_mat`)
	Expect(spotIlkMat).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))
}

func assertUrnStorageRecordCount(db *postgres.DB, count int) {
	var recordCount int

	vatUrnArtErr := db.Get(&recordCount, `SELECT count(*) FROM maker.vat_urn_art`)
	Expect(vatUrnArtErr).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))

	vatUrnInkErr := db.Get(&recordCount, `SELECT count(*) FROM maker.vat_urn_ink`)
	Expect(vatUrnInkErr).NotTo(HaveOccurred())
	Expect(recordCount).To(Equal(count))
}