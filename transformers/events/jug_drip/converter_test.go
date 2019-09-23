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

package jug_drip_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"

	"github.com/vulcanize/mcd_transformers/transformers/events/jug_drip"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Jug drip converter", func() {
	var converter = jug_drip.JugDripConverter{}
	It("returns err if log is missing topics", func() {
		badLog := types.Log{}
		_, err := converter.ToModels(constants.JugABI(), []types.Log{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to an model", func() {
		model, err := converter.ToModels(constants.JugABI(), []types.Log{test_data.EthJugDripLog})
		Expect(err).NotTo(HaveOccurred())
		Expect(model).To(Equal([]shared.InsertionModel{test_data.JugDripModel}))
	})
})
