package handlers

type RegisterIntegrationResponse struct {
	ServiceName    string   `json:"service_name"`
	TargetChannels []string `json:"target_channels"`
	Token          string   `json:"token"`
}
