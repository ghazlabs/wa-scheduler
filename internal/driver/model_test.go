package driver

import (
	"net/http"
	"testing"

	"github.com/ghazlabs/wa-scheduler/internal/core"
)

func TestNewErrorRespMapsRetryNonFailedToBadRequest(t *testing.T) {
	resp := NewErrorResp(core.ErrRetryNonFailed)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
	if resp.Err != "ERR_BAD_REQUEST" {
		t.Fatalf("Err = %q, want %q", resp.Err, "ERR_BAD_REQUEST")
	}
	if resp.Message != retryNonFailedMessage {
		t.Fatalf("Message = %q, want %q", resp.Message, retryNonFailedMessage)
	}
}
