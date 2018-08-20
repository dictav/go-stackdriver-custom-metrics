BASE_NAME=autoscale-test
DOCKER_IMAGE=$(BASE_NAME):0.0.12
METRIC=custom.googleapis.com/autoscaling/count

GOOGLE_APPLICATION_CREDENTIALS=
PROJECT=
ZONE=asia-northeast1-a
MACHINE=f1-micro
LOGCMD=gcloud --project=$(PROJECT) beta logging read 'resource.type="global" AND "autoscale-test"' --format=json --freshness=60s | jq -r '.[] | .receiveTimestamp + " " + .jsonPayload.data'
AUTOSCALE_UTILIZATION_TARGET=19

all: docker-build push template instance-group

push:
	docker tag $(DOCKER_IMAGE) gcr.io/$(PROJECT)/$(DOCKER_IMAGE)
	docker push gcr.io/$(PROJECT)/$(DOCKER_IMAGE)

docker-run:
	docker run -it --rm \
		-v $(GOOGLE_APPLICATION_CREDENTIALS):/cred.json \
		-e GOOGLE_APPLICATION_CREDENTIALS=/cred.json \
		autoscale_test:0.0.1 -project=$(PROJECT) test-instance

docker-build: Dockerfile autoscale_test entrypoint.sh
	docker build -t $(DOCKER_IMAGE) .
	rm autoscale_test

autoscale_test:
	GOOS=linux GOARCH=amd64 go build ../cmd/autoscale_test
	upx autoscale_test

template:
	gcloud --project=dictav-net beta compute instance-templates \
		create-with-container $$(echo $(DOCKER_IMAGE) | tr :. -) \
		--preemptible \
		--machine-type=$(MACHINE) \
		--scopes=https://www.googleapis.com/auth/cloud-platform \
		--container-image=gcr.io/$(PROJECT)/$(DOCKER_IMAGE)

instance-group: instance-group-create instance-group-autoscale

instance-group-create:
	gcloud beta compute --project=$(PROJECT) instance-groups managed \
		create $(BASE_NAME) \
		--base-instance-name=$(BASE_NAME) \
		--template=$$(echo $(DOCKER_IMAGE) | tr :. -) \
		--size=1 \
		--zone=$(ZONE) \
		--initial-delay=300

instance-group-autoscale:
	gcloud compute --project=$(PROJECT) instance-groups managed \
		set-autoscaling $(BASE_NAME) \
		--zone=$(ZONE) \
		--cool-down-period=60 \
		--max-num-replicas=5 \
		--min-num-replicas=1 \
		--custom-metric-utilization "metric=$(METRIC),utilization-target-type=GAUGE,utilization-target=$(AUTOSCALE_UTILIZATION_TARGET)"

instance-group-delete:
	gcloud beta compute --project=$(PROJECT) instance-groups managed \
		delete $(BASE_NAME) \
		--zone=$(ZONE)

logs:
	while :; do $(LOGCMD); sleep 60; done
