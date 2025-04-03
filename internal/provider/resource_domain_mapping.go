package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-cloudconnector/internal/api"
	apiobjects "github.com/SAP/terraform-provider-cloudconnector/internal/api/apiObjects"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &DomainMappingResource{}

func NewDomainMappingResource() resource.Resource {
	return &DomainMappingResource{}
}

type DomainMappingResource struct {
	client *api.RestApiClient
}

func (r *DomainMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_mapping"
}

func (r *DomainMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Cloud Connector Domain Mapping Resource.

__Tips:__
* You must be assigned to the following roles:
	* Administrator
	* Subaccount Administrator

__Further documentation:__
<https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/domain-mappings>`,
		Attributes: map[string]schema.Attribute{
			"region_host": schema.StringAttribute{
				MarkdownDescription: "Region Host Name.",
				Required:            true,
			},
			"subaccount": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"virtual_domain": schema.StringAttribute{
				MarkdownDescription: "Domain used on the cloud side.",
				Required:            true,
			},
			"internal_domain": schema.StringAttribute{
				MarkdownDescription: "Domain used on the on-premise side.",
				Required:            true,
			},
		},
	}
}

func (r *DomainMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.RestApiClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.RestApiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DomainMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DomainMappingConfig
	var respObj apiobjects.DomainMappings
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region_host := plan.RegionHost.ValueString()
	subaccount := plan.Subaccount.ValueString()
	internal_domain := plan.InternalDomain.ValueString()
	endpoint := fmt.Sprintf("/api/v1/configuration/subaccounts/%s/%s/domainMappings", region_host, subaccount)

	planBody := map[string]string{
		"virtualDomain":  plan.VirtualDomain.ValueString(),
		"internalDomain": plan.InternalDomain.ValueString(),
	}

	err := requestAndUnmarshal(r.client, &respObj.DomainMappings, "POST", endpoint, planBody, false)
	if err != nil {
		resp.Diagnostics.AddError("error creating the cloud connector domain mapping", err.Error())
		return
	}

	err = requestAndUnmarshal(r.client, &respObj.DomainMappings, "GET", endpoint, nil, true)
	if err != nil {
		resp.Diagnostics.AddError("error fetching the cloud connector domain mapping", err.Error())
		return
	}

	mappingRespObj, err := GetDomainMapping(respObj, internal_domain)
	if err != nil {
		resp.Diagnostics.AddError("error getting Domain Mapping", fmt.Sprintf("%s", err))
		return
	}

	responseModel, err := DomainMappingValueFrom(ctx, plan, *mappingRespObj)
	if err != nil {
		resp.Diagnostics.AddError("error mapping domain mapping value", fmt.Sprintf("%s", err))
		return
	}

	diags = resp.State.Set(ctx, responseModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DomainMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DomainMappingConfig
	var respObj apiobjects.DomainMappings
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region_host := state.RegionHost.ValueString()
	subaccount := state.Subaccount.ValueString()
	internal_domain := state.InternalDomain.ValueString()
	endpoint := fmt.Sprintf("/api/v1/configuration/subaccounts/%s/%s/domainMappings", region_host, subaccount)

	err := requestAndUnmarshal(r.client, &respObj.DomainMappings, "GET", endpoint, nil, true)
	if err != nil {
		resp.Diagnostics.AddError("error fetching the cloud connector domain mapping", err.Error())
		return
	}

	mappingRespObj, err := GetDomainMapping(respObj, internal_domain)
	if err != nil {
		resp.Diagnostics.AddError("error getting Domain Mapping", fmt.Sprintf("%s", err))
		return
	}

	responseModel, err := DomainMappingValueFrom(ctx, state, *mappingRespObj)
	if err != nil {
		resp.Diagnostics.AddError("error mapping domain mapping value", fmt.Sprintf("%s", err))
		return
	}

	diags = resp.State.Set(ctx, &responseModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DomainMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DomainMappingConfig
	var respObj apiobjects.DomainMappings

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region_host := plan.RegionHost.ValueString()
	subaccount := plan.Subaccount.ValueString()
	internal_domain := state.InternalDomain.ValueString()

	if (plan.RegionHost.ValueString() != region_host) ||
		(plan.Subaccount.ValueString() != subaccount) {
		resp.Diagnostics.AddError("error updating the cloud connector domain mapping.", "Failed to update the cloud connector domain mapping due to mismatched configuration values.")
		return
	}
	endpoint := fmt.Sprintf("/api/v1/configuration/subaccounts/%s/%s/domainMappings/%s", region_host, subaccount, internal_domain)

	planBody := map[string]string{
		"virtualDomain":  plan.VirtualDomain.ValueString(),
		"internalDomain": plan.InternalDomain.ValueString(),
	}

	err := requestAndUnmarshal(r.client, &respObj.DomainMappings, "PUT", endpoint, planBody, false)
	if err != nil {
		resp.Diagnostics.AddError("error updating the cloud connector domain mapping", err.Error())
		return
	}

	endpoint = fmt.Sprintf("/api/v1/configuration/subaccounts/%s/%s/domainMappings", region_host, subaccount)

	err = requestAndUnmarshal(r.client, &respObj.DomainMappings, "GET", endpoint, nil, true)
	if err != nil {
		resp.Diagnostics.AddError("error fetching the cloud connector domain mapping", err.Error())
		return
	}

	mappingRespObj, err := GetDomainMapping(respObj, plan.InternalDomain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error getting Domain Mapping", fmt.Sprintf("%s", err))
		return
	}

	responseModel, err := DomainMappingValueFrom(ctx, plan, *mappingRespObj)
	if err != nil {
		resp.Diagnostics.AddError("error mapping domain mapping value", fmt.Sprintf("%s", err))
		return
	}

	diags = resp.State.Set(ctx, responseModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DomainMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DomainMappingConfig
	var respObj apiobjects.DomainMapping
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region_host := state.RegionHost.ValueString()
	subaccount := state.Subaccount.ValueString()
	internal_domain := state.InternalDomain.ValueString()
	endpoint := fmt.Sprintf("/api/v1/configuration/subaccounts/%s/%s/domainMappings/%s", region_host, subaccount, internal_domain)

	err := requestAndUnmarshal(r.client, &respObj, "DELETE", endpoint, nil, false)
	if err != nil {
		resp.Diagnostics.AddError("error deleting the domain mapping", err.Error())
		return
	}

	responseModel, err := DomainMappingValueFrom(ctx, state, respObj)
	if err != nil {
		resp.Diagnostics.AddError("error mapping domain mapping value", fmt.Sprintf("%s", err))
		return
	}

	diags = resp.State.Set(ctx, responseModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// func GetDomainMapping(domainMappings apiobjects.DomainMappings, targetInternalDomain string) (*apiobjects.DomainMapping, error) {
// 	for _, mapping := range domainMappings.DomainMappings {
// 		if mapping.InternalDomain == targetInternalDomain {
// 			return &mapping, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("%s", "mapping doesn't exist")
// }
