package provider

import (
	"context"
	"strings"

	apiobjects "github.com/SAP/terraform-provider-cloudconnector/internal/api/apiObjects"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SystemMappingResourceData struct {
	ID                      types.String `tfsdk:"id"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	ExactMatchOnly          types.Bool   `tfsdk:"exact_match_only"`
	WebsocketUpgradeAllowed types.Bool   `tfsdk:"websocket_upgrade_allowed"`
	CreationDate            types.String `tfsdk:"creation_date"`
	Description             types.String `tfsdk:"description"`
}

type SystemMappingResourceConfig struct {
	RegionHost              types.String `tfsdk:"region_host"`
	Subaccount              types.String `tfsdk:"subaccount"`
	VirtualHost             types.String `tfsdk:"virtual_host"`
	VirtualPort             types.String `tfsdk:"virtual_port"`
	ID                      types.String `tfsdk:"id"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	ExactMatchOnly          types.Bool   `tfsdk:"exact_match_only"`
	WebsocketUpgradeAllowed types.Bool   `tfsdk:"websocket_upgrade_allowed"`
	CreationDate            types.String `tfsdk:"creation_date"`
	Description             types.String `tfsdk:"description"`
}

type SystemMappingResourcesConfig struct {
	RegionHost             types.String                `tfsdk:"region_host"`
	Subaccount             types.String                `tfsdk:"subaccount"`
	VirtualHost            types.String                `tfsdk:"virtual_host"`
	VirtualPort            types.String                `tfsdk:"virtual_port"`
	SystemMappingResources []SystemMappingResourceData `tfsdk:"system_mapping_resources"`
}

func SystemMappingResourceValueFrom(ctx context.Context, plan SystemMappingResourceConfig, value apiobjects.SystemMappingResource) (SystemMappingResourceConfig, error) {
	model := &SystemMappingResourceConfig{
		RegionHost:              plan.RegionHost,
		Subaccount:              plan.Subaccount,
		VirtualHost:             plan.VirtualHost,
		VirtualPort:             plan.VirtualPort,
		ID:                      types.StringValue(value.ID),
		Enabled:                 types.BoolValue(value.Enabled),
		ExactMatchOnly:          types.BoolValue(value.ExactMatchOnly),
		WebsocketUpgradeAllowed: types.BoolValue(value.WebsocketUpgradeAllowed),
		CreationDate:            types.StringValue(value.CreationDate),
		Description:             types.StringValue(value.Description),
	}

	return *model, nil
}

func SystemMappingResourcesValueFrom(ctx context.Context, plan SystemMappingResourcesConfig, value apiobjects.SystemMappingResources) (SystemMappingResourcesConfig, error) {
	system_mapping_resources := []SystemMappingResourceData{}
	for _, smr := range value.SystemMappingResources {
		r := SystemMappingResourceData{
			ID:                      types.StringValue(smr.ID),
			Enabled:                 types.BoolValue(smr.Enabled),
			ExactMatchOnly:          types.BoolValue(smr.ExactMatchOnly),
			WebsocketUpgradeAllowed: types.BoolValue(smr.WebsocketUpgradeAllowed),
			CreationDate:            types.StringValue(smr.CreationDate),
			Description:             types.StringValue(smr.Description),
		}
		system_mapping_resources = append(system_mapping_resources, r)
	}

	model := &SystemMappingResourcesConfig{
		RegionHost:             plan.RegionHost,
		Subaccount:             plan.Subaccount,
		VirtualHost:            plan.VirtualHost,
		VirtualPort:            plan.VirtualPort,
		SystemMappingResources: system_mapping_resources,
	}

	return *model, nil
}

/*
CreateEncodedResourceID encodes the given resource ID to make it safe for use in a URI path.

According to the encoding rules, it replaces specific characters to avoid collisions:
- '+' is replaced with '+2B'
- '-' is replaced with '+2D'
- '/' is replaced with '-'

This ensures the resource ID can be safely used in URI paths without misinterpretation.
*/
func CreateEncodedResourceID(input string) (encodedResourceID string) {
	input = strings.ReplaceAll(input, "+", "+2B")
	input = strings.ReplaceAll(input, "-", "+2D")
	input = strings.ReplaceAll(input, "/", "-")

	return input
}
