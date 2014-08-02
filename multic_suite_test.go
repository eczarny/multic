package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMultic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Multic Suite")
}
