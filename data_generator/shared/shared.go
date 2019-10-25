package shared

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_poke"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
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
	insertIlkQuery = `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`
	insertUrnQuery = `INSERT INTO maker.urns (identifier, ilk_id) VALUES ($1, $2) RETURNING id`

	insertLogSql = `
		WITH insertedAddressId AS (
			INSERT INTO public.addresses (address) VALUES ('0x1234567890123456789012345678901234567890') ON CONFLICT DO NOTHING RETURNING id
		),
		selectedAddressId AS (
			SELECT id FROM public.addresses WHERE address = '0x1234567890123456789012345678901234567890'
		)
		INSERT INTO public.header_sync_logs (header_id, address) VALUES ($1, (
			SELECT id FROM insertedAddressId
			UNION
			SELECT id FROM selectedAddressId
		)) RETURNING id`
	insertVatFrobSql = `INSERT INTO maker.vat_frob (header_id, urn_id, v, w, dink, dart, log_id)
		VALUES($1, $2::NUMERIC, $3, $4, $5::NUMERIC, $6::NUMERIC, $7)
		ON CONFLICT (header_id, log_id)
		DO UPDATE SET urn_id = $2, v = $3, w = $4, dink = $5, dart = $6;`
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
)

type GeneratorState struct {
	DB            *postgres.DB
	CurrentHeader core.Header // Current work header (Read-only everywhere except in Run)
	Ilks          []int64     // Generated ilks
	Urns          []int64     // Generated urns
	PgTx          *sqlx.Tx
}

func (state *GeneratorState) InsertCurrentHeader() error {
	header := state.CurrentHeader
	var id int64
	err := state.PgTx.QueryRow(headerSql, header.Hash, header.BlockNumber, header.Raw, header.Timestamp, 1, node.ID).Scan(&id)
	state.CurrentHeader.Id = id
	return err
}

func (state *GeneratorState) TouchIlks() error {
	p := rand.Float32()
	if p < 0.05 {
		return state.CreateIlk()
	} else {
		return state.updateIlk()
	}
}

func (state *GeneratorState) TouchUrns() error {
	p := rand.Float32()
	if p < 0.1 {
		return state.CreateUrn()
	} else {
		return state.updateUrn()
	}
}

func (state *GeneratorState) InsertEthNode() (core.Node, error) {
	node := state.DB.Node
	_, nodeErr := state.PgTx.Exec(nodeSql, "GENESIS", 1, node.ID)
	if nodeErr != nil {
		return core.Node{}, fmt.Errorf("could not insert initial node: %v", nodeErr)
	}
	return node, nil
}

func (state *GeneratorState) CreateIlk() error {
	ilkName := strings.ToUpper(test_data.AlreadySeededRandomString(7))
	hexIlk := GetHexIlk(ilkName)

	ilkId, insertIlkErr := state.insertIlk(hexIlk, ilkName)
	if insertIlkErr != nil {
		return insertIlkErr
	}

	initIlkErr := state.insertInitialIlkData(ilkId)
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

	var eventErr, logErr, storageErr error
	p := rand.Float64()
	newValue := rand.Int()
	if p < 0.1 {
		_, storageErr = state.PgTx.Exec(vat.InsertIlkRateQuery, blockNumber, blockHash, randomIlkId, newValue)
		// Rate is changed in fold, event which isn't included in spec
	} else {
		_, storageErr = state.PgTx.Exec(vat.InsertIlkSpotQuery, blockNumber, blockHash, randomIlkId, newValue)
		var logID int64
		logErr = state.PgTx.QueryRow(insertLogSql, state.CurrentHeader.Id).Scan(&logID)
		_, eventErr = state.PgTx.Exec(spot_poke.InsertSpotPokeQuery,
			state.CurrentHeader.Id, randomIlkId, newValue, newValue, logID)

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
	if logErr != nil {
		return logErr
	}

	return nil
}

// Creates a new urn associated with a random ilk
func (state *GeneratorState) CreateUrn() error {
	randomIlkId := state.Ilks[rand.Intn(len(state.Ilks))]
	guy := getRandomAddress()
	urnId, insertUrnErr := state.insertUrn(randomIlkId, guy)
	if insertUrnErr != nil {
		return insertUrnErr
	}

	blockNumber := state.CurrentHeader.BlockNumber
	blockHash := state.CurrentHeader.Hash

	ink := rand.Int()
	art := rand.Int()
	_, artErr := state.PgTx.Exec(vat.InsertUrnArtQuery, blockNumber, blockHash, urnId, art)
	_, inkErr := state.PgTx.Exec(vat.InsertUrnInkQuery, blockNumber, blockHash, urnId, ink)
	var logID int64
	logErr := state.PgTx.QueryRow(insertLogSql, state.CurrentHeader.Id).Scan(&logID)
	_, frobErr := state.PgTx.Exec(insertVatFrobSql,
		state.CurrentHeader.Id, urnId, guy, guy, ink, art, logID)

	if artErr != nil || inkErr != nil || frobErr != nil || logErr != nil {
		return fmt.Errorf("error creating urn.\n artErr: %v\ninkErr: %v\nfrobErr: %v\n logErr: %v", artErr, inkErr, frobErr, logErr)
	}

	txErr := state.insertCurrentBlockTx()
	if txErr != nil {
		return fmt.Errorf("error creating matching tx: %v", txErr)
	}

	state.Urns = append(state.Urns, urnId)
	return nil
}

// Inserts into `urns` table, returning the urn_id from the database
func (state *GeneratorState) insertUrn(ilkId int64, guy string) (int64, error) {
	var id int64
	err := state.PgTx.QueryRow(insertUrnQuery, guy, ilkId).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting urn: %v", err)
	}
	state.Urns = append(state.Urns, id)
	return id, nil
}

// Updates ink or art on a random urn
func (state *GeneratorState) updateUrn() error {
	randomUrnId := state.Urns[rand.Intn(len(state.Urns))]
	blockNumber := state.CurrentHeader.BlockNumber
	blockHash := state.CurrentHeader.Hash
	randomGuy := getRandomAddress()
	newValue := rand.Int()

	// Computing correct diff complicated, also getting correct guy :(

	var frobErr, logErr, updateErr error
	p := rand.Float32()
	if p < 0.5 {
		// Update ink
		_, updateErr = state.PgTx.Exec(vat.InsertUrnInkQuery, blockNumber, blockHash, randomUrnId, newValue)
		var logID int64
		logErr = state.PgTx.QueryRow(insertLogSql, state.CurrentHeader.Id).Scan(&logID)
		_, frobErr = state.PgTx.Exec(insertVatFrobSql,
			state.CurrentHeader.Id, randomUrnId, randomGuy, randomGuy, newValue, 0, logID)
	} else {
		// Update art
		_, updateErr = state.PgTx.Exec(vat.InsertUrnArtQuery, blockNumber, blockHash, randomUrnId, newValue)
		var logID int64
		logErr = state.PgTx.QueryRow(insertLogSql, state.CurrentHeader.Id).Scan(&logID)
		_, frobErr = state.PgTx.Exec(insertVatFrobSql,
			state.CurrentHeader.Id, randomUrnId, randomGuy, randomGuy, 0, newValue, logID)
	}

	if updateErr != nil {
		return updateErr
	}
	if frobErr != nil {
		return frobErr
	}
	if logErr != nil {
		return logErr
	}

	txErr := state.insertCurrentBlockTx()
	if txErr != nil {
		return txErr
	}

	return nil
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

// Inserts into `ilks` table, returning the ilk_id from the database
func (state *GeneratorState) insertIlk(hexIlk, name string) (int64, error) {
	var id int64
	err := state.PgTx.QueryRow(insertIlkQuery, hexIlk, name).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting ilk: %v", err)
	}
	state.Ilks = append(state.Ilks, id)
	return id, nil
}

// Skips initial events for everything, annoying to do individually
func (state *GeneratorState) insertInitialIlkData(ilkId int64) error {
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

func (state *GeneratorState) getCurrentBlockAndHash() (int64, string) {
	return state.CurrentHeader.BlockNumber, state.CurrentHeader.Hash
}

// Inserts a tx for the current header, with index 0. This matches the events, that are all generated with index 0
func (state *GeneratorState) insertCurrentBlockTx() error {
	txHash := getRandomHash()
	txFrom := getRandomAddress()
	txIndex := 0
	txTo := getRandomAddress()
	_, txErr := state.PgTx.Exec(txSql, state.CurrentHeader.Id, txHash, txFrom, txIndex, txTo)
	return txErr
}

func getRandomAddress() string {
	hash := getRandomHash()
	address := hash[:42]
	return address
}

func getRandomHash() string {
	seed := test_data.AlreadySeededRandomString(10)
	hash := sha3.Sum256([]byte(seed))
	return fmt.Sprintf("0x%x", hash)
}

// Creates a new urn associated with a random ilk
//func (state *GeneratorState) CreateUrnAssociatedWithRandomIlk() error {
//	randomIlkId := state.Ilks[rand.Intn(len(state.Ilks))]
//
//	err := state.createUrn(randomIlkId)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
