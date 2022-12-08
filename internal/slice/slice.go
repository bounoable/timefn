package slice

// Map maps a slice using the given function.
func Map[In, Out any](s []In, fn func(In) Out) []Out {
	if s == nil {
		return nil
	}

	out := make([]Out, len(s))
	for i, v := range s {
		out[i] = fn(v)
	}

	return out
}
