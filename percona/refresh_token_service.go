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

	tokenBuilderCreator banking.TokenBuilderCreator
	timer               banking.Timer
}

// NewRefreshTokenService returns a new instance of RefreshTokenService.
func NewRefreshTokenService(
	preparer Preparer,
	tokenBuilderCreator banking.TokenBuilderCreator,
	timer banking.Timer,
) *RefreshTokenService {
	return &RefreshTokenService{
		preparer: preparer,

		tokenBuilderCreator: tokenBuilderCreator,
		timer:               timer,
	}
}

// StoreToken stores a single Token.
func (svc *RefreshTokenService) StoreToken(ctx context.Context, token banking.Token) error {
	return svc.storeToken(ctx, token, banking.TimeToMilliseconds(token.Until()))
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

func (svc *RefreshTokenService) storeToken(ctx context.Context, token banking.Token, until int64) error {
	var (
		tokenID       = token.ID().String()
		userAccountID = token.Account().ID.String()
		issuedAt      = banking.TimeToMilliseconds(token.IssuedAt())
		expiration    = banking.TimeToMilliseconds(token.Expiration())
		value         = token.String()
	)

	createdAt, err := svc.timer.Time(ctx)
	if err != nil {
		return errors.Wrap(err, "store token")
	}

	query, args, err := squirrel.Insert("refresh_tokens").
		Columns("token_id", "token_user_account_id", "token_issued_at", "token_expiration",
			"token_valid_until", "token_value", "created_at").
		Values(tokenID, userAccountID, issuedAt, expiration, until, value, createdAt).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "store token")
	}

	tokenStmt, err := svc.preparer.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "store token")
	}

	defer tokenStmt.Close(ctx)

	if _, err = tokenStmt.ExecContext(ctx, args...); err != nil {
		return errors.Wrap(err, "store token")
	}

	return nil
}

// FindTokenByID returns a single Token.
func (svc *RefreshTokenService) FindTokenByID(ctx context.Context, id banking.ID) (banking.Token, error) {
	query, args, err := squirrel.Select("token_id", "token_user_account_id", "token_issued_at",
		"token_expiration", "token_valid_until").
		From("refresh_tokens").
		Where(squirrel.Eq{
			"token_id": id.String(),
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "find token by id")
	}

	tokenStmt, err := svc.preparer.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "find token by id")
	}

	defer tokenStmt.Close(ctx)

	userAccountStmt, err := svc.createUserAccountStmt(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "find token by id")
	}

	defer userAccountStmt.Close(ctx)

	token, err := svc.scanTokenRow(ctx, tokenStmt.QueryRowContext(ctx, args...), userAccountStmt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(banking.ErrTokenDoesNotExist, "find token by id")
	}

	if err != nil {
		return nil, errors.Wrap(err, "find token by id")
	}

	now, err := svc.timer.Time(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "find token by id")
	}

	if now.After(token.Until()) {
		return nil, errors.Wrap(banking.ErrTokenDoesNotExist, "find token by id")
	}

	return token, nil
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
	until, err := svc.timer.Time(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	query, args, err := squirrel.Select("token_id", "token_user_account_id", "token_issued_at",
		"token_expiration", "token_valid_until").
		Distinct().
		From("refresh_tokens").
		Where(squirrel.Lt{
			"token_valid_until": banking.TimeToMilliseconds(until),
		}).
		Limit(opts.Limit()).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	tokenStmt, err := svc.preparer.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	defer tokenStmt.Close(ctx)

	rows, err := tokenStmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	defer rows.Close()

	userAccountStmt, err := svc.createUserAccountStmt(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	defer userAccountStmt.Close(ctx)

	tt, err := svc.scanTokenRows(ctx, rows, userAccountStmt, opts)
	if err != nil {
		return nil, errors.Wrap(err, "find expired tokens")
	}

	return tt, nil
}

func (svc *RefreshTokenService) scanTokenRows(
	ctx context.Context,
	rows *sql.Rows,
	userAccountStmt Stmt,
	opts banking.FindOptions,
) (
	[]banking.Token,
	error,
) {
	tt := make([]banking.Token, 0, opts.Limit())

	for rows.Next() {
		token, err := svc.scanTokenRow(ctx, rows, userAccountStmt)
		if err != nil {
			return nil, errors.Wrap(err, "scan token rows")
		}

		tt = append(tt, token)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "scan token rows")
	}

	return tt, nil
}

func (svc *RefreshTokenService) scanTokenRow(
	ctx context.Context,
	scanner squirrel.RowScanner,
	userAccountStmt Stmt,
) (
	banking.Token,
	error,
) {
	var (
		jti   string
		sub   string
		iat   int64
		exp   int64
		until int64
	)

	if err := scanner.Scan(&jti, &sub, &iat, &exp, &until); err != nil {
		return nil, errors.Wrap(err, "scan token row")
	}

	account, err := fetchUserAccount(ctx, userAccountStmt, sub)
	if err != nil {
		return nil, errors.Wrap(err, "scan token row")
	}

	token, err := svc.tokenBuilderCreator.CreateTokenBuilder(ctx).
		WithID(banking.ID(jti)).
		WithAccount(account).
		WithIssuedAt(banking.MillisecondsToTime(iat)).
		WithExpiration(banking.MillisecondsToTime(exp)).
		WithValidUntil(banking.MillisecondsToTime(until)).
		Build(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "scan token row")
	}

	return token, nil
}

func (svc *RefreshTokenService) createUserAccountStmt(ctx context.Context) (Stmt, error) {
	query, _, err := squirrel.Select("account_id", "username", "email_address", "password_hash", "user_id").
		From("user_accounts").
		Where(squirrel.Expr("account_id = ?")).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "create user accounts stmt")
	}

	return svc.preparer.PrepareContext(ctx, query)
}

func fetchUserAccount(ctx context.Context, stmt Stmt, accountID string) (*banking.UserAccount, error) {
	var userID string

	account := new(banking.UserAccount)
	account.ID = banking.ID(accountID)

	err := stmt.QueryRowContext(ctx, accountID).Scan(&account.UserName, &account.EmailAddress, &account.PasswordHash,
		&userID)
	if err != nil {
		return nil, errors.Wrap(err, "find token by id")
	}

	account.User = new(banking.User)
	account.User.ID = banking.ID(userID)

	return account, nil
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

	tokenStmt, err := svc.preparer.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "remove expired tokens")
	}

	defer tokenStmt.Close(ctx)

	if _, err = tokenStmt.ExecContext(ctx, args...); err != nil {
		return errors.Wrap(err, "remove expired tokens")
	}

	return nil
}
