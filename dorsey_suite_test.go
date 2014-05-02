package dorsey

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDorsey(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dorsey Suite")
}
