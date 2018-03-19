FROM=$(shell cat $(DOCKERFILE) | grep FROM | cut -d' ' -f2)
DOCKERFILE_PATH=$(PWD)/image/$(OS)

SLASH := /
DASH := -
DOT := .
COLON := :

PREBUILT := N

# These values are changed in each version branch
# This is the only place they need to be changed
# other than the README.md file.
include versions.mk

IMG_STRING=$(shell echo $(IMAGE_NAME) | cut -d'/' -f2 | sed -e 's/$(OS)/nearform/g;')
RH_TARGET=registry.rhc4tp.openshift.com:443/$(RH_PID)/$(IMG_STRING):$(IMAGE_TAG)
TARGET=$(IMAGE_NAME):$(IMAGE_TAG)
ARCHIVE_NAME=$(IMAGE_NAME)-$(IMAGE_TAG)
ARCHIVE=sources-$(subst $(SLASH),$(DASH),$(ARCHIVE_NAME)).tgz

envinfo:
	@echo $(call .FEATURES)
	@env
.PHONY: all
all: build squash test

.PHONY: build
build:
	PREBUILT=$(PREBUILT) OS=$(OS) ./contrib/etc/get_node_source.sh "${NODE_VERSION}" $(PWD)/src/ && \
	docker build -f $(DOCKERFILE_PATH)/Dockerfile \
	--build-arg NODE_VERSION=$(NODE_VERSION) \
	--build-arg NPM_VERSION=$(NPM_VERSION) \
	--build-arg V8_VERSION=$(V8_VERSION) \
	--build-arg PREBUILT=$(PREBUILT) \
	-t $(TARGET) .

.PHONY: squash
squash:
	docker-squash -f $(FROM) $(TARGET) -t $(TARGET)
	docker run $(TARGET) ls -Alh /usr/libexec/s2i

.PHONY: test
test:
	 BUILDER=$(TARGET) NODE_VERSION=$(NODE_VERSION) ./test/run.sh

.PHONY: clean
clean:
	docker rmi `docker images $(TARGET) -q`

.PHONY: publish
publish:
	@echo $(DOCKER_PASS) | docker login --username $(DOCKER_USER) --password-stdin
	docker push $(TARGET)
ifndef DEBUG_BUILD
ifdef LATEST
	docker tag $(TARGET) $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):latest
endif
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
endif


.PHONY: redhat_publish
redhat_publish:
	echo "Publishing to RedHat repository"
ifndef DEBUG_BUILD
	docker tag nearform/$(OS)-s2i-nodejs:$(TAG) $(RH_TARGET)
	docker push $(RH_TARGET)
endif

.PHONY: archive
archive:
	mkdir -p dist
	git archive --prefix=build-tools/ --format=tar HEAD | gzip >dist/build-tools.tgz
	cp -v versions.mk dist/versions.mk
	git rev-parse HEAD >dist/build-tools.revision
	cd src/node-v$(NODE_VERSION) &&  git archive --format=tar v$(NODE_VERSION) | gzip >../../dist/node-v$(NODE_VERSION).tgz && cd ../..
	shasum dist/* >checksum
	cp -v checksum dist/dist.checksum
	tar czvf $(ARCHIVE) dist/*


.PHONY: upload
upload:
	echo "Attempting Upload of sources to S3 bucket $(S3BUCKET)"
	s3cmd put $(ARCHIVE) "$(S3BUCKET)"
