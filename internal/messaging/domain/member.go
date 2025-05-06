package domain

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	id       uuid.UUID
	role     string
	joinedAt time.Time
	lastRead time.Time
}

func newMember(id uuid.UUID, role string, joinedAt, lastRead time.Time) Member {
	return Member{
		id:       id,
		role:     role,
		joinedAt: joinedAt,
		lastRead: lastRead,
	}
}

func (m *Member) GetId() uuid.UUID {
	return m.id
}

func (m *Member) setId(id uuid.UUID) {
	m.id = id
}

func (m *Member) GetRole() string {
	return m.role
}

func (m *Member) setRole(role string) {
	m.role = role
}

func (m *Member) GetJoinedAt() time.Time {
	return m.joinedAt
}

func (m *Member) setJoinedAt(joinedAt time.Time) {
	m.joinedAt = joinedAt
}

func (m *Member) GetLastRead() time.Time {
	return m.lastRead
}

func (m *Member) setLastRead(lastRead time.Time) {
	m.lastRead = lastRead
}

func RehydrateMember(
	id uuid.UUID,
	role string,
	joinedAt time.Time,
	lastRead time.Time,
) Member {
	return Member{
		id:       id,
		role:     role,
		joinedAt: joinedAt,
		lastRead: lastRead,
	}
}
