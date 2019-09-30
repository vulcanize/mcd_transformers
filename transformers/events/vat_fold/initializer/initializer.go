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

package initializer

import (
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/transformers/events/vat_fold"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

var EventTransformerInitializer transformer.EventTransformerInitializer = shared.EventTransformer{
	Config:     shared.GetEventTransformerConfig(constants.VatFoldLabel, constants.VatFoldSignature()),
	Converter:  &vat_fold.VatFoldConverter{},
	Repository: &vat_fold.VatFoldRepository{},
}.NewEventTransformer
