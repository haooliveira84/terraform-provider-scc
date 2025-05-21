package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/SAP/terraform-provider-scc/internal/api"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &cloudConnectorProvider{}
)

func New() provider.Provider {
	return NewWithClient(http.DefaultClient)
}

func NewWithClient(httpClient *http.Client) provider.Provider {
	return &cloudConnectorProvider{
		httpClient: httpClient,
	}
}

type cloudConnectorProvider struct {
	httpClient *http.Client
}

type cloudConnectorProviderData struct {
	InstanceURL   types.String `tfsdk:"instance_url"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	CaCertificate types.String `tfsdk:"ca_certificate"`
}

func (c *cloudConnectorProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "scc"
}

func (c *cloudConnectorProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Terraform Provider for SAP Cloud Connector allows users to manage and configure SAP Cloud Connector instances within SAP BTP (Business Technology Platform). It enables automation of connectivity between SAP BTP subaccounts and on-premise systems using Terraform.",
		Attributes: map[string]schema.Attribute{
			"instance_url": schema.StringAttribute{
				MarkdownDescription: "The URL of Cloud Connector Instance. This can also be sourced from the `SCC_INSTANCE_URL` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^https?://`), "must be a valid URL starting with http:// or https://"),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username used to connect to Cloud Connector Instance. This can also be sourced from the `SCC_USERNAME` environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password used to connect to Cloud Connector Instance. This can also be sourced from the `SCC_PASSWORD` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"ca_certificate": schema.StringAttribute{
				MarkdownDescription: "Contents of a PEM-encoded CA certificate. Use `file(\"path/to/cert.pem\")` in the provider block to read from a file. This can also be sourced from the `SCC_CA_CERTIFICATE` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (c *cloudConnectorProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config cloudConnectorProviderData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for unknowns in required fields
	if config.InstanceURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_url"),
			"Unknown Cloud Connector Instance URL",
			"The provider cannot create the Cloud Connector client as the Cloud Connector Instance URL is unknown.",
		)
	}
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Username",
			"The provider cannot create the Cloud Connector client as the username is unknown.",
		)
	}
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Password",
			"The provider cannot create the Cloud Connector client as the password is unknown.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Load values from config or fallback to environment
	instance_url := os.Getenv("SCC_INSTANCE_URL")
	username := os.Getenv("SCC_USERNAME")
	password := os.Getenv("SCC_PASSWORD")
	ca_certificate := os.Getenv("SCC_CA_CERTIFICATE")

	if !config.InstanceURL.IsNull() {
		instance_url = config.InstanceURL.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}
	if !config.CaCertificate.IsNull() {
		ca_certificate = config.CaCertificate.ValueString()
	}

	// Validate required values
	if instance_url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_url"),
			"Missing Cloud Connector Instance URL",
			"The provider cannot create the Cloud Connector client because the Cloud Connector Instance URL is empty.",
		)
	}
	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Username",
			"The provider cannot create the Cloud Connector client because the username is empty.",
		)
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Password",
			"The provider cannot create the Cloud Connector client because the password is empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the base URL
	parsedURL, err := url.Parse(instance_url)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_url"),
			"Invalid Cloud Connector Instance URL",
			fmt.Sprintf("Failed to parse the provided Cloud Connector Instance URL: %s. Error: %v", instance_url, err),
		)
		return
	}

	// Convert CA certificate to []byte only if provided
	var certBytes []byte
	if ca_certificate != "" {
		certBytes = []byte(ca_certificate)
	}

	client, err := api.NewRestApiClient(c.httpClient, parsedURL, username, password, certBytes)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Creation Failed",
			fmt.Sprintf("Failed to create Cloud Connector client: %v", err),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (c *cloudConnectorProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSubaccountsDataSource,
		NewSubaccountDataSource,
		NewSystemMappingsDataSource,
		NewSystemMappingDataSource,
		NewSystemMappingResourcesDataSource,
		NewSystemMappingResourceDataSource,
		NewDomainMappingsDataSource,
		NewDomainMappingDataSource,
		NewSubaccountServiceChannelK8SDataSource,
		NewSubaccountServiceChannelsK8SDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (c *cloudConnectorProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSubaccountResource,
		NewSystemMappingResource,
		NewSystemMappingResourceResource,
		NewDomainMappingResource,
		NewSubaccountServiceChannelK8SResource,
	}
}
