package tunnel

const (
	// same size as GStreamer's rtspsrc
	udpKernelReadBufferSize = 0x80000

	// 1500 (UDP MTU) - 20 (IP header) - 8 (UDP header)
	udpMaxPayloadSize = 1472

	// tunnel:
	RtspProtocol10 = "RTSP/1.0"
	HttpProtocol10 = "HTTP/1.0"
)
