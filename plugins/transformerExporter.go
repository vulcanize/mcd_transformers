// This is a plugin generated to export the configured transformer initializers

package main

import (
	bite "github.com/vulcanize/mcd_transformers/transformers/events/bite/initializer"
	cat_file_chop_lump "github.com/vulcanize/mcd_transformers/transformers/events/cat_file/chop_lump/initializer"
	cat_file_flip "github.com/vulcanize/mcd_transformers/transformers/events/cat_file/flip/initializer"
	cat_file_vow "github.com/vulcanize/mcd_transformers/transformers/events/cat_file/vow/initializer"
	deal "github.com/vulcanize/mcd_transformers/transformers/events/deal/initializer"
	dent "github.com/vulcanize/mcd_transformers/transformers/events/dent/initializer"
	flap_kick "github.com/vulcanize/mcd_transformers/transformers/events/flap_kick/initializer"
	flip_kick "github.com/vulcanize/mcd_transformers/transformers/events/flip_kick/initializer"
	flop_kick "github.com/vulcanize/mcd_transformers/transformers/events/flop_kick/initializer"
	jug_drip "github.com/vulcanize/mcd_transformers/transformers/events/jug_drip/initializer"
	jug_file_base "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/base/initializer"
	jug_file_ilk "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/ilk/initializer"
	jug_file_vow "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/vow/initializer"
	jug_init "github.com/vulcanize/mcd_transformers/transformers/events/jug_init/initializer"
	new_cdp "github.com/vulcanize/mcd_transformers/transformers/events/new_cdp/initializer"
	spot_file_mat "github.com/vulcanize/mcd_transformers/transformers/events/spot_file/mat/initializer"
	spot_file_pip "github.com/vulcanize/mcd_transformers/transformers/events/spot_file/pip/initializer"
	spot_poke "github.com/vulcanize/mcd_transformers/transformers/events/spot_poke/initializer"
	tend "github.com/vulcanize/mcd_transformers/transformers/events/tend/initializer"
	tick "github.com/vulcanize/mcd_transformers/transformers/events/tick/initializer"
	vat_file_debt_ceiling "github.com/vulcanize/mcd_transformers/transformers/events/vat_file/debt_ceiling/initializer"
	vat_file_ilk "github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk/initializer"
	vat_flux "github.com/vulcanize/mcd_transformers/transformers/events/vat_flux/initializer"
	vat_fold "github.com/vulcanize/mcd_transformers/transformers/events/vat_fold/initializer"
	vat_fork "github.com/vulcanize/mcd_transformers/transformers/events/vat_fork/initializer"
	vat_frob "github.com/vulcanize/mcd_transformers/transformers/events/vat_frob/initializer"
	vat_grab "github.com/vulcanize/mcd_transformers/transformers/events/vat_grab/initializer"
	vat_heal "github.com/vulcanize/mcd_transformers/transformers/events/vat_heal/initializer"
	vat_init "github.com/vulcanize/mcd_transformers/transformers/events/vat_init/initializer"
	vat_move "github.com/vulcanize/mcd_transformers/transformers/events/vat_move/initializer"
	vat_slip "github.com/vulcanize/mcd_transformers/transformers/events/vat_slip/initializer"
	vat_suck "github.com/vulcanize/mcd_transformers/transformers/events/vat_suck/initializer"
	vow_fess "github.com/vulcanize/mcd_transformers/transformers/events/vow_fess/initializer"
	vow_file "github.com/vulcanize/mcd_transformers/transformers/events/vow_file/initializer"
	vow_flog "github.com/vulcanize/mcd_transformers/transformers/events/vow_flog/initializer"
	yank "github.com/vulcanize/mcd_transformers/transformers/events/yank/initializer"
	cat "github.com/vulcanize/mcd_transformers/transformers/storage/cat/initializer"
	cdp_manager "github.com/vulcanize/mcd_transformers/transformers/storage/cdp_manager/initializer"
	flap_storage "github.com/vulcanize/mcd_transformers/transformers/storage/flap/initializer"
	bat_flip "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/bat_flip"
	eth_flip_a "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/eth_flip_a"
	eth_flip_b "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/eth_flip_b"
	eth_flip_c "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/eth_flip_c"
	zrx_flip "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/zrx_flip"
	flop_storage "github.com/vulcanize/mcd_transformers/transformers/storage/flop/initializer"
	jug "github.com/vulcanize/mcd_transformers/transformers/storage/jug/initializer"
	spot "github.com/vulcanize/mcd_transformers/transformers/storage/spot/initializer"
	vat "github.com/vulcanize/mcd_transformers/transformers/storage/vat/initializer"
	vow "github.com/vulcanize/mcd_transformers/transformers/storage/vow/initializer"
	interface1 "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]interface1.EventTransformerInitializer, []interface1.StorageTransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []interface1.EventTransformerInitializer{flop_kick.EventTransformerInitializer, jug_file_ilk.EventTransformerInitializer, vat_move.EventTransformerInitializer, spot_poke.EventTransformerInitializer, vat_grab.EventTransformerInitializer, flip_kick.EventTransformerInitializer, cat_file_chop_lump.EventTransformerInitializer, deal.EventTransformerInitializer, jug_drip.EventTransformerInitializer, jug_init.EventTransformerInitializer, spot_file_pip.EventTransformerInitializer, vat_fork.EventTransformerInitializer, vat_frob.EventTransformerInitializer, tick.EventTransformerInitializer, vat_file_ilk.EventTransformerInitializer, vat_flux.EventTransformerInitializer, vow_flog.EventTransformerInitializer, jug_file_base.EventTransformerInitializer, vat_heal.EventTransformerInitializer, bite.EventTransformerInitializer, dent.EventTransformerInitializer, flap_kick.EventTransformerInitializer, jug_file_vow.EventTransformerInitializer, vow_file.EventTransformerInitializer, cat_file_flip.EventTransformerInitializer, vat_fold.EventTransformerInitializer, vat_init.EventTransformerInitializer, vat_slip.EventTransformerInitializer, vow_fess.EventTransformerInitializer, cat_file_vow.EventTransformerInitializer, spot_file_mat.EventTransformerInitializer, tend.EventTransformerInitializer, vat_file_debt_ceiling.EventTransformerInitializer, vat_suck.EventTransformerInitializer, yank.EventTransformerInitializer, new_cdp.EventTransformerInitializer}, []interface1.StorageTransformerInitializer{spot.StorageTransformerInitializer, bat_flip.StorageTransformerInitializer, vow.StorageTransformerInitializer, flop_storage.StorageTransformerInitializer, jug.StorageTransformerInitializer, flap_storage.StorageTransformerInitializer, eth_flip_a.StorageTransformerInitializer, eth_flip_b.StorageTransformerInitializer, eth_flip_c.StorageTransformerInitializer, cat.StorageTransformerInitializer, zrx_flip.StorageTransformerInitializer, vat.StorageTransformerInitializer, cdp_manager.StorageTransformerInitializer}, []interface1.ContractTransformerInitializer{}
}
