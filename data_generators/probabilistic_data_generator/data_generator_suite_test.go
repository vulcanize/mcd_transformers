package probabilistic_data_generator

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDataGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data generator suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
