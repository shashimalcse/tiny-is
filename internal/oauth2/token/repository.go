package token

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TokenRepository interface {
	PersistToken(ctx context.Context, jti, entryId, clientId, organizationId string, createdAt, expiresAt int64) error
	DeleteToken(ctx context.Context, jti string) error
	IsTokenExists(ctx context.Context, jti string) (bool, error)
}

type tokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (r *tokenRepository) PersistToken(ctx context.Context, jti, entryId, clientId, organizationId string, createdAt, expiresAt int64) error {
	_, err := r.db.Exec("INSERT INTO token (id, entry_id, client_id, organization_id, created_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6)", jti, entryId, clientId, organizationId, createdAt, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *tokenRepository) DeleteToken(ctx context.Context, jti string) error {
	_, err := r.db.Exec("DELETE FROM token WHERE id=$1", jti)
	if err != nil {
		return err
	}
	return nil
}

func (r *tokenRepository) IsTokenExists(ctx context.Context, jti string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM token WHERE id=$1", jti)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
