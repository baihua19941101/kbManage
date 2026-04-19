package backuprestore

import "context"

func (s *Service) ListPolicyViews(ctx context.Context, userID uint64, filter PolicyListFilter) ([]any, error) {
	items, err := s.ListPolicies(ctx, userID, filter)
	if err != nil {
		return nil, err
	}
	out := make([]any, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	return out, nil
}
