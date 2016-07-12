package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/howler-chat/api-service/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("Service", func() {
	var server http.Handler
	var req *http.Request
	var serviceCtx *service.ServiceContext
	var resp *httptest.ResponseRecorder

	BeforeEach(func() {
		// Get our Rethink Config from our local Environment
		parser := service.ParseRethinkArgs(nil)
		// Create a new service context for our service
		serviceCtx = service.NewServiceContext(parser)
		// Create a new handler instance
		server = service.NewService(serviceCtx)
		// Record HTTP responses.
		resp = httptest.NewRecorder()
	})

	AfterEach(func() {
		serviceCtx.Stop()
	})

	Describe("Error Conditions", func() {
		Context("When requested path doesn't exist", func() {
			It("should return 404", func() {
				req, _ = http.NewRequest("GET", "/path-not-found", nil)
				server.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(404))
			})
		})
	})
})
