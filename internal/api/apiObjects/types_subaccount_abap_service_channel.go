package apiobjects

type SubaccountABAPServiceChannel struct {
	ABAPCloudTenantHost string                            `json:"abapCloudTenantHost"`
	InstanceNumber      int64                             `json:"instanceNumber"`
	ID                  int64                             `json:"id"`
	Type                string                            `json:"type"`
	Port                int64                             `json:"port"`
	Enabled             bool                              `json:"enabled"`
	Connections         int64                             `json:"connections"`
	Comment             string                            `json:"comment"`
	State               SubaccountABAPServiceChannelState `json:"state"`
}

type SubaccountABAPServiceChannelState struct {
	Connected               bool  `json:"connected"`
	OpenedConnections       int64 `json:"openedConnections"`
	ConnectedSinceTimeStamp int64 `json:"connectedSinceTimeStamp"`
}

type SubaccountABAPServiceChannels struct {
	SubaccountABAPServiceChannels []SubaccountABAPServiceChannel `json:"service_channels_k8s"`
}
