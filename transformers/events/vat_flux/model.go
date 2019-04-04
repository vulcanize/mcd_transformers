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

package vat_flux

type VatFluxModel struct {
	Ilk              string `db:"ilk_id"`
	Src              string
	Dst              string
	Wad              string
	TransactionIndex uint   `db:"tx_idx"`
	LogIndex         uint   `db:"log_idx"`
	Raw              []byte `db:"raw_log"`
}
