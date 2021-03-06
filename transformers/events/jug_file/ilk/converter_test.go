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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/vulcanize/mcd_transformers/transformers/events/jug_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Jug file ilk converter", func() {
	var converter = ilk.JugFileIlkConverter{}
	It("returns err if log missing topics", func() {
		badLog := core.HeaderSyncLog{
			Log: types.Log{
				Topics: []common.Hash{{}},
				Data:   []byte{1, 1, 1, 1, 1},
			}}

		_, err := converter.ToModels(constants.JugABI(), []core.HeaderSyncLog{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("returns err if log missing data", func() {
		badLog := core.HeaderSyncLog{
			Log: types.Log{
				Topics: []common.Hash{{}, {}, {}, {}},
			}}

		_, err := converter.ToModels(constants.JugABI(), []core.HeaderSyncLog{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, err := converter.ToModels(constants.JugABI(), []core.HeaderSyncLog{test_data.JugFileIlkHeaderSyncLog})
		Expect(err).NotTo(HaveOccurred())
		Expect(models).To(Equal([]shared.InsertionModel{test_data.JugFileIlkModel()}))
	})
})
