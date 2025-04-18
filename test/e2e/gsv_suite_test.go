package gsv_e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGSV(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gsv Test Suite")
}
