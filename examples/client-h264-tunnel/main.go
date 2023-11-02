package main

import (
	"fmt"
	tunnel "gortsplib-tunnel"
	"log"

	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/bluenviron/gortsplib/v4/pkg/format/rtph264"
	"github.com/bluenviron/gortsplib/v4/pkg/url"
	"github.com/pion/rtp"
)

// This example shows how to
// 1. connect to a RTSP server
// 2. check if there's an H264 media stream
// 3. decode the H264 media stream into RGBA frames

// This example requires the FFmpeg libraries, that can be installed with this command:
// apt install -y libavformat-dev libswscale-dev gcc pkg-config

func main() {

	// parse URL
	rtsp_url := "rtsp://localhost:8554/mystream?tunnelPort=8112"

	u, err := url.Parse(rtsp_url)
	if err != nil {
		panic(err)
	}

	c := tunnel.NewClientTunnel(8112)
	// c.SetTunnelOverPort(8112)

	// connect to the server
	err = c.Start(u.Scheme, u.Host)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	resp, err := c.GET(u)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GET Resp: %+v", resp)

	_, err = c.POST(u)
	if err != nil {
		panic(err)
	}

	// find published medias
	desc, _, err := c.Describe(u)
	if err != nil {
		panic(err)
	}

	// find the H264 media and format
	var forma *format.H264
	medi := desc.FindFormat(&forma)
	if medi == nil {
		panic("media not found")
	}

	// setup RTP/H264 -> H264 decoder
	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		panic(err)
	}

	// setup H264 -> raw frames decoder
	frameDec, err := newH264Decoder()
	if err != nil {
		panic(err)
	}
	defer frameDec.close()

	// if SPS and PPS are present into the SDP, send them to the decoder
	if forma.SPS != nil {
		frameDec.decode(forma.SPS)
	}
	if forma.PPS != nil {
		frameDec.decode(forma.PPS)
	}

	// setup a single media
	_, err = c.Setup(desc.BaseURL, medi, 0, 0)
	if err != nil {
		panic(err)
	}

	// called when a RTP packet arrives
	c.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// decode timestamp
		pts, ok := c.PacketPTS(medi, pkt)
		if !ok {
			log.Printf("waiting for timestamp")
			return
		}

		// extract access units from RTP packets
		au, err := rtpDec.Decode(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				log.Printf("ERR: %v", err)
			}
			return
		}

		for _, nalu := range au {
			// convert NALUs into RGBA frames
			img, err := frameDec.decode(nalu)
			if err != nil {
				panic(err)
			}

			// wait for a frame
			if img == nil {
				continue
			}

			log.Printf("decoded frame with PTS %v and size %v", pts, img.Bounds().Max)
		}
	})

	// start playing
	_, err = c.Play(nil)
	if err != nil {
		panic(err)
	}

	// wait until a fatal error
	panic(c.Wait())
}
