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
	var serviceCtx *service.ServiceContext
	var err error

	Describe("service un-available", func() {
		BeforeEach(func() {
			cmdLine := []string{"endpoints", "http://unknown-host:8000"}
			// Get our Rethink Config from our local Environment
			parser := service.ParseRethinkArgs(&cmdLine)
			// Create a new service context for our service
			serviceCtx = service.NewServiceContext(parser)
			// Create a new instance
			server = httptest.NewServer(service.NewService(serviceCtx))
			// New Instance of the client
			client, err = service.NewServiceClient(server.URL)
			if err != nil {
				Fail(err.Error())
			}

		})

		AfterEach(func() {
			serviceCtx.Stop()
			server.Close()
		})

		Context("When api service is not connected to rethinkdb", func() {
			It("service should return code 503", func() {
				msg, err := client.GetMessage(context.Background(), "non-existant", "non-existant")
				Expect(msg).To(BeNil())
				Expect(err).To(Not(BeNil()))
				Expect(service.GetErrorMsg(err)).
					To(Equal("Rethinkdb Error - gorethink: the connection is closed"))
				Expect(service.GetErrorCode(err)).To(Equal(503))
			})
		})
	})

	/*Describe("/api", func() {
		BeforeEach(func() {
			// Get our Rethink Config from our local Environment
			parser := service.ParseRethinkArgs(nil)
			// Create a rethink factory for our service
			factory = rethink.NewFactory(parser)
			// Create a new instance
			server = httptest.NewServer(service.NewService(factory))
			// New Instance of the client
			client, err = service.NewServiceClient(server.URL)
			if err != nil {
				Fail(err.Error())
			}
		})

		AfterEach(func() {
			factory.Close()
			server.Close()
		})

		Describe("/message.get", func() {
			Context("When requested messageId and channelId doesn't exist", func() {
				It("should return code 404", func() {
					msg, err := client.GetMessage(context.Background(), "non-existant", "non-existant")
					Expect(msg).To(BeNil())
					Expect(err).To(Not(BeNil()))
					Expect(string(service.GetErrorRaw(err))).To(Equal(""))
					Expect(service.GetErrorMsg(err)).To(Equal(""))
					Expect(service.GetErrorCode(err)).To(Equal(404))
				})
			})
		})
	})*/
})
