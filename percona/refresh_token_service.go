package percona

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	banking "github.com/morozovcookie/agat-banking"
	"github.com/pkg/errors"
)

// InvalidTokenTime is the time that will be set in the token_valid_until column as marker that this token was manually
// invalidated.
const InvalidTokenTime = 0

var _ banking.TokenService = (*RefreshTokenService)(nil)

// RefreshTokenService represents a service for managing token data.
type RefreshTokenService struct {
	preparer Preparer

	timer         banking.Timer
	tokenFactory  banking.TokenFactory
	secretFactory banking.SecretFactory
}

// NewRefreshTokenService returns a new instance of RefreshTokenService.
func NewRefreshTokenService(preparer Preparer) *RefreshTokenService {
	return &RefreshTokenService{
		preparer: preparer,
	}
}

// StoreToken stores a single Token.
func (svc *RefreshTokenService) StoreToken(ctx context.Context, token banking.Token) error {
	return svc.storeToken(ctx, token, banking.NanosecondsToMilliseconds(token.Expiration().UnixNano()))
}

// ExpireToken expires single Token.
// Return the new Token state after update.
func (svc *RefreshTokenService) ExpireToken(ctx context.Context, id banking.ID) (banking.Token, error) {
	token, err := svc.FindTokenByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "expire token")
	}

	if err = svc.storeToken(ctx, token, InvalidTokenTime); err != nil {
		return nil, errors.Wrap(err, "expire token")
	}

	return token, nil
}

func (svc *RefreshTokenService) storeToken(ctx context.Context, token banking.Token, validUntil int64) error {
	var (
		tokenID       = token.ID().String()
		userAccountID = token.Account().ID.String()
		issuedAt      = banking.NanosecondsToMilliseconds(token.IssuedAt().UnixNano())
		expiration    = banking.NanosecondsToMilliseconds(token.Expiration().UnixNano())
		value         = token.String()
	)

	createdAt, err := svc.timer.Time(ctx)
	if err != nil {
		return errors.Wrap(err, "store token")
	}

	query, args, err := squirrel.Insert("refresh_tokens").
		Columns("token_id", "token_user_account_id", "token_issued_at", "token_expiration",
			"token_valid_until", "token_value", "created_at").
		Values(tokenID, userAccountID, issuedAt, expiration, validUntil, value, createdAt).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "store token")
	}

	stmt, err := svc.preparer.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "store token")
	}

	defer func(ctx context.Context, stmt Stmt) {
		if closeErr := stmt.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}(ctx, stmt)

	if _, err = stmt.ExecContext(ctx, args...); err != nil {
		return errors.Wrap(err, "store token")
	}

	return nil
}

// FindTokenByID returns a single Token.
func (svc *RefreshTokenService) FindTokenByID(ctx context.Context, id banking.ID) (banking.Token, error) {
	query, args, err := squirrel.Select().
		From("refresh_tokens").
		Where(squirrel.Eq{
			"token_id": id.String(),
	}).
		Limit(1).
		ToSql()

	return nil, nil
}

// RemoveExpiredTokens removes expired tokens.
// Return tokens list after remove.
func (svc *RefreshTokenService) RemoveExpiredTokens(
	ctx context.Context,
	opts banking.FindOptions,
) (
	[]banking.Token,
	error,
) {
	tt, err := svc.findExpiredTokens(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "remove expired tokens")
	}

	if err = svc.removeExpiredTokens(ctx, tt); err != nil {
		return nil, errors.Wrap(err, "remove expired tokens")
	}

	return tt, nil
}

func (svc *RefreshTokenService) findExpiredTokens(
	ctx context.Context,
	opts banking.FindOptions,
) (
	[]banking.Token,
	error,
) {
	validUntil, err := svc.timer.Time(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	query, args, err := squirrel.Select("token_id", "token_user_account_id", "token_issued_at",
		"token_expiration").
		Distinct().
		From("refresh_tokens").
		Where(squirrel.Lt{
			"valid_until": banking.NanosecondsToMilliseconds(validUntil.UnixNano()),
		}).
		Limit(opts.Limit()).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	stmt, err := svc.preparer.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	defer func(ctx context.Context, stmt Stmt){
		if closeErr := stmt.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}(ctx, stmt)

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	defer func(rows *sql.Rows){
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}(rows)

	var ()

	for rows.Next() {

	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	return nil, nil
}

func (svc *RefreshTokenService) removeExpiredTokens(ctx context.Context, tt []banking.Token) error {
	idd := make([]string, 0, len(tt))

	for _, t := range tt {
		idd = append(idd, t.ID().String())
	}

	query, args, err := squirrel.Delete("refresh_tokens").
		Where(squirrel.Eq{
			"token_id": idd,
		}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "remove expired tokens")
	}

	stmt, err := svc.preparer.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "remove expired tokens")
	}

	defer func(ctx context.Context, stmt Stmt){
		if closeErr := stmt.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}(ctx, stmt)

	if _, err = stmt.ExecContext(ctx, args...); err != nil {
		return errors.Wrap(err, "remove expired tokens")
	}

	return nil
}
