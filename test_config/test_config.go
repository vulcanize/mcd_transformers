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

package test_config

import (
	"errors"
	"fmt"
	"os"

	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

var TestConfig *viper.Viper
var DBConfig config.Database
var TestClient config.Client
var Infura *viper.Viper
var InfuraClient config.Client
var ABIFilePath string

func init() {
	setTestConfig()
	setInfuraConfig()
	setABIPath()
}

func setTestConfig() {
	TestConfig = viper.New()
	TestConfig.SetConfigName("private")
	TestConfig.AddConfigPath("$GOPATH/src/github.com/vulcanize/mcd_transformers/environments/")
	err := TestConfig.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	ipc := TestConfig.GetString("client.ipcPath")
	hn := TestConfig.GetString("database.hostname")
	port := TestConfig.GetInt("database.port")
	name := TestConfig.GetString("database.name")
	DBConfig = config.Database{
		Hostname: hn,
		Name:     name,
		Port:     port,
	}
	TestClient = config.Client{
		IPCPath: ipc,
	}
}

func setInfuraConfig() {
	Infura = viper.New()
	Infura.SetConfigName("infura")
	Infura.AddConfigPath("$GOPATH/src/github.com/vulcanize/mcd_transformers/environments/")
	err := Infura.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	ipc := Infura.GetString("client.ipcpath")

	// If we don't have an ipc path in the config file, check the env variable
	if ipc == "" {
		Infura.BindEnv("url", "INFURA_URL")
		ipc = Infura.GetString("url")
	}
	if ipc == "" {
		log.Fatal(errors.New("infura.toml IPC path or $INFURA_URL env variable need to be set"))
	}

	InfuraClient = config.Client{
		IPCPath: ipc,
	}
}

func setABIPath() {
	gp := os.Getenv("GOPATH")
	ABIFilePath = gp + "/src/github.com/vulcanize/vulcanizedb/pkg/geth/testing/"
}

func NewTestDB(node core.Node) *postgres.DB {
	db, err := postgres.NewDB(DBConfig, node)
	if err != nil {
		panic(fmt.Sprintf("Could not create new test db: %v", err))
	}
	return db
}

func CleanTestDB(db *postgres.DB) {
	db.MustExec("DELETE FROM blocks")
	db.MustExec("DELETE FROM checked_headers")
	db.MustExec("DELETE FROM full_sync_receipts")
	db.MustExec("DELETE FROM full_sync_transactions")
	db.MustExec("DELETE FROM headers")
	db.MustExec("DELETE FROM light_sync_receipts")
	db.MustExec("DELETE FROM light_sync_transactions")
	db.MustExec("DELETE FROM log_filters")
	db.MustExec("DELETE FROM logs")
	db.MustExec("DELETE FROM maker.bite")
	db.MustExec("DELETE FROM maker.cat_file_chop_lump")
	db.MustExec("DELETE FROM maker.cat_file_flip")
	db.MustExec("DELETE FROM maker.cat_file_vow")
	db.MustExec("DELETE FROM maker.cat_flip_ilk")
	db.MustExec("DELETE FROM maker.cat_flip_ink")
	db.MustExec("DELETE FROM maker.cat_flip_urn")
	db.MustExec("DELETE FROM maker.cat_flip_tab")
	db.MustExec("DELETE FROM maker.cat_ilk_chop")
	db.MustExec("DELETE FROM maker.cat_ilk_flip")
	db.MustExec("DELETE FROM maker.cat_ilk_lump")
	db.MustExec("DELETE FROM maker.cat_live")
	db.MustExec("DELETE FROM maker.cat_nflip")
	db.MustExec("DELETE FROM maker.cat_pit")
	db.MustExec("DELETE FROM maker.cat_vat")
	db.MustExec("DELETE FROM maker.cat_vow")
	db.MustExec("DELETE FROM maker.deal")
	db.MustExec("DELETE FROM maker.dent")
	db.MustExec("DELETE FROM maker.flap_kick")
	db.MustExec("DELETE FROM maker.flip_kick")
	db.MustExec("DELETE FROM maker.flop_kick")
	db.MustExec("DELETE FROM maker.jug_drip")
	db.MustExec("DELETE FROM maker.jug_file_base")
	db.MustExec("DELETE FROM maker.jug_file_ilk")
	db.MustExec("DELETE FROM maker.jug_file_vow")
	db.MustExec("DELETE FROM maker.jug_ilk_rho")
	db.MustExec("DELETE FROM maker.jug_ilk_duty")
	db.MustExec("DELETE FROM maker.jug_base")
	db.MustExec("DELETE FROM maker.jug_vat")
	db.MustExec("DELETE FROM maker.jug_vow")
	db.MustExec("DELETE FROM maker.pip_log_value")
	db.MustExec("DELETE FROM maker.tend")
	db.MustExec("DELETE FROM maker.vat_dai")
	db.MustExec("DELETE FROM maker.vat_debt")
	db.MustExec("DELETE FROM maker.vat_file_debt_ceiling")
	db.MustExec("DELETE FROM maker.vat_file_ilk")
	db.MustExec("DELETE FROM maker.vat_flux")
	db.MustExec("DELETE FROM maker.vat_fold")
	db.MustExec("DELETE FROM maker.vat_frob")
	db.MustExec("DELETE FROM maker.vat_gem")
	db.MustExec("DELETE FROM maker.vat_grab")
	db.MustExec("DELETE FROM maker.vat_heal")
	db.MustExec("DELETE FROM maker.vat_ilk_art")
	db.MustExec("DELETE FROM maker.vat_ilk_dust")
	db.MustExec("DELETE FROM maker.vat_ilk_line")
	db.MustExec("DELETE FROM maker.vat_ilk_rate")
	db.MustExec("DELETE FROM maker.vat_ilk_spot")
	db.MustExec("DELETE FROM maker.vat_init")
	db.MustExec("DELETE FROM maker.vat_line")
	db.MustExec("DELETE FROM maker.vat_live")
	db.MustExec("DELETE FROM maker.vat_move")
	db.MustExec("DELETE FROM maker.vat_sin")
	db.MustExec("DELETE FROM maker.vat_slip")
	db.MustExec("DELETE FROM maker.vat_urn_art")
	db.MustExec("DELETE FROM maker.vat_urn_ink")
	db.MustExec("DELETE FROM maker.vat_vice_header")
	db.MustExec("DELETE FROM maker.vat_vice")
	db.MustExec("DELETE FROM maker.vow_ash")
	db.MustExec("DELETE FROM maker.vow_bump")
	db.MustExec("DELETE FROM maker.vow_cow")
	db.MustExec("DELETE FROM maker.vow_fess")
	db.MustExec("DELETE FROM maker.vow_flog")
	db.MustExec("DELETE FROM maker.vow_hump")
	db.MustExec("DELETE FROM maker.vow_row")
	db.MustExec("DELETE FROM maker.vow_sin_integer")
	db.MustExec("DELETE FROM maker.vow_sump")
	db.MustExec("DELETE FROM maker.vow_vat")
	db.MustExec("DELETE FROM maker.vow_wait")
	db.MustExec("DELETE FROM watched_contracts")
	// TODO: add ON DELETE CASCADE? otherwise these need to come after deleting tables that reference it
	db.MustExec("DELETE FROM maker.urns")
	db.MustExec("DELETE FROM maker.ilks")
}

// Returns a new test node, with the same ID
func NewTestNode() core.Node {
	return core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}
}

func NewTestBlock(blockNumber int64, repository repositories.BlockRepository) (blockId int64) {
	blockId, err := repository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
	Expect(err).NotTo(HaveOccurred())

	return blockId
}
