package query_benchmarking_generator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"math/rand"
)


var _ = Describe("data generator for query benchmarking", func() {
	Describe("with a given seed", func() {
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

		It("generates the specified number of ilks", func() {
			ilkErr := state.GenerateDataForIlkQueryTesting(10, 0)
			Expect(ilkErr).NotTo(HaveOccurred())

			Expect(len(state.Ilks)).To(Equal(10))

			var ilkCount int
			getErr := db.Get(&ilkCount, `SELECT count(*) FROM maker.ilks`)

			Expect(getErr).NotTo(HaveOccurred())
			Expect(ilkCount).To(Equal(10))
		})

		It("generates the specified number of storage records associated with ilks", func() {
			ilkErr := state.GenerateDataForIlkQueryTesting(1, 9)
			Expect(ilkErr).NotTo(HaveOccurred())

			Expect(len(state.Ilks)).To(Equal(1))

			var vatIlkRateCount int
			getErr := db.Get(&vatIlkRateCount, `SELECT count(*) FROM maker.vat_ilk_rate`)

			Expect(getErr).NotTo(HaveOccurred())
			Expect(vatIlkRateCount).To(Equal(10))
		})

		It("rolls back if there's an error inserting", func() { })
	})
})
