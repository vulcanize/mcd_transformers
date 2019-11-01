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

package ilk_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var _ = Describe("Vat file ilk converter", func() {
	var converter = ilk.VatFileIlkConverter{}
	It("returns err if log is missing topics", func() {
		badLog := core.HeaderSyncLog{
			Log: types.Log{
				Data: []byte{1, 1, 1, 1, 1},
			}}

		_, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{badLog})
		Expect(err).To(HaveOccurred())
	})

	Describe("when log is valid", func() {
		It("converts to model with data converted to ray when what is 'spot'", func() {
			models, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{test_data.VatFileIlkSpotHeaderSyncLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(models).To(Equal([]shared.InsertionModel{test_data.VatFileIlkSpotModel()}))
		})

		It("converts to model with data converted to wad when what is 'line'", func() {
			models, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{test_data.VatFileIlkLineHeaderSyncLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models).To(Equal([]shared.InsertionModel{test_data.VatFileIlkLineModel()}))
		})

		It("converts to model with data converted to rad when what is 'dust'", func() {
			models, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{test_data.VatFileIlkDustHeaderSyncLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models).To(Equal([]shared.InsertionModel{test_data.VatFileIlkDustModel()}))
		})
	})
})
