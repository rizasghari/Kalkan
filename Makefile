run-demo-origins:
	docker run --rm -d -p 9007:80 --name origin1 kennethreitz/httpbin
	docker run --rm -d -p 9008:80 --name origin2 kennethreitz/httpbin
	docker run --rm -d -p 9009:80 --name origin3 kennethreitz/httpbin

stop-demo-origins:
	docker stop origin1
	docker stop origin2
	docker stop origin3