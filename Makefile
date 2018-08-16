
watch:
	go run cmd/sdcustom/*.go -debug -project=$(PROJECT) -metric=$(METRIC)

run:
	go run cmd/autoscale-test/*.go -project=$(PROJECT) 12345678
