package model

func stringPtrOrNull(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return &s
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
