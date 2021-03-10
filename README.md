# csv-arb converter
Convertor from csv to arb and from arb to csv

#### Install from source:
```
dir=$(mktemp -d) 
git clone https://github.com/evg1605/csv_arb "$dir" 
cd "$dir"
go install -ldflags "-s -w" -v  ./cmd/arbc.go
```

#### How to use:
```
arbc --mode=csv2arb --csv-path=[PATH_TO_CSV_FILE] --arb-path=[PATH_TO_FOLDER_CONTAINS_ARB_FILES]
```
#### Full params list:
`--mode` conversion mode, possible values:<br/>
&nbsp;&nbsp;&nbsp;&nbsp;`csv2arb` from csv to arb<br/>
&nbsp;&nbsp;&nbsp;&nbsp;`arb2csv` from arb to csv<br/>

`--csv-path` path to csv file