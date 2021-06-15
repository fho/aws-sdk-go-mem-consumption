#!/usr/bin/env bash

dir="logs"
parallel_uploads=10

rm -r "$dir"
mkdir "$dir"
for (( i=1; i< 15;i++)); do
  echo -n "$i," >> "$dir/with-read-seeker.log"
  /usr/bin/time -a -o "$dir/with-read-seeker.log" -f '%M' bin/v2 -parallel-uploads $parallel_uploads -repeat $i -with-read-seeker

  echo -n "$i," >> "$dir/with-reader.log"
  /usr/bin/time -a -o "$dir/with-reader.log" -f '%M' bin/v2 -parallel-uploads $parallel_uploads -repeat $i
done

gnuplot graph.gnuplot
mv graph.png "graph-$i-parallel-uploads.png"
