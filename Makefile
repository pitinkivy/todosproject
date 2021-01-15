test
	set ADDR=0.0.0.0:8889
	go env -w ADDR=0.0.0.0:8889
	$env: ADDR = "0.0.0.0:8889";
	set ADDR="0.0.0.0:8889" & go run main.go
	
	ADDR=0.0.0.0:8889 go run main.go
	APP_ADDR=0.0.0.0:8889 go run main.go

run: test
	go run main.go
	