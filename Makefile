export GOOS=linux
export GOARCH=mipsle
export GOMIPS=softfloat

all: idig-station idig-api

idig-station:
	go build -o gl-mt300n-v2/files/usr/bin/idig-station -trimpath -ldflags="-s -w"

idig-api:
	go build -o gl-mt300n-v2/files/www/cgi-bin/idig-api -trimpath -ldflags="-s -w" idig-station/gl-mt300n-v2/idig-api

clean:
	rm -f gl-mt300n-v2/files/usr/bin/idig-station gl-mt300n-v2/files/www/cgi-bin/idig-api
