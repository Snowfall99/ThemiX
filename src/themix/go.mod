module go.themix.io/themix

go 1.23.0

require (
	go.themix.io/client v0.0.0-00010101000000-000000000000
	go.themix.io/crypto v0.0.0-00010101000000-000000000000
	go.themix.io/transport v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.16.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/oasislabs/ed25519 v0.0.0-20200302143042-29f6767a7c3e // indirect
	github.com/perlin-network/noise v1.1.3 // indirect
	go.dedis.ch/fixbuf v1.0.3 // indirect
	go.dedis.ch/kyber/v3 v3.0.13 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace go.themix.io/transport => ../transport

replace go.themix.io/crypto => ../crypto

replace go.themix.io/client => ../client
