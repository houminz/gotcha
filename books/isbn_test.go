package books_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SimpCosm/godemo/ginkgo/books"
)

var _ = Describe("Looking up ISBN numbers", func() {                                                   // L27
    Context("When the book can be found", func() {                                                     // L28
        It("returns the correct ISBN number", func() {                                                 // L29
            Expect(books.ISBNFor("The Chronicles of Narnia", "C.S. Lewis")).To(Equal("9780060598242")) // L30
        })                                                                                             // L31
    })                                                                                                 // L32

    Context("When the book can't be found", func() {                                                   // L33
        It("returns an error", func() {                                                                // L34
            isbn, err := books.ISBNFor("The Chronicles of Blarnia", "C.S. Lewis")                      // L35
            Expect(isbn).To(BeZero())                                                                  // L36
            Expect(err).To(HaveOccurred())                                                             // L37
        })                                                                                             // L38
    })                                                                                                 // L39
})                                                                                                     // L40
