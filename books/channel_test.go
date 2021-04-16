package books_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SimpCosm/godemo/ginkgo/books"
)

var _ = Describe("Channel", func() {
    It("panics in a goroutine", func(done Done) {
    go func() {
        defer GinkgoRecover()

        Ω(doSomething()).Should(BeTrue())

        close(done)
    }()
})
