BASE_NAME=autoscale-test
DOCKER_IMAGE=$(BASE_NAME):0.0.2
METRIC=custom.googleapis.com/autoscaling/count

GOOGLE_APPLICATION_CREDENTIALS=
PROJECT=
ZONE=asia-northeast1-a

template_name:
	echo $(TEMPLATE_NAME)
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
		--custom-metric-utilization "metric=$(METRIC),utilization-target-type=GAUGE,utilization-target=5"

