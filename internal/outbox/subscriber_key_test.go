package outbox

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresSubscriberKeyRepository_CreateAndGetActive(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresSubscriberKeyRepository(db)
	keyID := uuid.New()
	now := time.Now()
	jwk := json.RawMessage(`{"kty":"RSA","kid":"k1"}`)

	mock.ExpectExec("INSERT INTO subscriber_keys").
		WithArgs(sqlmock.AnyArg(), "sub-1", "k1", jwk, SubscriberKeyActive, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(&SubscriberKey{
		ID:           keyID,
		SubscriberID: "sub-1",
		KeyID:        "k1",
		JWK:          jwk,
		Status:       SubscriberKeyActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{
		"id", "subscriber_id", "key_id", "jwk", "status", "expires_at", "created_at", "updated_at",
	}).AddRow(keyID, "sub-1", "k1", jwk, SubscriberKeyActive, nil, now, now)

	mock.ExpectQuery("SELECT id, subscriber_id").
		WithArgs("sub-1", SubscriberKeyActive, sqlmock.AnyArg()).
		WillReturnRows(rows)

	active, err := repo.GetActiveKey("sub-1")
	require.NoError(t, err)
	assert.Equal(t, "k1", active.KeyID)
	require.NoError(t, mock.ExpectationsWereMet())
}
