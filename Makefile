build:
	docker build --rm --no-cache --progress=plain -t burhanmubarok/helloworld:0.5.3 .

run:
	docker run --rm -d -p 8080:8080 --name helloworld burhanmubarok/helloworld:0.5.3

shell:
	docker exec -it helloworld /bin/sh

start:
	docker start helloworld

stop:
	docker stop helloworld

ping:
	curl localhost:8080
