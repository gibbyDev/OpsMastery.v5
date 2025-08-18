
import (
	"testing"

	"OpsMastery.v5/internal/webrtc_service"
)

func TestSignalingServerStart(t *testing.T) {
	t.Run("Should start signaling server without error", func(t *testing.T) {
		// Capture output or errors from server start
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Signaling server crashed: %v", r)
			}
		}()

		// Call the StartSignalingServer function
		webrtc_service.StartSignalingServer()
	})
}

