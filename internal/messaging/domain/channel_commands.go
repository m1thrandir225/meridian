package domain

import (
	"github.com/google/uuid"
)

type Command interface {
	CommandName() string
}

type CreateChannelCommand struct {
	Name          string
	Topic         string
	CreatorUserID uuid.UUID
}

func (c CreateChannelCommand) CommandName() string {
	return "CreateChannel"
}

type GetChannelCommand struct {
	ChannelID uuid.UUID
}

func (c GetChannelCommand) CommandName() string {
	return "GetChannel"
}

type JoinChannelCommand struct {
	ChannelID uuid.UUID
	UserID    uuid.UUID
}

func (c JoinChannelCommand) CommandName() string {
	return "JoinChannel"
}

type LeaveChannelCommand struct {
	ChannelID uuid.UUID
	UserID    uuid.UUID
}

func (c LeaveChannelCommand) CommandName() string {
	return "LeaveChanel"
}

type SendMessageCommand struct {
	ChannelID       uuid.UUID
	SenderUserID    uuid.UUID
	Content         MessageContent
	ParentMessageID *uuid.UUID
}

func (c SendMessageCommand) CommandName() string {
	return "SendMessage"
}

type SendNotificationCommand struct {
	ChannelID     uuid.UUID
	IntegrationID uuid.UUID
	Content       MessageContent
}

func (c SendNotificationCommand) CommandName() string {
	return "SendNotification"
}

type AddReactionCommand struct {
	ChannelID    uuid.UUID
	MessageID    uuid.UUID
	UserID       uuid.UUID
	ReactionType string
}

func (c AddReactionCommand) CommandName() string {
	return "AddReaction"
}

type RemoveReactionCommand struct {
	ChannelID    uuid.UUID
	MessageID    uuid.UUID
	UserID       uuid.UUID
	ReactionType string
}

func (c RemoveReactionCommand) CommandName() string {
	return "RemoveReaction"
}

type SetChannelTopicCommand struct {
	ChannelID uuid.UUID
	UserID    uuid.UUID
	Topic     string
}

func (c SetChannelTopicCommand) CommandName() string {
	return "SetChannelTopic"
}

type ArchiveChannelCommand struct {
	ChannelID uuid.UUID
	UserID    uuid.UUID
}

func (c ArchiveChannelCommand) CommandName() string {
	return "ArchiveChannel"
}

type UnarchiveChannelCommand struct {
	ChannelID uuid.UUID
	UserID    uuid.UUID
}

func (c UnarchiveChannelCommand) CommandName() string {
	return "UnarchiveChannel"
}

type ListMessagesForChannelCommand struct {
	ChannelID uuid.UUID
	Limit     int
	Offset    int
}

func (c ListMessagesForChannelCommand) CommandName() string {
	return "ListMessagesForChannel"
}

type CommandResult interface {
	IsSuccess() bool
	GetError() error
}

type commandResult struct {
	success bool
	err     error
}

func (r commandResult) IsSuccess() bool {
	return r.success
}

func (r commandResult) GetError() error {
	return r.err
}

func Success() CommandResult {
	return commandResult{success: true, err: nil}
}

func Failure(err error) CommandResult {
	return commandResult{success: false, err: err}
}

type CommandResultWithData interface {
	CommandResult
	GetData() any
}

type commandResultWithData struct {
	commandResult
	data any
}

func (r commandResultWithData) GetData() any {
	return r.data
}

func SuccessWithData(data any) CommandResultWithData {
	return commandResultWithData{
		commandResult: commandResult{success: true},
		data:          data,
	}
}
