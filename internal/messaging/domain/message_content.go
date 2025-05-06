package domain

import "github.com/google/uuid"

// NOTE: might be redundant
type MessageContent struct {
	text      string
	mentions  []uuid.UUID
	links     []string
	formatted bool
}

func NewMessageContent(text string, mentions []uuid.UUID, links []string, formatted bool) MessageContent {
	return MessageContent{
		text:      text,
		mentions:  mentions,
		links:     links,
		formatted: formatted,
	}
}

func (mc *MessageContent) GetText() string {
	return mc.text
}

func (mc *MessageContent) setText(text string) {
	mc.text = text
}

func (mc *MessageContent) GetMentions() []uuid.UUID {
	return mc.mentions
}

func (mc *MessageContent) setMentions(mentions []uuid.UUID) {
	mc.mentions = mentions
}

func (mc *MessageContent) GetLinks() []string {
	return mc.links
}

func (mc *MessageContent) setLinks(links []string) {
	mc.links = links
}

func (mc *MessageContent) GetIsFormatted() bool {
	return mc.formatted
}

func (mc *MessageContent) setIsFormatted(formatted bool) {
	mc.formatted = formatted
}
