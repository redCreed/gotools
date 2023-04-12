package concurrently

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroup(t *testing.T) {
	group := Group{}
	group.Go(func() error {
		return fmt.Errorf("first")
	})
	group.Go(func() error {
		return nil
	})
	group.Go(func() error {
		return fmt.Errorf("second")
	})
	group.Go(func() error {
		return fmt.Errorf("third")
	})

	allErrors := group.Wait()
	fmt.Println(allErrors)
	require.Len(t, allErrors, 3)
}
