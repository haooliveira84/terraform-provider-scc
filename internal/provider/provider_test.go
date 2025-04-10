package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/SAP/terraform-provider-cloudconnector/validation/uuidvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	regexpValidUUID        = uuidvalidator.UuidRegexp
	regexValidTimeStamp    = regexp.MustCompile(`^\d{13}$`)
	regexValidSerialNumber = regexp.MustCompile(`^(?:[0-9a-fA-F]{2}:){15}[0-9a-fA-F]{2}$`)
)

type User struct {
	Username      string
	Password      string
	CACertificate string
}

var redactedTestUser = User{
	Username:      "Administrator",
	Password:      "Terraform",
	CACertificate: "/workspaces/terraform-provider-for-sap-cloud-connector/rootCA.pem",
}

func providerConfig(_ string, testUser User) string {
	instance_url := "https://10.52.101.149:8443"
	return fmt.Sprintf(`
	provider "cloudconnector" {
	instance_url= "%s"
	username= "%s"
	password= "%s"
	ca_certificate= file("%s")
	}
	`, instance_url, testUser.Username, testUser.Password, testUser.CACertificate)
}

func getTestProviders(httpClient *http.Client) map[string]func() (tfprotov6.ProviderServer, error) {
	cloudconnectorProvider := NewWithClient(httpClient).(*cloudConnectorProvider)

	return map[string]func() (tfprotov6.ProviderServer, error){
		"cloudconnector": providerserver.NewProtocol6WithError(cloudconnectorProvider),
	}
}

func setupVCR(t *testing.T, cassetteName string) (*recorder.Recorder, User) {
	t.Helper()

	mode := recorder.ModeRecordOnce
	if testRecord, _ := strconv.ParseBool(os.Getenv("TEST_RECORD")); testRecord {
		mode = recorder.ModeRecordOnly
	}

	user := redactedTestUser

	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       cassetteName,
		Mode:               mode,
		SkipRequestLatency: true,
		RealTransport:      http.DefaultTransport,
	})

	if rec.IsRecording() {
		t.Logf("ATTENTION: Recording '%s'", cassetteName)
		user.Username = os.Getenv("CC_USERNAME")
		user.Password = os.Getenv("CC_PASSWORD")
		user.CACertificate = os.Getenv("CC_CA_CERTIFICATE")
		if len(user.Username) == 0 || len(user.Password) == 0 || len(user.CACertificate) == 0 {
			t.Fatal("Env vars CC_USERNAME and CC_PASSWORD are required when recording test fixtures")
		}
	} else {
		t.Logf("Replaying '%s'", cassetteName)
	}

	if err != nil {
		t.Fatal()
	}

	rec.SetMatcher(requestMatcher(t))
	rec.AddHook(hookRedactSensitiveHeaders(), recorder.BeforeSaveHook)
	rec.AddHook(hookRedactBodyLinks(), recorder.BeforeSaveHook)
	rec.AddHook(hookRedactSensitiveBody(), recorder.BeforeSaveHook)

	return rec, user
}

func requestMatcher(t *testing.T) cassette.MatcherFunc {
	return func(r *http.Request, i cassette.Request) bool {
		if r.Method != i.Method || r.URL.String() != i.URL {
			return false
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal("Unable to read body from request")
		}
		requestBody := string(bytes)
		return requestBody == i.Body
	}
}

func hookRedactSensitiveHeaders() func(i *cassette.Interaction) error {
	return func(i *cassette.Interaction) error {
		redact := func(headers map[string][]string) {
			for key := range headers {
				if strings.Contains(strings.ToLower(key), "x-csrf-token") ||
					strings.Contains(strings.ToLower(key), "set-cookie") ||
					strings.Contains(strings.ToLower(key), "authorization") {
					headers[key] = []string{"redacted"}
				}
			}
		}

		redact(i.Request.Headers)
		redact(i.Response.Headers)

		return nil
	}
}

func hookRedactSensitiveBody() func(i *cassette.Interaction) error {
	return func(i *cassette.Interaction) error {
		if strings.Contains(i.Request.Body, "cloudPassword") {
			reBindingSecret := regexp.MustCompile(`cloudPassword":"(.*?)"`)
			i.Request.Body = reBindingSecret.ReplaceAllString(i.Request.Body, `cloudPassword":"redacted"`)
		}

		if strings.Contains(i.Request.Body, "cloudUser") {
			reBindingSecret := regexp.MustCompile(`cloudUser":"(.*?)"`)
			i.Request.Body = reBindingSecret.ReplaceAllString(i.Request.Body, `cloudUser":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "subaccountCertificate") {
			reBindingSecret := regexp.MustCompile(`subaccountCertificate":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `subaccountCertificate":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "user") {
			reBindingSecret := regexp.MustCompile(`user":"(.*?)"`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `user":"redacted"`)
		}

		return nil
	}
}

func hookRedactBodyLinks() func(i *cassette.Interaction) error {
	return func(i *cassette.Interaction) error {
		if strings.Contains(i.Response.Body, "_links") {
			reBindingSecret := regexp.MustCompile(`_links":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `_links":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "Kyma-channels") {
			reBindingSecret := regexp.MustCompile(`Kyma-channels":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `Kyma-channels":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "HANA-channels") {
			reBindingSecret := regexp.MustCompile(`HANA-channels":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `HANA-channels":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "systemMappings") {
			reBindingSecret := regexp.MustCompile(`systemMappings":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `systemMappings":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "VirtualMachine-channels") {
			reBindingSecret := regexp.MustCompile(`VirtualMachine-channels":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `VirtualMachine-channels":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "domainMappings") {
			reBindingSecret := regexp.MustCompile(`domainMappings":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `domainMappings":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "self") {
			reBindingSecret := regexp.MustCompile(`self":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `self":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "state") {
			reBindingSecret := regexp.MustCompile(`state":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `state":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "validity") {
			reBindingSecret := regexp.MustCompile(`validity":{"(.*?)"}`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `validity":"redacted"`)
		}

		return nil
	}
}

func stopQuietly(rec *recorder.Recorder) {
	if err := rec.Stop(); err != nil {
		panic(err)
	}
}

func TestCCProvider_AllResources(t *testing.T) {

	expectedResources := []string{
		"cloudconnector_domain_mapping",
		"cloudconnector_subaccount",
		"cloudconnector_system_mapping_resource",
		"cloudconnector_system_mapping",
		"cloudconnector_subaccount_k8s_service_channel",
	}

	ctx := context.Background()
	registeredResources := []string{}

	for _, resourceFunc := range New().Resources(ctx) {
		var resp resource.MetadataResponse

		resourceFunc().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "cloudconnector"}, &resp)

		registeredResources = append(registeredResources, resp.TypeName)
	}

	assert.ElementsMatch(t, expectedResources, registeredResources)
}

func TestCCProvider_AllDataSources(t *testing.T) {

	expectedDataSources := []string{
		"cloudconnector_domain_mapping",
		"cloudconnector_domain_mappings",
		"cloudconnector_subaccount_configuration",
		"cloudconnector_subaccounts",
		"cloudconnector_system_mapping_resource",
		"cloudconnector_system_mapping_resources",
		"cloudconnector_system_mapping",
		"cloudconnector_system_mappings",
		"cloudconnector_subaccount_k8s_service_channel",
		"cloudconnector_subaccount_k8s_service_channels",
	}

	ctx := context.Background()
	registeredDataSources := []string{}

	for _, datasourceFunc := range New().DataSources(ctx) {
		var resp datasource.MetadataResponse

		datasourceFunc().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "cloudconnector"}, &resp)

		registeredDataSources = append(registeredDataSources, resp.TypeName)
	}

	assert.ElementsMatch(t, expectedDataSources, registeredDataSources)
}
