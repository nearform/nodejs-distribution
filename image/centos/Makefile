FROM=openshift/base-centos7


# These values are changed in each version branch
# This is the only place they need to be changed
# other than the README.md file.
include versions.mk

TARGET=$(IMAGE_NAME):$(IMAGE_TAG)
SCRIPT_PATH=$(PWD)
REPO_ROOT=$(PWD)/../..

.PHONY: all
all: build squash test

.PHONY: build
build:
	cd $(REPO_ROOT) \
	PREBUILT=$(PREBUILT) ./contrib/etc/get_node_source.sh "${NODE_VERSION}" $(REPO_ROOT)/src/ && \
	docker build -f $(SCRIPT_PATH)/Dockerfile \
	--build-arg NODE_VERSION=$(NODE_VERSION) \
	--build-arg NPM_VERSION=$(NPM_VERSION) \
	--build-arg V8_VERSION=$(V8_VERSION) \
	-t $(TARGET) .

.PHONY: squash
squash:
	docker-squash -f $(FROM) $(TARGET) -t $(TARGET)

.PHONY: test
test:
	 BUILDER=$(TARGET) NODE_VERSION=$(NODE_VERSION) ./test/run.sh

.PHONY: clean
clean:
	docker rmi `docker images $(TARGET) -q`

.PHONY: tag
tag:
	if [ ! -z $(LTS_TAG) ]; then docker tag $(TARGET) $(IMAGE_NAME):$(LTS_TAG); fi

.PHONY: publish
publish:
	echo $(DOCKER_PASS) | docker login --username $(DOCKER_USER) --password-stdin
	docker push $(TARGET)
ifdef MAJOR_TAG
	docker tag $(TARGET) $(IMAGE_NAME):$(MAJOR_TAG)
	docker push $(IMAGE_NAME):$(MAJOR_TAG)
endif
ifdef MINOR_TAG
	docker tag $(TARGET) $(IMAGE_NAME):$(MINOR_TAG)
	docker push $(IMAGE_NAME):$(MINOR_TAG)
endif
ifdef LTS_TAG
	docker tag $(TARGET) $(IMAGE_NAME):$(LTS_TAG)
	docker push $(IMAGE_NAME):$(LTS_TAG)
endif
