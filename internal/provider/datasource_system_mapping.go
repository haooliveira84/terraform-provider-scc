package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-cloudconnector/internal/api"
	apiobjects "github.com/SAP/terraform-provider-cloudconnector/internal/api/apiObjects"
	"github.com/SAP/terraform-provider-cloudconnector/internal/api/endpoints"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &SystemMappingDataSource{}

func NewSystemMappingDataSource() datasource.DataSource {
	return &SystemMappingDataSource{}
}

type SystemMappingDataSource struct {
	client *api.RestApiClient
}

func (d *SystemMappingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_mapping"
}

func (r *SystemMappingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Cloud Connector System Mapping Data Source.
				
__Tips:__
* You must be assigned to the following roles:
	* Administrator
	* Subaccount Administrator
	* Display
	* Support

__Further documentation:__
<https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/system-mappings>`,
		Attributes: map[string]schema.Attribute{
			"region_host": schema.StringAttribute{
				MarkdownDescription: "Region Host Name.",
				Required:            true,
			},
			"subaccount": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"virtual_host": schema.StringAttribute{
				MarkdownDescription: "Virtual host used on the cloud side.",
				Required:            true,
			},
			"virtual_port": schema.StringAttribute{
				MarkdownDescription: "Virtual port used on the cloud side.",
				Required:            true,
			},
			"local_host": schema.StringAttribute{
				MarkdownDescription: "Host on the on-premise side.",
				Computed:            true,
			},
			"local_port": schema.StringAttribute{
				MarkdownDescription: "Port on the on-premise side.",
				Computed:            true,
			},
			"creation_date": schema.StringAttribute{
				MarkdownDescription: "Date of creation of system mapping.",
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol used when sending requests and receiving responses, which must be one of the following values:" +
					getFormattedValueAsTableRow("protocol", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("HTTP", "") +
					getFormattedValueAsTableRow("HTTPS", "") +
					getFormattedValueAsTableRow("RFC", "") +
					getFormattedValueAsTableRow("RFCS", "") +
					getFormattedValueAsTableRow("LDAP", "") +
					getFormattedValueAsTableRow("LDAPS", "") +
					getFormattedValueAsTableRow("TCP", "") +
					getFormattedValueAsTableRow("TCPS", ""),
				Computed: true,
			},
			"backend_type": schema.StringAttribute{
				MarkdownDescription: "Type of the backend system. Valid values are:" +
					getFormattedValueAsTableRow("backend", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("abapSys", "") +
					getFormattedValueAsTableRow("netweaverCE", "") +
					getFormattedValueAsTableRow("netweaverGW", "") +
					getFormattedValueAsTableRow("applServerJava", "") +
					getFormattedValueAsTableRow("PI", "") +
					getFormattedValueAsTableRow("hana", "") +
					getFormattedValueAsTableRow("otherSAPsys", "") +
					getFormattedValueAsTableRow("nonSAPsys", ""),
				Computed: true,
			},
			"authentication_mode": schema.StringAttribute{
				MarkdownDescription: "Authentication mode to be used on the backend side, which must be one of the following:" +
					getFormattedValueAsTableRow("authentication mode", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("NONE", "") +
					getFormattedValueAsTableRow("NONE_RESTRICTED", "") +
					getFormattedValueAsTableRow("X509_GENERAL", "") +
					getFormattedValueAsTableRow("X509_RESTRICTED", "") +
					getFormattedValueAsTableRow("KERBEROS", "") +
					"The authentication modes NONE_RESTRICTED and X509_RESTRICTED prevent the Cloud Connector from sending the system certificate in any case, whereas NONE and X509_GENERAL will send the system certificate if the circumstances allow it.",
				Computed: true,
			},
			"host_in_header": schema.StringAttribute{
				MarkdownDescription: "Policy for setting the host in the response header. This property is applicable to HTTP(S) protocols only. If set, it must be one of the following strings:" +
					getFormattedValueAsTableRow("policy", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("internal/INTERNAL", "") +
					getFormattedValueAsTableRow("virtual/VIRTUAL", "") + "The default is virtual.",
				Computed: true,
			},
			"sid": schema.StringAttribute{
				MarkdownDescription: "The ID of the system.",
				Computed:            true,
			},
			"total_resources_count": schema.Int64Attribute{
				MarkdownDescription: "The total number of resources.",
				Computed:            true,
			},
			"enabled_resources_count": schema.Int64Attribute{
				MarkdownDescription: "The number of enabled resources.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for the system mapping.",
				Computed:            true,
			},
			"sap_router": schema.StringAttribute{
				MarkdownDescription: "SAP router route, required only if an SAP router is used.",
				Computed:            true,
			},
		},
	}
}

func (d *SystemMappingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.RestApiClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.RestApiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *SystemMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemMappingConfig
	var respObj apiobjects.SystemMapping
	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region_host := data.RegionHost.ValueString()
	subaccount := data.Subaccount.ValueString()
	virtual_host := data.VirtualHost.ValueString()
	virtual_port := data.VirtualPort.ValueString()
	endpoint := endpoints.GetSystemMappingEndpoint(region_host, subaccount, virtual_host, virtual_port)

	err := requestAndUnmarshal(d.client, &respObj, "GET", endpoint, nil, true)
	if err != nil {
		resp.Diagnostics.AddError("error fetching the cloud connector system mapping", err.Error())
		return
	}

	responseModel, err := SystemMappingValueFrom(ctx, data, respObj)
	if err != nil {
		resp.Diagnostics.AddError("error mapping system mapping value", fmt.Sprintf("%s", err))
		return
	}
	diags = resp.State.Set(ctx, &responseModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
