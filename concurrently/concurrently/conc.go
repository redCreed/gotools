// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information

package concurrently

// Concurrently runs fns concurrently and returns the non-nil errors.
func Concurrently(fns ...func() error) []error {
	var g Group
	for _, fn := range fns {
		g.Go(fn)
	}
	return g.Wait()
}
