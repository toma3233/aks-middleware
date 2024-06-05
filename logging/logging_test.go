package logging_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/aks-middleware/logging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogRequest", func() {
	var (
		logger    *slog.Logger
		logBuffer *bytes.Buffer
	)

	BeforeEach(func() {
		logBuffer = &bytes.Buffer{}
		logger = slog.New(slog.NewJSONHandler(logBuffer, nil))
	})

	AfterEach(func() {
		logBuffer.Reset()
	})

	Context("when method is GET and URL has nested resource type", func() {
		It("logs the correct method info for a READ operation", func() {
			parsedURL, err := url.Parse("https://management.azure.com/subscriptions/sub_id/resourceGroups/rg_name/providers/Microsoft.Storage/storageAccounts/account_name?api-version=version")
			Expect(err).To(BeNil())

			params := logging.LogRequestParams{
				Logger:    logger,
				StartTime: time.Now(),
				Request:   &http.Request{Method: "GET", URL: parsedURL},
				Response:  &http.Response{StatusCode: 200},
				Error:     nil,
			}
			expected := "GET storageaccounts - READ"
			logging.LogRequest(params)
			Expect(logBuffer.String()).To(ContainSubstring(expected))
		})
	})

	Context("when method is GET and URL has top-level resource type", func() {
		It("logs the correct method info for a LIST operation", func() {
			parsedURL, err := url.Parse("https://management.azure.com/subscriptions/sub_id/resourceGroups?api-version=version")
			Expect(err).To(BeNil())

			params := logging.LogRequestParams{
				Logger:    logger,
				StartTime: time.Now(),
				Request:   &http.Request{Method: "GET", URL: parsedURL},
				Response:  &http.Response{StatusCode: 200},
				Error:     nil,
			}
			expected := "GET resourcegroups - LIST"
			logging.LogRequest(params)
			Expect(logBuffer.String()).To(ContainSubstring(expected))
		})
	})

	Context("when method is not GET", func() {
		It("logs the correct method info without operation type", func() {
			parsedURL, err := url.Parse("https://management.azure.com/subscriptions/sub_id/resourceGroups/rg_name/providers/Microsoft.Storage/storageAccounts?api-version=version")
			Expect(err).To(BeNil())

			params := logging.LogRequestParams{
				Logger:    logger,
				StartTime: time.Now(),
				Request:   &http.Request{Method: "POST", URL: parsedURL},
				Response:  &http.Response{StatusCode: 200},
				Error:     nil,
			}
			expected := "POST storageaccounts"
			logging.LogRequest(params)
			Expect(logBuffer.String()).To(ContainSubstring(expected))
		})
	})

	Context("when sending an azcore policy req", func() {
		It("logs the correct method info without operation type", func() {
			req, err := runtime.NewRequest(context.Background(), http.MethodPost, "https://management.azure.com/subscriptions/sub_id/resourceGroups/rg_name/providers/Microsoft.Storage/storageAccounts?api-version=version")
			Expect(err).To(BeNil())

			params := logging.LogRequestParams{
				Logger:    logger,
				StartTime: time.Now(),
				Request:   req,
				Response:  &http.Response{StatusCode: 200},
				Error:     nil,
			}
			expected := "POST storageaccounts"
			logging.LogRequest(params)
			Expect(logBuffer.String()).To(ContainSubstring(expected))
		})
	})

	Context("when URL is malformed", func() {
		It("logs the correct method with the entire URL", func() {
			parsedURL, err := url.Parse("https://management.azure.com/subscriptions/sub_id/resourceGroups/rg_name/providers/Microsoft.Storage/storageAccounts/account_name/api-version=version")
			Expect(err).To(BeNil())

			params := logging.LogRequestParams{
				Logger:    logger,
				StartTime: time.Now(),
				Request:   &http.Request{Method: "GET", URL: parsedURL},
				Response:  &http.Response{StatusCode: 200},
				Error:     nil,
			}
			expected := "GET https://management.azure.com/subscriptions/sub_id/resourceGroups/rg_name/providers/Microsoft.Storage/storageAccounts/account_name/api-version=version"
			logging.LogRequest(params)
			Expect(logBuffer.String()).To(ContainSubstring(expected))
		})
	})
	Context("when the resource type is not found", func() {
		It("logs the correct method with the entire URL", func() {
			parsedURL, err := url.Parse("https://management.azure.com/subscriptions/sub_id/customResourceGroup/resource_name/providers/Microsoft.Storage/customResource/resource_name?api-version=version")
			Expect(err).To(BeNil())

			params := logging.LogRequestParams{
				Logger:    logger,
				StartTime: time.Now(),
				Request:   &http.Request{Method: "GET", URL: parsedURL},
				Response:  &http.Response{StatusCode: 200},
				Error:     nil,
			}
			expected := "GET https://management.azure.com/subscriptions/sub_id/customResourceGroup/resource_name/providers/Microsoft.Storage/customResource/resource_name - READ"
			logging.LogRequest(params)
			Expect(logBuffer.String()).To(ContainSubstring(expected))
		})
	})
	Context("when there are multiple query parameters", func() {
		It("logs the correct method with the entire URL", func() {
			parsedURL, err := url.Parse("https://management.azure.com/subscriptions/sub_id/resourceGroups/rg_name/providers/Microsoft.Storage/storageAccounts?api-version=version&param1=value1&param2=value2")
			Expect(err).To(BeNil())

			params := logging.LogRequestParams{
				Logger:    logger,
				StartTime: time.Now(),
				Request:   &http.Request{Method: "GET", URL: parsedURL},
				Response:  &http.Response{StatusCode: 200},
				Error:     nil,
			}
			expected := "GET storageaccounts - LIST"
			logging.LogRequest(params)
			Expect(logBuffer.String()).To(ContainSubstring(expected))
		})
	})
	Context("when there are zero query parameters", func() {
		It("logs the correct method with the entire URL", func() {
			parsedURL, err := url.Parse("https://management.azure.com/subscriptions/sub_id/resourceGroups/rg_name/providers/Microsoft.Storage/storageAccounts")
			Expect(err).To(BeNil())

			params := logging.LogRequestParams{
				Logger:    logger,
				StartTime: time.Now(),
				Request:   &http.Request{Method: "GET", URL: parsedURL},
				Response:  &http.Response{StatusCode: 200},
				Error:     nil,
			}
			expected := "GET https://management.azure.com/subscriptions/sub_id/resourceGroups/rg_name/providers/Microsoft.Storage/storageAccounts"
			logging.LogRequest(params)
			Expect(logBuffer.String()).To(ContainSubstring(expected))
		})
	})
})
