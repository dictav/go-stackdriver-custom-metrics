# go-stackdriver-custom-metrics

## Build

```
$ make -C deploy docker-build
```

## Run

**Please set GOOGLE_APPLICATION_CREDENTIALS before run**

```
$ make -eC deploy docker-run PROJECT=<YOUR_GCP_PROJECT>
```

## Deploy

**Please configure your docker using `gcloud auth configure-docker` before push**

```
$ make -C deploy push
```

```
$ make template PROJECT=<YOUR_GCP_PROJECT>
```

```
$ make -C deploy instance-group PROJECT=<YOUR_GCP_PROJECT>
```
