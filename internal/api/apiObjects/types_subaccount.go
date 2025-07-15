package apiobjects

type Subaccount struct {
	RegionHost  string           `json:"regionHost"`
	Subaccount  string           `json:"subaccount"`
	LocationID  string           `json:"locationID"`
	DisplayName string           `json:"displayName"`
	Description string           `json:"description"`
	Tunnel      SubaccountTunnel `json:"tunnel"`
}

type Subaccounts struct {
	RegionHost string `tfsdk:"regionHost"`
	Subaccount string `tfsdk:"subaccount"`
	LocationID string `tfsdk:"locationID"`
}

type SubaccountsDataSource struct {
	Subaccounts []Subaccounts `json:"subaccounts"`
}

type SubaccountCertificate struct {
	NotAfterTimeStamp  int64  `json:"notAfterTimeStamp"`
	NotBeforeTimeStamp int64  `json:"notBeforeTimeStamp"`
	SubjectDN          string `json:"subjectDN"`
	Issuer             string `json:"issuer"`
	SerialNumber       string `json:"serialNumber"`
}

type SubaccountTunnel struct {
	State                   string                             `json:"state"`
	ConnectedSinceTimeStamp int64                              `json:"connectedSinceTimeStamp"`
	Connections             int64                              `json:"connections"`
	SubaccountCertificate   SubaccountCertificate              `json:"subaccountCertificate"`
	User                    string                             `json:"user"`
	ApplicationConnections  []SubaccountApplicationConnections `json:"applicationConnections"`
	ServiceChannels         []SubaccountServiceChannels        `json:"serviceChannels"`
}

type SubaccountApplicationConnections struct {
	ConnectionCount int64  `json:"connectionCount"`
	Name            string `json:"name"`
	Type            string `json:"type"`
}

type SubaccountServiceChannels struct {
	Type    string `json:"type"`
	State   string `json:"state"`
	Details string `json:"details"`
	Comment string `json:"comment"`
}

type SubaccountResource struct {
	RegionHost    string           `json:"regionHost"`
	Subaccount    string           `json:"subaccount"`
	CloudUser     string           `json:"cloudUser"`
	CloudPassword string           `json:"cloudPassword"`
	LocationID    string           `json:"locationID,omitempty"`
	DisplayName   string           `json:"displayName,omitempty"`
	Description   string           `json:"description,omitempty"`
	Tunnel        SubaccountTunnel `json:"tunnel"`
}

type SubaccountUsingAuthResource struct {
	RegionHost         string           `json:"regionHost"`
	Subaccount         string           `json:"subaccount"`
	AuthenticationData string           `json:"authenticationData"`
	LocationID         string           `json:"locationID,omitempty"`
	DisplayName        string           `json:"displayName,omitempty"`
	Description        string           `json:"description,omitempty"`
	Tunnel             SubaccountTunnel `json:"tunnel"`
}
