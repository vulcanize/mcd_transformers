package shared_behaviors

import (
	"database/sql"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"reflect"
	"strings"
)

type StorageVariableBehaviorInputs struct {
	KeyFieldName     string
	ValueFieldName   string
	Key              string
	Value            string
	IsAMapping       bool
	StorageTableName string
	Repository       storage.Repository
	Metadata         utils.StorageValueMetadata
}

func SharedStorageRepositoryVariableBehaviors(inputs *StorageVariableBehaviorInputs) {
	Describe("Create", func() {
		var (
			repo            = inputs.Repository
			fakeBlockNumber = rand.Int()
			fakeHash        = fakes.FakeHash.Hex()
			database        = test_config.NewTestDB(test_config.NewTestNode())
		)

		BeforeEach(func() {
			test_config.CleanTestDB(database)
			repo.SetDB(database)
		})

		It("persists a record", func() {
			err := repo.Create(fakeBlockNumber, fakeHash, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			if inputs.IsAMapping == true {
				var result MappingRes
				query := fmt.Sprintf("SELECT block_number, block_hash, %s AS key, %s AS value FROM %s",
					inputs.KeyFieldName, inputs.ValueFieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeHash, inputs.Key, inputs.Value)
			} else {
				var result VariableRes
				query := fmt.Sprintf("SELECT block_number, block_hash, %s AS value FROM %s", inputs.ValueFieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeHash, inputs.Value)
			}
		})

		It("doesn't duplicate a record", func() {
			err := repo.Create(fakeBlockNumber, fakeHash, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			err = repo.Create(fakeBlockNumber, fakeHash, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM %s", inputs.StorageTableName)
			err = database.Get(&count, query)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})
}

type IlkTriggerTestInput struct {
	Repository       storage.Repository
	Metadata         utils.StorageValueMetadata
	PropertyName     string
	PropertyValueOne string
	PropertyValueTwo string
}

func SharedIlkTriggerTests(input IlkTriggerTestInput) {
	Describe("updating current_ilk_state trigger table", func() {
		var (
			repo            = input.Repository
			database        = test_config.NewTestDB(test_config.NewTestNode())
			fakeBlockNumber = rand.Int()
			fakeBlockHash   = "expected_block_hash"
			columnName      = strings.ToLower(input.PropertyName)
			setUpStateQuery = fmt.Sprintf(`INSERT INTO api.current_ilk_state (ilk_identifier, %s, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`, columnName)
			getStateQuery   = fmt.Sprintf(`SELECT ilk_identifier, %s, created, updated FROM api.current_ilk_state`, columnName)
		)

		BeforeEach(func() {
			test_config.CleanTestDB(database)
			repo.SetDB(database)
		})

		It("inserts a row for new ilk identifier", func() {
			rawTimestamp := int64(rand.Int31())
			CreateHeader(rawTimestamp, fakeBlockNumber, database)
			expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, input.Metadata, input.PropertyValueOne)
			Expect(err).NotTo(HaveOccurred())

			var ilkState test_helpers.IlkState
			queryErr := database.Get(&ilkState, getStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
			Expect(getIlkProperty(ilkState, input.PropertyName)).To(Equal(input.PropertyValueOne))
			Expect(ilkState.Created).To(Equal(expectedTime))
			Expect(ilkState.Updated).To(Equal(expectedTime))
		})

		It("updates time created if new diff is from earlier block", func() {
			rawTimestamp := int64(rand.Int31())
			CreateHeader(rawTimestamp, fakeBlockNumber, database)
			formattedTimestamp := FormatTimestamp(rawTimestamp)
			expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

			// set up old ilk state in later block
			_, insertErr := database.Exec(setUpStateQuery,
				test_helpers.FakeIlk.Identifier, input.PropertyValueOne, formattedTimestamp)
			Expect(insertErr).NotTo(HaveOccurred())

			// set up earlier header
			earlierBlockNumber := fakeBlockNumber - 1
			earlierTimestamp := rawTimestamp - 1
			CreateHeader(earlierTimestamp, earlierBlockNumber, database)
			formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
			expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

			// trigger new ilk state from earlier block
			err := repo.Create(earlierBlockNumber, fakeBlockHash, input.Metadata, input.PropertyValueTwo)
			Expect(err).NotTo(HaveOccurred())

			var ilkState test_helpers.IlkState
			queryErr := database.Get(&ilkState, getStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
			Expect(getIlkProperty(ilkState, input.PropertyName)).To(Equal(input.PropertyValueOne))
			Expect(ilkState.Created).To(Equal(expectedTimeCreated))
			Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
		})

		It("updates value and time updated if new diff is from later block", func() {
			rawTimestamp := int64(rand.Int31())
			CreateHeader(rawTimestamp, fakeBlockNumber, database)
			formattedTimestamp := FormatTimestamp(rawTimestamp)
			expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

			// set up old ilk state in earlier block
			_, insertErr := database.Exec(setUpStateQuery,
				test_helpers.FakeIlk.Identifier, input.PropertyValueOne, formattedTimestamp)
			Expect(insertErr).NotTo(HaveOccurred())

			// set up later header
			laterBlockNumber := fakeBlockNumber + 1
			laterTimestamp := rawTimestamp + 1
			CreateHeader(laterTimestamp, laterBlockNumber, database)
			formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
			expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

			// trigger new ilk state from later block
			err := repo.Create(laterBlockNumber, fakeBlockHash, input.Metadata, input.PropertyValueTwo)
			Expect(err).NotTo(HaveOccurred())

			var ilkState test_helpers.IlkState
			queryErr := database.Get(&ilkState, getStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
			Expect(getIlkProperty(ilkState, input.PropertyName)).To(Equal(input.PropertyValueTwo))
			Expect(ilkState.Created).To(Equal(expectedTimeCreated))
			Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
		})

		It("otherwise leaves row as is", func() {
			rawTimestamp := int64(rand.Int31())
			CreateHeader(rawTimestamp, fakeBlockNumber, database)
			formattedTimestamp := FormatTimestamp(rawTimestamp)
			expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

			_, insertErr := database.Exec(setUpStateQuery,
				test_helpers.FakeIlk.Identifier, input.PropertyValueOne, formattedTimestamp)
			Expect(insertErr).NotTo(HaveOccurred())

			err := repo.Create(fakeBlockNumber, fakeBlockHash, input.Metadata, input.PropertyValueTwo)
			Expect(err).NotTo(HaveOccurred())

			var ilkState test_helpers.IlkState
			queryErr := database.Get(&ilkState, getStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
			Expect(getIlkProperty(ilkState, input.PropertyName)).To(Equal(input.PropertyValueOne))
			Expect(ilkState.Created).To(Equal(expectedTime))
			Expect(ilkState.Updated).To(Equal(expectedTime))
		})
	})
}

func getIlkProperty(ilk test_helpers.IlkState, fieldName string) string {
	r := reflect.ValueOf(ilk)
	property := reflect.Indirect(r).FieldByName(fieldName)
	return property.String()
}
