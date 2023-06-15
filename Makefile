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
	docker run --rm \
        --memory 2G \
        --log-opt max-size=5M \
        --log-opt max-file=3 \
        --name forum \
        -p 5432:5432 \
		-p 5000:5000 -d \
	  vk_edu/forum

