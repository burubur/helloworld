build:
	@docker build . -t burhanmubarok/helloworld:0.0.5

run:
	@docker run -d -p 8080:8080 --name helloworld burhanmubarok/helloworld:0.0.5

start:
	@docker start helloworld

stop:
	@docker stop helloworld

ping:
	@curl localhost:8080
