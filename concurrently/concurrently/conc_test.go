package concurrently

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConc(t *testing.T) {
	errs := Concurrently(
		func() error { return errors.New("err1") },
		func() error { return errors.New("err2") })

	fmt.Println(errs)
	require.Len(t, errs, 3) //返回2个错误
}
