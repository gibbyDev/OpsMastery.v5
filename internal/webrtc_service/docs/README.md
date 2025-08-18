# WebRTC Microservice Documentation

## Overview
The WebRTC microservice provides peer-to-peer video chat functionality by implementing a signaling server. It uses the `pion/webrtc` library for WebRTC protocols.

## Features
- Peer-to-peer connection setup
- Video and audio streaming
- Scalable architecture with signaling server

## Setup Instructions
1. Ensure Go is installed on your system.
2. Navigate to the project directory:
   ```bash
   cd /path/to/project
   ```
3. Install dependencies:
   ```bash
   go get github.com/pion/webrtc/v3
   ```
4. Run the tests to ensure everything is functioning:
   ```bash
   go test ./internal/webrtc_service/tests/...
   ```

## API Endpoints
### Start WebRTC Signaling Server
**Endpoint**: `/api/v1/webrtc/start`

**Method**: `POST`

**Description**: Starts the signaling server for WebRTC connections.

**Response**:
```json
{
  "status": "signaling server started"
}
```

## Integration Guidelines
- Use the `/api/v1/webrtc/start` endpoint to initiate the signaling server.
- Ensure the signaling server is running before initiating any WebRTC connections.

## Notes
- The signaling server is a critical component for WebRTC setup, handling the exchange of SDP and ICE candidates.
- For advanced configurations, refer to the [pion/webrtc documentation](https://github.com/pion/webrtc).

