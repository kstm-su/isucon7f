GOPATH := ${PWD}
export GOPATH

.PHONY: deploy
deploy:
	ssh root@app0252 systemctl stop cco.golang.service
	ssh root@app0253 systemctl stop cco.golang.service
	go build -v app
	scp app root@app0252:/home/isucon/webapp/go/app
	scp app root@app0253:/home/isucon/webapp/go/app
	ssh root@app0252 systemctl start cco.golang.service
	ssh root@app0253 systemctl start cco.golang.service
	sudo systemctl restart cco.golang.service

build:
		go build -v app

ensure:
		cd src/app && dep ensure -vendor-only

update:
		cd src/app && dep ensure

test:
		go test -v app

vet:
		go vet ./src/app/...
