package domain

import (
	"regexp"

	"github.com/google/uuid"
)

// NOTE: might be redundant
type MessageContent struct {
	text      string
	mentions  []uuid.UUID
	links     []string
	formatted bool
}

var (
	urlRegex     = regexp.MustCompile(`https?:\/\/[^\s]+`)
	mentionRegex = regexp.MustCompile(`@[\w-]+`)
)

func NewMessageContent(message string) MessageContent {
	foundLinks := urlRegex.FindAllString(message, -1)
	//foundUsernames := mentionRegex.FindAllString(message, -1)

	if foundLinks == nil {
		foundLinks = []string{}
	}
	return MessageContent{
		text:      message,
		mentions:  make([]uuid.UUID, 0),
		links:     foundLinks,
		formatted: true,
	}
}

func RehydrateMessageContent(text string, mentions []uuid.UUID, links []string, formatted bool) MessageContent {
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
