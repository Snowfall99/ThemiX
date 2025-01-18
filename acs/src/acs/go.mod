module go.themix.io/acs

go 1.23.0

require (
	go.themix.io/crypto v0.0.0-00010101000000-000000000000
	go.themix.io/transport v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.16.0
)

require (
	go.dedis.ch/fixbuf v1.0.3 // indirect
	go.dedis.ch/kyber/v3 v3.0.13 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529 // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
)

replace go.themix.io/transport => ../transport

replace go.themix.io/crypto => ../crypto
