// Package conn contains a RTSP connection implementation.
package conn

import (
	"bufio"
	"encoding/base64"
	"io"

	"github.com/bluenviron/gortsplib/v4/pkg/base"
)

const (
	readBufferSize = 4096
)

// Conn is a RTSP connection.
type Conn struct {
	w  io.Writer
	br *bufio.Reader

	// reuse interleaved frames. they should never be passed to secondary routines
	fr base.InterleavedFrame
}

// NewConn allocates a Conn.
func NewConn(rw io.ReadWriter) *Conn {
	return &Conn{
		w:  rw,
		br: bufio.NewReaderSize(rw, readBufferSize),
	}
}

// Read reads a Request, a Response or an Interleaved frame.
func (c *Conn) Read() (interface{}, error) {
	byts, err := c.br.Peek(2)
	if err != nil {
		return nil, err
	}

	if byts[0] == base.InterleavedFrameMagicByte {
		return c.ReadInterleavedFrame()
	}

	// RT or HT
	if (byts[0] == 'R' || byts[0] == 'H') && byts[1] == 'T' {
		return c.ReadResponse()
	}

	return c.ReadRequest()
}

// ReadRequest reads a Request.
func (c *Conn) ReadRequest() (*base.Request, error) {
	var req base.Request
	err := req.Unmarshal(c.br)
	return &req, err
}

// ReadResponse reads a Response.
func (c *Conn) ReadResponse() (*base.Response, error) {
	var res base.Response
	err := res.Unmarshal(c.br)
	return &res, err
}

// ReadInterleavedFrame reads a InterleavedFrame.
func (c *Conn) ReadInterleavedFrame() (*base.InterleavedFrame, error) {
	err := c.fr.Unmarshal(c.br)
	return &c.fr, err
}

// WriteRequest writes a request.
func (c *Conn) WriteRequest(req *base.Request) error {
	buf, _ := req.Marshal()
	_, err := c.w.Write(buf)
	return err
}

// tunnel:
// WriteRequest writes a tunnel request.
func (c *Conn) WriteTunnelRequest(req *base.Request) error {
	// fmt.Printf(">>> TunnelRequest: %+v", req)

	buf, _ := req.Marshal()
	// When we're tunneling RTSP-over-HTTP, we Base-64-encode the request before we send it.
	// (However, we don't do this for the HTTP "GET" and "POST" commands that we use to set up the tunnel.)
	if req.Method != base.Get && req.Method != base.Post {
		encodedStr := base64.StdEncoding.EncodeToString(buf)
		bufout := []byte(encodedStr)

		_, err := c.w.Write(bufout)

		return err
	} else {
		_, err := c.w.Write(buf)
		return err
	}
}

// WriteResponse writes a response.
func (c *Conn) WriteResponse(res *base.Response) error {
	buf, _ := res.Marshal()
	_, err := c.w.Write(buf)
	return err
}

// WriteInterleavedFrame writes an interleaved frame.
func (c *Conn) WriteInterleavedFrame(fr *base.InterleavedFrame, buf []byte) error {
	n, _ := fr.MarshalTo(buf)
	_, err := c.w.Write(buf[:n])
	return err
}
