go run server/cmd/main.go --id 0 --port 11200 --keys ../crypto --cluster http://127.0.0.1:11200,http://127.0.0.1:11210,http://127.0.0.1:11220
go run server/cmd/main.go --id 1 --port 11210 --keys ../crypto --cluster http://127.0.0.1:11200,http://127.0.0.1:11210,http://127.0.0.1:11220
go run server/cmd/main.go --id 2 --port 11220 --keys ../crypto --cluster http://127.0.0.1:11200,http://127.0.0.1:11210,http://127.0.0.1:11220
# go run server/cmd/main.go --id 3 --port 11230 --keys ../crypto --cluster http://127.0.0.1:11200,http://127.0.0.1:11210,http://127.0.0.1:11220,http://127.0.0.1:11230
