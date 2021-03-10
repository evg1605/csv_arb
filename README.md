# csv-arb converter
Convertor from csv to arb and from arb to csv

Install from source:
```
dir=$(mktemp -d) 
git clone https://github.com/evg1605/csv_arb "$dir" 
cd "$dir"
go install -ldflags "-s -w" -v  ./cmd/arbc.go
```
