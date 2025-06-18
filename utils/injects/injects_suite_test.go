package injects_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInjects(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Injects Suite")
}
