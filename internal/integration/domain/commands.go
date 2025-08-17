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

type GetIntegrationCommand struct {
	IntegrationID string
}

type ListIntegrationsCommand struct {
	CreatorUserID string
}

type UpdateIntegrationCommand struct {
	IntegrationID    string
	RequestorID      string
	TargetChannelIDs []string
}

type RegisterIntegrationAsBotInChannelCommand struct {
	IntegrationID string
	ChannelIDs    []string
}
