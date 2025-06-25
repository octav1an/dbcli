package commands

func intPtr(i int) *int {
	return &i
}

func intsEqual(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func valOrNil(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
