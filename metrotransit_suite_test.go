package metrotransit_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMetrotransit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metrotransit Suite")
}
