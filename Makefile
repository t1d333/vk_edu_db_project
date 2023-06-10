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
        --name forum \
        -p 5432:5432 \
		-p 5000:5000 -d \
	  vk_edu/forum

