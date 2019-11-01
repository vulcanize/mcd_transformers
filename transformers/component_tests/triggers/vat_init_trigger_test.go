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

package trigger_test

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_init"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"math/rand"
)

var _ = Describe("Updating ilk_state_history table", func() {
	var (
		blockOne,
		blockTwo int
		headerOne,
		headerTwo core.Header
		rawTimestampOne,
		rawTimestampTwo int64
		logIdOne            int64
		vatInitModel        shared.InsertionModel
		repo                = vat_init.VatInitRepository{}
		database            = test_config.NewTestDB(test_config.NewTestNode())
		getTimeCreatedQuery = `SELECT created FROM api.ilk_state_history ORDER BY block_number`
		insertRecordQuery   = `INSERT INTO api.ilk_state_history (ilk_identifier, block_number, created) VALUES ($1, $2, $3::TIMESTAMP)`
		insertEmptyRowQuery = `INSERT INTO api.ilk_state_history (ilk_identifier, block_number) VALUES ($1, $2)`
	)

	BeforeEach(func() {
		test_config.CleanTestDB(database)
		repo.SetDB(database)
		blockOne = rand.Int()
		blockTwo = blockOne + 1
		rawTimestampOne = int64(rand.Int31())
		rawTimestampTwo = rawTimestampOne + 1
		headerOne = CreateHeader(rawTimestampOne, blockOne, database)
		headerTwo = CreateHeader(rawTimestampTwo, blockTwo, database)
		logIdOne = test_data.CreateTestLog(headerOne.Id, database).ID
		vatInitModel = createVatInitModel(headerOne.Id, logIdOne, test_helpers.FakeIlk.Hex)
	})

	It("updates time created of all records for an ilk", func() {
		_, setupErr := database.Exec(insertEmptyRowQuery, test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber)
		Expect(setupErr).NotTo(HaveOccurred())
		expectedTimeCreated := sql.NullString{Valid: true, String: FormatTimestamp(rawTimestampOne)}

		err := repo.Create([]shared.InsertionModel{vatInitModel})
		Expect(err).NotTo(HaveOccurred())

		var ilkStates []test_helpers.IlkState
		queryErr := database.Select(&ilkStates, getTimeCreatedQuery)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(len(ilkStates)).To(Equal(1))
		Expect(ilkStates[0].Created).To(Equal(expectedTimeCreated))
	})

	It("does not update time created if old time created is not null", func() {
		_, setupErr := database.Exec(insertRecordQuery, test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber,
			FormatTimestamp(rawTimestampTwo))
		Expect(setupErr).NotTo(HaveOccurred())
		expectedTimeCreated := sql.NullString{Valid: true, String: FormatTimestamp(rawTimestampTwo)}

		err := repo.Create([]shared.InsertionModel{vatInitModel})
		Expect(err).NotTo(HaveOccurred())

		var ilkStates []test_helpers.IlkState
		queryErr := database.Select(&ilkStates, getTimeCreatedQuery)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(len(ilkStates)).To(Equal(1))
		Expect(ilkStates[0].Created).To(Equal(expectedTimeCreated))
	})

	It("does not update records with a different ilk", func() {
		_, setupErr := database.Exec(insertEmptyRowQuery, test_helpers.AnotherFakeIlk.Identifier, headerTwo.BlockNumber)
		Expect(setupErr).NotTo(HaveOccurred())
		expectedTimeCreated := sql.NullString{Valid: false, String: ""}

		err := repo.Create([]shared.InsertionModel{vatInitModel})
		Expect(err).NotTo(HaveOccurred())

		var ilkStates []test_helpers.IlkState
		queryErr := database.Select(&ilkStates, getTimeCreatedQuery)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(len(ilkStates)).To(Equal(1))
		Expect(ilkStates[0].Created).To(Equal(expectedTimeCreated))
	})
})

func createVatInitModel(headerId, logId int64, ilkHex string) shared.InsertionModel {
	vatInit := test_data.VatInitModel
	vatInit.ForeignKeyValues[constants.IlkFK] = ilkHex
	vatInit.ColumnValues[constants.HeaderFK] = headerId
	vatInit.ColumnValues[constants.LogFK] = logId
	return vatInit
}
