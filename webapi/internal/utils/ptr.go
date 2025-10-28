package utils

func SliceToSliceOfPtrs[T any](slice []T) []*T {
	r := make([]*T, len(slice))
	for i := range slice {
		r[i] = &slice[i]
	}
	return r
}
