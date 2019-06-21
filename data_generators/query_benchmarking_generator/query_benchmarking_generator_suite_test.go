package query_benchmarking_generator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueryBenchmarkingGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "QueryBenchmarkingGenerator Suite")
}
