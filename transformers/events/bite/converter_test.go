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

package bite_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"

	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Bite Converter", func() {
	var converter = bite.BiteConverter{}

	Describe("ToEntity", func() {
		It("converts an eth log to a bite entity", func() {
			entities, err := converter.ToEntities(constants.CatABI(), []types.Log{test_data.EthBiteLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(entities)).To(Equal(1))
			entity := entities[0]
			Expect(entity).To(Equal(test_data.BiteEntity))
		})
	})

	Describe("ToModel", func() {
		It("converts a log to a Model", func() {
			models, err := converter.ToModels(constants.CatABI(), []types.Log{test_data.EthBiteLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models).To(Equal([]shared.InsertionModel{test_data.BiteModel}))
		})

		It("returns an error if converting log to entity fails", func() {
			_, err := converter.ToModels("error abi", []types.Log{test_data.EthBiteLog})

			Expect(err).To(HaveOccurred())
		})
	})
})
