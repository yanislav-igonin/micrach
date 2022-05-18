  echo "GET http://micrach.igonin.dev" | vegeta attack -rate 100 -duration=120s | tee results.bin | vegeta report
  # vegeta report -type=json results.bin > metrics.json
  cat results.bin | vegeta plot > plot.html
  open plot.html
  # cat results.bin | vegeta report -type="hist[0,100ms,200ms,300ms]"
