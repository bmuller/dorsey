package dorsey

import (
	"log"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/op/go-logging"
)

func (s *Server) FakeRequest(method, url string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	return w
}

var _ = Describe("A server", func() {
	// shut the logger up
	logging.SetBackend(logging.NewMemoryBackend(0))

	var (
		server *Server
	)

	BeforeEach(func() {
		server = &Server{Log: logging.MustGetLogger("dorsey")}
	})

	Context("with a string path", func() {
		It("should be able to match it", func() {
			server.Get("/path", func(w *ResponseWriter, r *Request) {
				w.Render("hi")
			})

			w := server.FakeRequest("GET", "/path")
			Expect(w.Body.String()).To(Equal("hi"))
		})

		It("should be able to not match a bad path", func() {
			server.Get("/path", func(w *ResponseWriter, r *Request) {
				w.Render("hi")
			})

			w := server.FakeRequest("GET", "/badpath")
			Expect(w.Code).To(Equal(404))
		})
	})

	Context("with a named path", func() {
		It("should make a new path", func() {
			server.Get("/path/with/:name/something", func(w *ResponseWriter, r *Request) {
				Expect(r.GetURLParam("name")).To(Equal("snakePlissken"))
				w.Render("")
			})

			w := server.FakeRequest("GET", "/path/with/snakePlissken/something")
			Expect(w.Code).To(Equal(200))
		})
	})
})
