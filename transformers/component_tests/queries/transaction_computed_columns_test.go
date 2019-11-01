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

package queries

import (
	"database/sql"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Transaction computed columns", func() {
	Describe("tx_era", func() {
		It("returns an era object for a transaction", func() {
			db := test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			defer db.Close()

			headerRepository := repositories.NewHeaderRepository(db)
			fakeBlock := rand.Int()
			fakeHeader := fakes.GetFakeHeaderWithTimestamp(int64(111111111), int64(fakeBlock))
			headerID, insertHeaderErr := headerRepository.CreateOrUpdateHeader(fakeHeader)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			txFrom := "fromAddress"
			txTo := "toAddress"
			txIndex := rand.Intn(10)
			_, insertTxErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
                VALUES ($1, $2, $3, $4, $5)`, headerID, fakeHeader.Hash, txFrom, txIndex, txTo)
			Expect(insertTxErr).NotTo(HaveOccurred())

			var actualEra Era
			getEraErr := db.Get(&actualEra, `SELECT * FROM api.tx_era(
                    (SELECT (txs.hash, txs.tx_index, h.block_number, h.hash, txs.tx_from, txs.tx_to)::api.tx
			        FROM header_sync_transactions txs
			        LEFT JOIN headers h ON h.id = txs.header_id)
			    )`)
			Expect(getEraErr).NotTo(HaveOccurred())

			expectedEra := Era{
				Epoch: fakeHeader.Timestamp,
				Iso:   "1973-07-10T00:11:51Z", // Z for Zulu, meaning UTC
			}
			Expect(actualEra).To(Equal(expectedEra))
		})
	})
})

type Era struct {
	Epoch string
	Iso   string
}

type Tx struct {
	TransactionHash  sql.NullString `db:"transaction_hash"`
	TransactionIndex sql.NullInt64  `db:"transaction_index"`
	BlockHeight      sql.NullInt64  `db:"block_height"`
	BlockHash        sql.NullString `db:"block_hash"`
	TxFrom           sql.NullString `db:"tx_from"`
	TxTo             sql.NullString `db:"tx_to"`
}
