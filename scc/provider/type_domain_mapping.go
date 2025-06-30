package provider

import (
	"context"
	"fmt"

	apiobjects "github.com/SAP/terraform-provider-scc/internal/api/apiObjects"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DomainMappingConfig struct {
	RegionHost     types.String `tfsdk:"region_host"`
	Subaccount     types.String `tfsdk:"subaccount"`
	VirtualDomain  types.String `tfsdk:"virtual_domain"`
	InternalDomain types.String `tfsdk:"internal_domain"`
}

type DomainMapping struct {
	VirtualDomain  types.String `tfsdk:"virtual_domain"`
	InternalDomain types.String `tfsdk:"internal_domain"`
}

type DomainMappingsConfig struct {
	RegionHost     types.String    `tfsdk:"region_host"`
	Subaccount     types.String    `tfsdk:"subaccount"`
	DomainMappings []DomainMapping `tfsdk:"domain_mappings"`
}

func DomainMappingsValueFrom(ctx context.Context, plan DomainMappingsConfig, value apiobjects.DomainMappings) (DomainMappingsConfig, error) {
	domain_mappings := []DomainMapping{}
	for _, mappings := range value.DomainMappings {
		c := DomainMapping{
			VirtualDomain:  types.StringValue(mappings.VirtualDomain),
			InternalDomain: types.StringValue(mappings.InternalDomain),
		}
		domain_mappings = append(domain_mappings, c)
	}

	model := &DomainMappingsConfig{
		RegionHost:     plan.RegionHost,
		Subaccount:     plan.Subaccount,
		DomainMappings: domain_mappings,
	}

	return *model, nil
}

func DomainMappingValueFrom(ctx context.Context, plan DomainMappingConfig, value apiobjects.DomainMapping) (DomainMappingConfig, error) {
	model := &DomainMappingConfig{
		RegionHost:     plan.RegionHost,
		Subaccount:     plan.Subaccount,
		VirtualDomain:  types.StringValue(value.VirtualDomain),
		InternalDomain: types.StringValue(value.InternalDomain),
	}
	return *model, nil
}

func GetDomainMapping(domainMappings apiobjects.DomainMappings, targetInternalDomain string) (*apiobjects.DomainMapping, error) {
	for _, mapping := range domainMappings.DomainMappings {
		if mapping.InternalDomain == targetInternalDomain {
			return &mapping, nil
		}
	}
	return nil, fmt.Errorf("%s", "mapping doesn't exist")
}
