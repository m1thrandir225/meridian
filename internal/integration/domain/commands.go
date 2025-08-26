package domain

type Command interface {
	CommandName() string
}

type RegisterIntegrationCommand struct {
	ServiceName    string
	CreatorUserID  string
	TargetChannels []string
}

func (c RegisterIntegrationCommand) CommandName() string {
	return "RegisterIntegration"
}

type RevokeTokenCommand struct {
	IntegrationID string
	RequestorID   string
}

func (c RevokeTokenCommand) CommandName() string {
	return "RevokeToken"
}

type UpvokeIntegrationCommand struct {
	IntegrationID string
	RequestorID   string
}

func (c UpvokeIntegrationCommand) CommandName() string {
	return "UpvokeIntegration"
}

type GetIntegrationCommand struct {
	IntegrationID string
}

func (c GetIntegrationCommand) CommandName() string {
	return "GetIntegration"
}

type ListIntegrationsCommand struct {
	CreatorUserID string
}

func (c ListIntegrationsCommand) CommandName() string {
	return "ListIntegrations"
}

type UpdateIntegrationCommand struct {
	IntegrationID    string
	RequestorID      string
	TargetChannelIDs []string
}

func (c UpdateIntegrationCommand) CommandName() string {
	return "UpdateIntegration"
}

type RegisterIntegrationAsBotInChannelCommand struct {
	IntegrationID string
	ChannelIDs    []string
}

func (c RegisterIntegrationAsBotInChannelCommand) CommandName() string {
	return "RegisterIntegrationAsBotInChannel"
}

type RegisterIntegrationAsWebhookCommand struct {
	IntegrationID string
	WebhookURL    string
}

func (c RegisterIntegrationAsWebhookCommand) CommandName() string {
	return "RegisterIntegrationAsWebhook"
}

type DeleteIntegrationCommand struct {
	IntegrationID string
	RequestorID   string
}

func (c DeleteIntegrationCommand) CommandName() string {
	return "DeleteIntegration"
}
