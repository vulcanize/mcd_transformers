package shared

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_poke"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/spot"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"strings"
)

const (
	headerSql = `INSERT INTO public.headers (hash, block_number, raw, block_timestamp, eth_node_id, eth_node_fingerprint)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	nodeSql = `INSERT INTO public.eth_nodes (genesis_block, network_id, eth_node_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	txSql   = `INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
		VALUES ($1, $2, $3, $4, $5)`

	// Event data
	// TODO add event data
	// TODO add tx for events
)

var (
	node = core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}
	emptyRaw, _ = json.Marshal("nothing")
)

type GeneratorState struct {
	DB            *postgres.DB
	CurrentHeader core.Header // Current work header (Read-only everywhere except in Run)
	Ilks          []int64     // Generated ilks
	Urns          []int64     // Generated urns
	PgTx          *sqlx.Tx
}

func NewGenerator(db *postgres.DB) GeneratorState {
	return GeneratorState{
		DB:            db,
		CurrentHeader: core.Header{},
		Ilks:          []int64{},
		Urns:          []int64{},
		PgTx:          nil,
	}
}

// Creates a starting ilk and urn, with the corresponding header.
func (state *GeneratorState) InsertEthNode() (core.Node, error) {
	node := state.DB.Node
	_, nodeErr := state.PgTx.Exec(nodeSql, "GENESIS", 1, node.ID)
	if nodeErr != nil {
		return core.Node{}, fmt.Errorf("could not insert initial node: %v", nodeErr)
	}
	return node, nil
}

// Creates a new ilk, or updates a random one
func (state *GeneratorState) TouchIlks() error {
	p := rand.Float32()
	if p < 0.05 {
		return state.CreateIlk()
	} else {
		return state.updateIlk()
	}
}

func (state *GeneratorState) CreateIlk() error {
	ilkName := strings.ToUpper(test_data.AlreadySeededRandomString(7))
	hexIlk := GetHexIlk(ilkName)

	ilkId, insertIlkErr := state.insertIlk(hexIlk, ilkName)
	if insertIlkErr != nil {
		return insertIlkErr
	}

	initIlkErr := state.InsertInitialIlkData(ilkId)
	if initIlkErr != nil {
		return initIlkErr
	}

	state.Ilks = append(state.Ilks, ilkId)
	return nil
}

// Updates a random property of a randomly chosen ilk
func (state *GeneratorState) updateIlk() error {
	randomIlkId := state.Ilks[rand.Intn(len(state.Ilks))]
	blockNumber, blockHash := state.getCurrentBlockAndHash()

	var storageErr error
	var eventErr error
	p := rand.Float64()
	newValue := rand.Int()
	if p < 0.1 {
		_, storageErr = state.PgTx.Exec(vat.InsertIlkRateQuery, blockNumber, blockHash, randomIlkId, newValue)
		// Rate is changed in fold, event which isn't included in spec
	} else {
		_, storageErr = state.PgTx.Exec(vat.InsertIlkSpotQuery, blockNumber, blockHash, randomIlkId, newValue)
		_, eventErr = state.PgTx.Exec(spot_poke.InsertSpotPokeQuery,
			state.CurrentHeader.Id, randomIlkId, newValue, newValue, 0, 0, emptyRaw) // tx_idx 0 to match tx

		txErr := state.insertCurrentBlockTx()
		if txErr != nil {
			return txErr
		}
	}

	if storageErr != nil {
		return storageErr
	}
	if eventErr != nil {
		return eventErr
	}

	return nil
}

func (state *GeneratorState) TouchUrns() error {
	p := rand.Float32()
	if p < 0.1 {
		return state.CreateUrnAssociatedWithRandomIlk()
	} else {
		return state.updateUrn()
	}
}

// Creates a new urn associated with the given ilk
func (state *GeneratorState) CreateUrn(ilkId int64) error {
	guy := GetRandomAddress()
	urnId, insertUrnErr := state.insertUrn(ilkId, guy)
	if insertUrnErr != nil {
		return insertUrnErr
	}

	insertUrnInitialErr := state.InsertInitialUrnData(urnId, guy)
	if insertUrnInitialErr != nil {
		return insertUrnInitialErr
	}

	state.Urns = append(state.Urns, urnId)
	return nil
}

// Creates a new urn associated with a random ilk
func (state *GeneratorState) CreateUrnAssociatedWithRandomIlk() error {
	randomIlkId := state.Ilks[rand.Intn(len(state.Ilks))]

	err := state.CreateUrn(randomIlkId)
	if err != nil {
		return err
	}

	return nil
}

func (state *GeneratorState) InsertInitialUrnData(urnId int64, guy string) error {
	blockNumber := state.CurrentHeader.BlockNumber
	blockHash := state.CurrentHeader.Hash
	ink := rand.Int()
	art := rand.Int()
	_, artErr := state.PgTx.Exec(vat.InsertUrnArtQuery, blockNumber, blockHash, urnId, art)
	_, inkErr := state.PgTx.Exec(vat.InsertUrnInkQuery, blockNumber, blockHash, urnId, ink)
	_, frobErr := state.PgTx.Exec(vat_frob.InsertVatFrobQuery,
		state.CurrentHeader.Id, urnId, guy, guy, ink, art, emptyRaw, 0, 0) // txIx 0 to match tx

	if artErr != nil || inkErr != nil || frobErr != nil {
		return fmt.Errorf("error creating urn.\n artErr: %v\ninkErr: %v\nfrobErr: %v", artErr, inkErr, frobErr)
	}

	txErr := state.insertCurrentBlockTx()
	if txErr != nil {
		return fmt.Errorf("error creating matching tx: %v", txErr)
	}
	return nil
}

// Updates ink or art on a random urn
func (state *GeneratorState) updateUrn() error {
	randomUrnId := state.Urns[rand.Intn(len(state.Urns))]
	blockNumber := state.CurrentHeader.BlockNumber
	blockHash := state.CurrentHeader.Hash
	randomGuy := GetRandomAddress()
	newValue := rand.Int()

	// Computing correct diff complicated, also getting correct guy :(

	var updateErr error
	var frobErr error
	p := rand.Float32()
	if p < 0.5 {
		// Update ink
		_, updateErr = state.PgTx.Exec(vat.InsertUrnInkQuery, blockNumber, blockHash, randomUrnId, newValue)
		_, frobErr = state.PgTx.Exec(vat_frob.InsertVatFrobQuery,
			state.CurrentHeader.Id, randomUrnId, randomGuy, randomGuy, newValue, 0, emptyRaw, 0, 0) // txIx 0 to match tx
	} else {
		// Update art
		_, updateErr = state.PgTx.Exec(vat.InsertUrnArtQuery, blockNumber, blockHash, randomUrnId, newValue)
		_, frobErr = state.PgTx.Exec(vat_frob.InsertVatFrobQuery,
			state.CurrentHeader.Id, randomUrnId, randomGuy, randomGuy, 0, newValue, emptyRaw, 0, 0) // txIx 0 to match tx
	}

	if updateErr != nil {
		return updateErr
	}
	if frobErr != nil {
		return frobErr
	}

	txErr := state.insertCurrentBlockTx()
	if txErr != nil {
		return txErr
	}

	return nil
}

// Inserts into `urns` table, returning the urn_id from the database
func (state *GeneratorState) insertUrn(ilkId int64, guy string) (int64, error) {
	var id int64
	err := state.PgTx.QueryRow(shared.InsertUrnQuery, guy, ilkId).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting urn: %v", err)
	}
	return id, nil
}

// Inserts into `ilks` table, returning the ilk_id from the database
func (state *GeneratorState) insertIlk(hexIlk, name string) (int64, error) {
	var id int64
	err := state.PgTx.QueryRow(shared.InsertIlkQuery, hexIlk, name).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting ilk: %v", err)
	}
	return id, nil
}

// Skips initial events for everything, annoying to do individually
func (state *GeneratorState) InsertInitialIlkData(ilkId int64) error {
	blockNumber, blockHash := state.getCurrentBlockAndHash()
	intInsertions := []string{
		vat.InsertIlkRateQuery,
		vat.InsertIlkSpotQuery,
		vat.InsertIlkArtQuery,
		vat.InsertIlkLineQuery,
		vat.InsertIlkDustQuery,
		jug.InsertJugIlkDutyQuery,
		jug.InsertJugIlkRhoQuery,
		cat.InsertCatIlkChopQuery,
		cat.InsertCatIlkLumpQuery,
		spot.InsertSpotIlkMatQuery,
		spot.InsertSpotIlkPipQuery,
	}

	for _, intInsertSql := range intInsertions {
		_, err := state.PgTx.Exec(intInsertSql, blockNumber, blockHash, ilkId, rand.Int())
		if err != nil {
			return fmt.Errorf("error inserting initial ilk data: %v", err)
		}
	}
	_, flipErr := state.PgTx.Exec(cat.InsertCatIlkFlipQuery,
		blockNumber, blockHash, ilkId, test_data.AlreadySeededRandomString(10))

	if flipErr != nil {
		return fmt.Errorf("error inserting initial ilk data: %v", flipErr)
	}

	return nil
}

func (state *GeneratorState) InsertCurrentHeader() error {
	header := state.CurrentHeader
	var id int64
	err := state.PgTx.QueryRow(headerSql, header.Hash, header.BlockNumber, header.Raw, header.Timestamp, 1, node.ID).Scan(&id)
	state.CurrentHeader.Id = id
	return err
}

// Inserts a tx for the current header, with index 0. This matches the events, that are all generated with index 0
func (state *GeneratorState) insertCurrentBlockTx() error {
	txHash := GetRandomHash()
	txFrom := GetRandomAddress()
	txIndex := 0
	txTo := GetRandomAddress()
	_, txErr := state.PgTx.Exec(txSql, state.CurrentHeader.Id, txHash, txFrom, txIndex, txTo)
	return txErr
}

func (state *GeneratorState) getCurrentBlockAndHash() (int64, string) {
	return state.CurrentHeader.BlockNumber, state.CurrentHeader.Hash
}

// UTF-oblivious, names generated with alphanums anyway
func GetHexIlk(ilkName string) string {
	hexIlk := fmt.Sprintf("%x", ilkName)
	unpaddedLength := len(hexIlk)
	for i := unpaddedLength; i < 64; i++ {
		hexIlk = hexIlk + "0"
	}
	return hexIlk
}

func GetRandomAddress() string {
	hash := GetRandomHash()
	address := hash[:42]
	return address
}

func GetRandomHash() string {
	seed := test_data.AlreadySeededRandomString(10)
	hash := sha3.Sum256([]byte(seed))
	return fmt.Sprintf("0x%x", hash)
}
