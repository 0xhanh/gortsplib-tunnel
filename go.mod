module gortsplib-tunnel

go 1.19

require (
	github.com/bluenviron/gortsplib/v4 v4.3.0
	github.com/pion/rtcp v1.2.10
	github.com/pion/rtp v1.8.2
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/bluenviron/mediacommon v1.5.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/sdp/v3 v3.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// replace github.com/bluenviron/gortsplib/v4 v4.3.0 => github.com/0xhanh/gortsplib v0.0.0-20231101081823-e5ff82d994fc

replace github.com/bluenviron/gortsplib/v4 v4.3.0 => /home/hanh/elcom/vms/bluenviron/gortsplib
