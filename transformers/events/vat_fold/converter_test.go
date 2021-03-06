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

package vat_fold_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/vulcanize/mcd_transformers/transformers/events/vat_fold"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Vat fold converter", func() {
	var converter = vat_fold.VatFoldConverter{}
	It("returns err if log missing topics", func() {
		badLog := core.HeaderSyncLog{}

		_, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("converts a log with positive rate to an model", func() {
		models, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{test_data.VatFoldHeaderSyncLogWithPositiveRate})

		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(Equal([]shared.InsertionModel{test_data.VatFoldModelWithPositiveRate}))
	})

	It("converts a log with negative rate to an model", func() {
		models, err := converter.ToModels(constants.VatABI(), []core.HeaderSyncLog{test_data.VatFoldHeaderSyncLogWithNegativeRate})

		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(Equal([]shared.InsertionModel{test_data.VatFoldModelWithNegativeRate}))
	})
})
