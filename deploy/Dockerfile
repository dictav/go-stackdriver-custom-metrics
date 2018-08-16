FROM alpine:3.8

RUN apk add --update --no-cache \
    curl \
  && rm -rf /var/cache/apk/*


ADD autoscale_test /autoscale_test
ADD entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
