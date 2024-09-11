package crafting_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCrafting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crafting Suite")
}
