#!/bin/sh

if [ $# = 0 ]; then
  project=$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/project/project-id)
  echo project=$project
  zone=$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone | sed -e 's,.*/,,')
  echo zone=$zone
  instance=$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/id)
  echo instance=$instance

  opts="-project=$project -zone=$zone $instance"
  echo opts=$opts
else
  opts=$@
fi

echo opts=$#:$opts

/autoscale_test $opts
