package books_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SimpCosm/godemo/ginkgo/books"
)

var _ = Describe("Reading", func() {
    var book *books.Book                                                // L6

    BeforeEach(func() {                                                 // L7
        book = books.New("The Chronicles of Narnia", 300)               // L8
        Expect(book.CurrentPage()).To(Equal(1))                         // L9
        Expect(book.NumPages()).To(Equal(300))                          // L10
    })                                                                  // L11

    It("should increment the page number", func() {                     // L12
        err := book.Read(3)                                             // L13
        Expect(err).NotTo(HaveOccurred())                               // L14
        Expect(book.CurrentPage()).To(Equal(4))                         // L15
    })                                                                  // L16

    Context("when the reader finishes the book", func() {               // L17
        It("should not allow them to read more pages", func() {         // L18
            err := book.Read(300)                                       // L19
            Expect(err).NotTo(HaveOccurred())                           // L20
            Expect(book.IsFinished()).To(BeTrue())                      // L21
            err = book.Read(1)                                          // L22
            Expect(err).To(HaveOccurred())                              // L23
        })                                                              // L24
    })                                                                  // L25
})                                                                      // L26

