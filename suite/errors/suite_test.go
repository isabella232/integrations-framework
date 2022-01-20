package errors_test

//revive:disable:dot-imports
import (
	"testing"

	"github.com/smartcontractkit/integrations-framework/utils"

	. "github.com/onsi/ginkgo/v2"
)

func Test_Suite(t *testing.T) {
	utils.GinkgoSuite()
	RunSpecs(t, "Geth Errors")
}