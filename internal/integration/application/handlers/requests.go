package handlers

type RegisterIntegrationRequest struct {
	ServiceName      string   `json:"service_name" binding:"required"`
	TargetChannelIDs []string `json:"target_channel_ids" binding:"required"`
}

type RevokeIntegrationRequest struct {
	IntegrationID string `json:"integration_id" binding:"required"`
}
