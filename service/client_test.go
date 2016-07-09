// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service_test

import (
	"net/http/httptest"
	"testing"

	"github.com/howler-chat/api-service/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

func TestHttpClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Client Suite")
}

var _ = Describe("HttpClient", func() {
	var client *service.ServiceClient
	var server *httptest.Server

	BeforeEach(func() {
		// Create a new instance
		server = httptest.NewServer(service.NewService())
		// New Instance of the client
		client = service.NewServiceClient(server.URL)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("/api/message.get", func() {

		Context("When requested messageId and channelId doesn't exist", func() {
			It("should return code 404", func() {
				msg, err := client.GetMessage(context.Background(), "non-existant", "non-existant")
				Expect(msg).To(BeNil())
				Expect(err.GetCode()).To(Equal(404))
				Expect(err.GetMessage()).To(Equal(""))
			})
		})
	})
})
