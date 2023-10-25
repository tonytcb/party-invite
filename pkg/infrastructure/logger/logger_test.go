package logger

import (
	"bytes"
	"context"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestSimpleLogger_Infof(t *testing.T) {
	t.Parallel()

	var (
		buff = &bytes.Buffer{}
		log  = NewLogger(buff)
		cid  = "1-123-4"
		ctx  = context.WithValue(context.Background(), config.CorrelationIDKeyName, cid)
	)

	log.Infof("Test Info method")
	assert.Contains(t, readBuffer(t, buff), "[INFO] [cid=] Test Info method")

	log = log.FromContext(ctx)
	log.Errorf("Test Error method")
	assert.Contains(t, readBuffer(t, buff), "[ERROR] [cid=1-123-4] Test Error method")

	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, readBuffer(t, buff), "[FATAL] [cid=1-123-4] Test Fatal method")
			return
		}

		t.Error("recover method was not called")
	}()

	log.Fatalf("Test Fatal method")
}

func readBuffer(t *testing.T, buff *bytes.Buffer) string {
	b, err := io.ReadAll(buff)
	if err != nil {
		t.Fatalf("failed to read buffer: %v", err)
	}

	return string(b)
}
