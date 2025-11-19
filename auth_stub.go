//go:build !windows

package vcenterauth

import (
	"context"
	"errors"
)

// ErrSSPINotSupported is returned when LoginSSPI is called on non-Windows platforms
var ErrSSPINotSupported = errors.New("SSPI is only supported on Windows")

// LoginSSPI is not available on non-Windows platforms
func LoginSSPI(ctx context.Context, host string, insecure bool) (*Client, error) {
	return nil, ErrSSPINotSupported
}
