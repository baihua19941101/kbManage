package identitytenancy

import "context"

func (s *Service) RevokeUserAccess(ctx context.Context, userID uint64, reason string) error {
	return s.revocations.Mark(ctx, userID, reason)
}
