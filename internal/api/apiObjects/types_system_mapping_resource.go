package apiobjects

type SystemMappingResource struct {
	URLPath                 string `json:"id"`
	Enabled                 bool   `json:"enabled"`
	PathOnly                bool   `json:"exactMatchOnly"`
	WebsocketUpgradeAllowed bool   `json:"websocketUpgradeAllowed"`
	CreationDate            string `json:"creationDate"`
	Description             string `json:"description"`
}

type SystemMappingResourceDataSource struct {
	RegionHost              string `json:"region_host"`
	Subaccount              string `json:"subaccount"`
	VirtualHost             string `json:"virtual_host"`
	VirtualPort             string `json:"virtual_port"`
	URLPath                 string `json:"id"`
	Enabled                 bool   `json:"enabled"`
	PathOnly                bool   `json:"exactMatchOnly"`
	WebsocketUpgradeAllowed bool   `json:"websocketUpgradeAllowed"`
	CreationDate            string `json:"creationDate"`
	Description             string `json:"description"`
}

type SystemMappingResources struct {
	RegionHost             string                  `json:"regionHost"`
	Subaccount             string                  `json:"subaccount"`
	VirtualHost            string                  `json:"virtualHost"`
	VirtualPort            string                  `json:"virtualPort"`
	SystemMappingResources []SystemMappingResource `json:"systemMappingResources"`
}
