package utl

func ref[T any](x T) *T { return &x }
