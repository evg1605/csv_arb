# csv-arb converter

Convertor from csv to arb and from arb to csv

#### Install from source:

```
dir=$(mktemp -d) 
git clone https://github.com/evg1605/csv_arb "$dir" 
cd "$dir"
go install -ldflags "-s -w"  ./cmd/arbc.go
```

#### How to use:

```
arbc --mode=csv2arb --csv-path=[PATH_TO_CSV_FILE] --arb-path=[PATH_TO_FOLDER_CONTAINS_ARB_FILES]
```

#### Full params list:

`--mode`- conversion mode, possible values:<br/>
&nbsp;&nbsp;&nbsp;&nbsp;`csv2arb`- from csv to arb<br/>
&nbsp;&nbsp;&nbsp;&nbsp;`arb2csv`- from arb to csv<br/>

`--csv-path`- path to csv file<br/><br/>
`--csv-url`- url to csv file<br/><br/>
`--arb-path`- arb folder path (folder contains arb files - one for every culture)<br/><br/>
`--arb-template`- template of arb file name, default value is **app_{culture}.arb**<br/><br/>
`--default-culture`- default culture, default value is **en**<br/><br/>
`--log`- log level (**trace**, **debug**, **info**, **warning**, **error**, **fatal**, **panic**), default value is **error**<br/><br/>
#### Example csv table

| name               	| description                   	| parameters 	| en                                       	| ru                             	|
|--------------------	|-------------------------------	|------------	|------------------------------------------	|--------------------------------	|
| mainTitle          	| Title of application          	|            	| My super application                     	| Моё супер приложение           	|
| passwordFieldError 	| Validation error for password 	| min;max    	| Password len must be from {min} to {max} 	| Длина пароля от {min} до {max} 	|

