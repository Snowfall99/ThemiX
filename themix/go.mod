module go.themix.io/themix

go 1.16

require (
	github.com/andelf/go-curl v0.0.0-20200630032108-fd49ff24ed97 // indirect
	go.themix.io/crypto v0.0.0-00010101000000-000000000000
	go.themix.io/transport v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.16.0
)

replace go.themix.io/transport => ../transport

replace go.themix.io/crypto => ../crypto
