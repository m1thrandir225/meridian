package domain

type RegisterIntegrationCommand struct {
	ServiceName    string
	CreatorUserID  string
	TargetChannels []string
}

type RevokeTokenCommand struct {
	IntegrationID string
	RequestorID   string
}
