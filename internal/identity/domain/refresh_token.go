package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/pkg/common"
)

const refreshTokenBytes = 32

type RefreshToken struct {
	ID        uuid.UUID
	UserID    UserID
	TokenHash string
	ExpiresAt time.Time
	IsRevoked bool
	CreatedAt time.Time
	Device    string
	IPAddress string
}

func newRefreshToken(userID UserID, device, ipAddress string, validity time.Duration) (*RefreshToken, string, error) {
	rawToken, err := common.GenerateRandomString(refreshTokenBytes)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate raw refresh token: %w", err)
	}

	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	rt := &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(validity),
		IsRevoked: false,
		CreatedAt: time.Now().UTC(),
		Device:    device,
		IPAddress: ipAddress,
	}
	return rt, rawToken, nil
}

func (rt *RefreshToken) Revoke() {
	rt.IsRevoked = true
}
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && time.Now().UTC().Before(rt.ExpiresAt)
}
