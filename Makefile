.PHONY: build-image
build-image:
	DOCKER_BUILDKIT=1 docker build -t vk_edu/forum .

.PHONY: stop
stop:
	docker stop forum

.PHONY: logs
logs:
	docker logs forum


.PHONY: run
run: build-image
	docker run --rm -d \
        --memory 2G \
        --log-opt max-size=5M \
        --log-opt max-file=3 \
        --name forum \
		-p 5000:5000 \
	  vk_edu/forum

