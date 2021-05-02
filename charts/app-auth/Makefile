CHART_REPO := http://jenkins-x-chartmuseum:8080
CURRENT=$(pwd)
NAME := app-auth
OS := $(shell uname)
RELEASE_VERSION := ""
TAG_VERSION := ""
SEMVER_VERSION := ""
CHARTMUSEUM_CREDS_USR := $(shell cat /tekton/home/basic-auth-user)
CHARTMUSEUM_CREDS_PSW := $(shell cat /tekton/home/basic-auth-pass)
DOCKER_REGISTRY := eu.gcr.io
MODIFIED := $(shell (test -s /workspace/source/scripts/outputs/app-auth || test -s /workspace/source/scripts/outputs/all) && echo true)


all:

ifeq ($(MODIFIED),true)
all: version ech docker-build tag release promote
else
all:
	exit 0
endif

version: install-semver
	$(eval CURRENT_VERSION := $(shell git ls-remote --tags --refs --sort="v:refname" | grep $(NAME) | sed -E 's/^[[:xdigit:]]+[[:space:]]+refs\/tags\/(.+)/\1/g' | tail -1))
	$(eval CURRENT_VERSION_NUMBER := $(shell awk -F\_ '{print $$2}' <<<$(CURRENT_VERSION)))
	$(eval SEMVER_VERSION := $(shell semver $(CURRENT_VERSION_NUMBER) -i))
	$(eval COMMIT_HEAD := $(shell git log --format="%H" -n 1 | cut -c1-7))
	$(eval TAG_VERSION := $(NAME)_v$(SEMVER_VERSION)_$(COMMIT_HEAD))
	$(eval RELEASE_VERSION := $(SEMVER_VERSION)_$(COMMIT_HEAD))
	git tag -a $(TAG_VERSION) -m '$(NAME) release $(RELEASE_VERSION)'
	git push origin $(TAG_VERSION)

install-semver:
	npm i semver -g
	git fetch origin

docker-build:
	gcloud auth activate-service-account clipfiler-ko@clipfiler.iam.gserviceaccount.com --key-file=/tekton/home/kaniko-secret.json --project=clipfiler
	gcloud builds submit --tag eu.gcr.io/clipfiler/$(NAME):$(RELEASE_VERSION) /workspace/source/services/app-auth

promote:
	jx promote -b --all-auto --timeout 1h --version $(SEMVER_VERSION)

ech:
	echo $(RELEASE_VERSION)
	echo $(SEMVER_VERSION)

build: clean
	rm -rf requirements.lock
	helm dependency build
	helm lint

install: clean build
	helm install . --name ${NAME}

upgrade: clean build
	helm upgrade ${NAME} .

delete:
	helm delete --purge ${NAME}

clean:
	rm -rf charts
	rm -rf ${NAME}*.tgz

release: clean
	helm init --stable-repo-url https://charts.helm.sh/stable
	helm dependency build
	helm lint
	helm init --client-only
	helm package .
	curl --fail -u $(CHARTMUSEUM_CREDS_USR):$(CHARTMUSEUM_CREDS_PSW) --data-binary "@$(NAME)-$(shell sed -n 's/^version: //p' Chart.yaml).tgz" $(CHART_REPO)/api/charts
	rm -rf ${NAME}*.tgz%

tag:
ifeq ($(OS),Darwin)
	sed -i "" -e "s/version:.*/version: $(SEMVER_VERSION)/" Chart.yaml
	sed -i "" -e "s/tag:.*/tag: $(RELEASE_VERSION)/" values.yaml
else ifeq ($(OS),Linux)
	sed -i -e "s/version:.*/version: $(SEMVER_VERSION)/" Chart.yaml
	sed -i -e "s|repository:.*|repository: $(DOCKER_REGISTRY)\/clipfiler\/$(NAME)|" values.yaml
	sed -i -e "s/tag:.*/tag: $(RELEASE_VERSION)/" values.yaml
else
	echo "platform $(OS) not supported to release from"
	exit -1

endif
	echo "All Good"
