package setup_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Setup Suite")
}
