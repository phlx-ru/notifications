package api

import (
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestMain(m *testing.M) {
	cleanup, err := bootstrap()
	if err != nil {
		log.Errorf(`bootstrap failed: %v`, err)
		os.Exit(1)
	}
	defer cleanup()
	os.Exit(m.Run())
}
