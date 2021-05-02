OS := $(shell uname)
DOCKER_REGISTRY := eu.gcr.io
NAME := app-auth
OS := $(shell uname)
PREVIEW_VERSION := ""
TAG_VERSION := ""
SEMVER_VERSION := ""
MODIFIED := $(shell (test -s /workspace/source/scripts/outputs/app-auth || test -s /workspace/source/scripts/outputs/all) && echo true)
PR_NUMBER := $(shell cat /workspace/source/PR-VERSION)


all: version ech update-cors docker-build tag-preview

version: install-semver
	$(eval CURRENT_VERSION := $(shell git ls-remote --tags --refs --sort="v:refname" | grep $(NAME) | sed -E 's/^[[:xdigit:]]+[[:space:]]+refs\/tags\/(.+)/\1/g' | tail -1))
	$(eval CURRENT_VERSION_NUMBER := $(shell awk -F\_ '{print $$2}' <<<$(CURRENT_VERSION)))
	$(eval SEMVER_VERSION := $(shell semver $(CURRENT_VERSION_NUMBER) -i prerelease --preid preview))
	$(eval COMMIT_HEAD := $(shell cat /tekton/home/MERGE_HEAD_SHA | cut -c1-7))
	$(eval TAG_VERSION := $(NAME)_v$(SEMVER_VERSION)_$(COMMIT_HEAD))
	$(eval PREVIEW_VERSION := $(SEMVER_VERSION)_$(COMMIT_HEAD))
	git tag -a $(TAG_VERSION) -m '$(NAME) release $(PREVIEW_VERSION)'
	git push origin $(TAG_VERSION)

static-version:
	$(eval CURRENT_VERSION := $(shell git ls-remote --tags --refs --sort="v:refname" | grep $(NAME) | sed -E 's/^[[:xdigit:]]+[[:space:]]+refs\/tags\/(.+)/\1/g' | tail -1))
	$(eval CURRENT_VERSION_NUMBER := $(shell awk -F\_ '{print $$2}' <<<$(CURRENT_VERSION)))
	$(eval SEMVER_VERSION := $(shell echo $(CURRENT_VERSION_NUMBER) | cut -c2-))
	$(eval COMMIT_HEAD := $(shell awk -F\_ '{print $$3}' <<<$(CURRENT_VERSION)))
	$(eval TAG_VERSION := $(NAME)_v$(SEMVER_VERSION)_$(COMMIT_HEAD))
	$(eval PREVIEW_VERSION := $(SEMVER_VERSION)_$(COMMIT_HEAD))

install-semver:
	npm i semver -g
	git fetch origin

ech:
	echo $(PREVIEW_VERSION)
	echo $(SEMVER_VERSION)

docker-build:
	gcloud auth activate-service-account clipfiler-ko@clipfiler.iam.gserviceaccount.com --key-file=/tekton/home/kaniko-secret.json --project=clipfiler
	gcloud builds submit --tag eu.gcr.io/clipfiler/$(NAME):$(PREVIEW_VERSION) /workspace/source/services/app-auth

preview:
	jx preview --app $(NAME) --version $(SEMVER_VERSION) --dir ../..

update-cors:
ifeq ($(OS),Darwin)
	sed -i "" -e "s/{{NUMBER}}/$(PR_NUMBER)/" ../../main.go
else ifeq ($(OS),Linux)
	sed -i -e "s/{{NUMBER}}/$(PR_NUMBER)/" ../../main.go
else
	echo "platfrom $(OS) not supported to release from"
	exit -1
endif
	echo "CORS UPDATED"

tag-preview:
ifeq ($(OS),Darwin)
	sed -i "" -e "s/version:.*/version: $(SEMVER_VERSION)/" Chart.yaml
	sed -i "" -e "s/version:.*/version: $(SEMVER_VERSION)/" ../*/Chart.yaml
	sed -i "" -e "s/tag:.*/tag: $(PREVIEW_VERSION)/" values.yaml
	sed -i "" -e "s/tag:.*/tag: $(PREVIEW_VERSION)/" ../*/values.yaml
	sed -i "" -e "s/$(NAME)-version:.*/version: $(SEMVER_VERSION)/" ../../../../test/bdd/charts/preview/requirements.yaml
else ifeq ($(OS),Linux)
	sed -i -e "s/version:.*/version: $(SEMVER_VERSION)/" Chart.yaml
	sed -i -e "s/version:.*/version: $(SEMVER_VERSION)/" ../*/Chart.yaml
	sed -i -e "s|repository:.*|repository: $(DOCKER_REGISTRY)\/clipfiler\/$(NAME)|" values.yaml
	sed -i -e "s/tag:.*/tag: $(PREVIEW_VERSION)/" values.yaml
	sed -i -e "s/tag:.*/tag: $(PREVIEW_VERSION)/" ../*/values.yaml
	sed -i -e "s/$(NAME)-version:.*/version: $(SEMVER_VERSION)/" ../../../../test/bdd/charts/preview/requirements.yaml
else
	echo "platfrom $(OS) not supported to release from"
	exit -1
endif
	echo "  version: $(SEMVER_VERSION)" >> requirements.yaml
	helm init --stable-repo-url https://charts.helm.sh/stable
	jx step helm build
