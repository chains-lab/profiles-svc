package username

import (
	"crypto/rand"
	"fmt"

	"github.com/chains-lab/profiles-svc/internal/errx"
)

func GenerateUsername() (string, error) {
	const (
		prefix = "citizen"
		digits = 8
	)
	buf := make([]byte, digits)
	if _, err := rand.Read(buf); err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("cannot generate random digits: %w", err),
		)
	}
	for i := 0; i < digits; i++ {
		buf[i] = '0' + (buf[i] % 10)
	}
	return prefix + string(buf), nil
}
