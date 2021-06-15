set datafile separator ","

set term pngcairo size 1024,768
set output 'graph.png'

set xlabel "iteration"
set ylabel "Max. RSS KB"
set xtics 1
set ytics 10240

set style line 1 lc rgb '#8b1a0e' pt 1 ps 1 lt 1 lw 2 # --- red
set style line 2 lc rgb '#5e9c36' pt 6 ps 1 lt 1 lw 2 # --- green

set style line 12 lc rgb '#808080' lt 0 lw 1
set grid back ls 12
plot 'logs/with-reader.log' using 1:2 with linespoints linestyle 1 title "reader", \
     'logs/with-read-seeker.log' using 1:2 with linespoints linestyle 2 title "read-seeker"
