package domain

import (
	"time"
)

type ChannelDTO struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Topic           string              `json:"topic"`
	CreatorUserID   string              `json:"creator_user_id"`
	CreationTime    time.Time           `json:"creation_time"`
	LastMessageTime time.Time           `json:"last_message_time"`
	IsArchived      bool                `json:"is_archived"`
	MembersCount    int                 `json:"members_count"`
	Members         []UserDTO           `json:"members"`
	IntegrationBOts []IntegrationBotDTO `json:"bots"`
}

func ToChannelDTO(channel *Channel, members []*User, integrationBots []*IntegrationBot) ChannelDTO {
	membersDTO := make([]UserDTO, len(members))
	for i, member := range members {
		membersDTO[i] = ToUserDTO(member)
	}

	integrationBotsDTO := make([]IntegrationBotDTO, len(integrationBots))

	for i, bot := range integrationBots {
		integrationBotsDTO[i] = ToIntegrationBotDTO(bot)
	}

	return ChannelDTO{
		ID:              channel.ID.String(),
		Name:            channel.Name,
		Topic:           channel.Topic,
		CreatorUserID:   channel.CreatorUserID.String(),
		CreationTime:    channel.CreationTime,
		LastMessageTime: channel.LastMessageTime,
		IsArchived:      channel.IsArchived,
		Members:         membersDTO,
		IntegrationBOts: integrationBotsDTO,
		MembersCount:    len(channel.Members),
	}
}

type MessageDTO struct {
	ID              string    `json:"id"`
	ChannelID       string    `json:"channel_id"`
	SenderUserID    *string   `json:"sender_user_id,omitempty"`
	IntegrationID   *string   `json:"integration_id,omitempty"`
	ContentText     string    `json:"content_text"`
	CreatedAt       time.Time `json:"created_at"`
	ParentMessageID *string
	SenderUser      *UserDTO           `json:"sender_user,omitempty"`
	IntegrationBot  *IntegrationBotDTO `json:"integration_bot,omitempty"`
	Reactions       []ReactionDTO      `json:"reactions,omitempty"`
}

type IntegrationBotDTO struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	CreatedAt   time.Time `json:"created_at"`
	IsRevoked   bool      `json:"is_revoked"`
}

func ToIntegrationBotDTO(integration *IntegrationBot) IntegrationBotDTO {
	return IntegrationBotDTO{
		ID:          integration.GetId().String(),
		ServiceName: integration.GetServiceName(),
		CreatedAt:   integration.GetCreatedAt(),
		IsRevoked:   integration.GetIsRevoked(),
	}
}

type UserDTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func ToUserDTO(user *User) UserDTO {
	return UserDTO{
		ID:        user.GetId().String(),
		Username:  user.GetUsername(),
		Email:     user.GetEmail(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
	}
}

func ToMessageDTO(message *Message, sender *User, integration *IntegrationBot) MessageDTO {
	var senderId, integrationId, parentId *string
	if message.GetSenderUserId() != nil {
		sId := message.GetSenderUserId().String()
		senderId = &sId
	}

	if message.GetIntegrationId() != nil {
		iId := message.GetIntegrationId().String()
		integrationId = &iId
	}

	if message.GetParentMessageId() != nil {
		pId := message.GetParentMessageId().String()
		parentId = &pId
	}

	reactionsDTO := make([]ReactionDTO, len(message.GetReactions()))
	for i, reaction := range message.GetReactions() {
		reactionsDTO[i] = ToReactionDTO(&reaction)
	}

	var integrationBot *IntegrationBotDTO
	if integration != nil {
		bot := ToIntegrationBotDTO(integration)
		integrationBot = &bot
	}

	var senderUser *UserDTO
	if sender != nil {
		user := ToUserDTO(sender)
		senderUser = &user
	}

	return MessageDTO{
		ID:              message.GetId().String(),
		ChannelID:       message.GetChannelId().String(),
		SenderUserID:    senderId,
		IntegrationID:   integrationId,
		ContentText:     message.GetContent().GetText(),
		CreatedAt:       message.GetCreatedAt(),
		ParentMessageID: parentId,
		Reactions:       reactionsDTO,
		SenderUser:      senderUser,
		IntegrationBot:  integrationBot,
	}
}

type ReactionDTO struct {
	ID           string    `json:"id"`
	MessageID    string    `json:"message_id"`
	UserID       string    `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
	Timestamp    time.Time `json:"timestamp"`
}

func ToReactionDTO(reaction *Reaction) ReactionDTO {
	return ReactionDTO{
		ID:           reaction.GetId().String(),
		MessageID:    reaction.GetMessageId().String(),
		UserID:       reaction.GetUserId().String(),
		ReactionType: reaction.GetReactionType(),
		Timestamp:    reaction.GetCreatedAt(),
	}
}
