## run: starts demo http services
run-containers:
	docker run --rm -d -p 9007:80 --name server1 kennethreitz/httpbin
	docker run --rm -d -p 9008:80 --name server2 kennethreitz/httpbin
	docker run --rm -d -p 9009:80 --name server3 kennethreitz/httpbin

## stop: stops all demo services
stop:
	docker stop server1
	docker stop server2
	docker stop server3