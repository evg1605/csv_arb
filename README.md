# csv-arb converter

Convertor from csv to arb and from arb to csv

#### Install from releases:
Download last release from https://github.com/evg1605/csv_arb/releases
#### Install from sources:

```
dir=$(mktemp -d) 
git clone https://github.com/evg1605/csv_arb "$dir" 
cd "$dir"
go install -ldflags "-s -w -X main.AppVersion=dev"  ./arbc
```

#### How to use:

```
arbc csv2arb --csv-path=[PATH_OR_URL_TO_CSV_FILE] --arb-path=[PATH_TO_FOLDER_CONTAINS_ARB_FILES]
arbc arb2csv --arb-path=[PATH_TO_FOLDER_CONTAINS_ARB_FILES] --csv-path=[PATH_TO_CSV_FILE]
```

#### Full params list for csv2arb command:
```
   --arb-path                    arb folder path (folder contains arb files - one for every culture)
   --csv-path                    url or path of csv file 
   --arb-template                arb file template (default: app_{culture}.arb)
   --col-descr                   name column name in csv table (default: description)
   --col-name                    name column name in csv table (default: name)
   --col-params                  name column name in csv table (default: parameters)
   --culture                     default culture (default: en)
   --help                        displays usage information of the application or a command (default: false)
   --log-level                   log level (trace, debug, info, warning, error, fatal, panic) (default: error)
```
<br/>

#### Full params list for arb2csv command:
```
   --arb-path                    arb folder path (folder contains arb files - one for every culture) 
   --csv-path                    path to csv file
   --col-descr                   name column name in csv table (default: description)
   --col-name                    name column name in csv table (default: name)
   --col-params                  name column name in csv table (default: parameters)
   --culture                     default culture (default: en)
   --help                        displays usage information of the application or a command (default: false)
   --log-level                   log level (trace, debug, info, warning, error, fatal, panic) (default: error)
```
#### Example csv table

| name               	| description                   	| parameters 	| en                                       	| ru                             	|
|--------------------	|-------------------------------	|------------	|------------------------------------------	|--------------------------------	|
| mainTitle          	| Title of application          	|            	| My super application                     	| Моё супер приложение           	|
| passwordFieldError 	| Validation error for password 	| min;max    	| Password len must be from {min} to {max} 	| Длина пароля от {min} до {max} 	|

