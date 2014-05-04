package dorsey

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/op/go-logging"
)

var _ = Describe("A server with routing", func() {
	// shut the logger up
	logging.SetBackend(logging.NewMemoryBackend(0))

	var (
		server *Server
	)

	BeforeEach(func() {
		server = &Server{Log: logging.MustGetLogger("dorsey")}
	})

	Context("and a string path", func() {
		It("should be able to match it", func() {
			server.Get("/path", func(w *ResponseWriter, r *Request) {
				w.Render("hi")
			})

			server.Get("/anotherpath/thats/long", func(w *ResponseWriter, r *Request) {
				w.Render("anotherhi")
			})

			w := server.FakeRequest("GET", "/path")
			Expect(w.Body.String()).To(Equal("hi"))

			w = server.FakeRequest("GET", "/anotherpath/thats/long")
			Expect(w.Body.String()).To(Equal("anotherhi"))
		})

		It("should be able to not match a bad path", func() {
			server.Get("/path/thats/long", func(w *ResponseWriter, r *Request) {
				w.Render("hi")
			})

			w := server.FakeRequest("GET", "/path/thats")
			Expect(w.Code).To(Equal(404))
		})
	})
})
