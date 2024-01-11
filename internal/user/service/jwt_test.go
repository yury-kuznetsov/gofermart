package service

import (
	"github.com/google/uuid"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	tokenService := NewTokenService()

	testCases := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name: "Valid UUID",
			id:   uuid.New(),
		},
		{
			name:    "Empty UUID",
			id:      uuid.Nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token := tokenService.GenerateToken(tc.id)
			if tc.wantErr && token != "" {
				t.Errorf("expected empty token, got '%s'", token)
				return
			}
			if !tc.wantErr && token == "" {
				t.Errorf("expected a token, got nothing")
				return
			}

			parsedID := tokenService.GetUserID(token)
			if parsedID != tc.id {
				t.Errorf("expected UUID '%s', got '%s'", tc.id, parsedID)
			}
		})
	}
}
