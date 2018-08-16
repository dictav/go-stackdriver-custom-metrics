if [ $# = 0 ]; then
  v=empty
else
  v="-arg=$1 -arg=2"
fi

echo $v
