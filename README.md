# HTTP Auto Responder

Useful tool for static responses. It's like Fiddler's AutoResponder.

## Features

-   Rule definitions with matching methods.
-   Simple UI to manage Rules and Responses. You can navigate /http-auto-responder after the execution.
-   FARX & SAZ file import.
-   JSON files or boltdb for rules and responses storage.
-   Besides of FARX importing, you can use directly serve FARX files but no UI support for managing it. 
-   [Templating](https://golang.org/pkg/text/template/) on response using a basic model created from incoming request.

    ### FARX Files
        FARX files can be served directly. If you specify a FARX containing dir in the config file, tool reads all FARX files recursively and serve it accordingly. No UI support if you are using FARX files.
    #### FARX files reloding
        There is a reload enpoint available for if you specify FARX dir. Example;
        curl http://localhost/http-auto-responder/reload?path=farx_test.farx
        
## Installation
```
go get -u github.com/asalih/http-auto-responder
```
Then edit config.json

## Example
The tool provides minimum UI except if you want to use FARX files. After you get the source or get the release binaries then just edit the config file as you wish then you good to go!

Config file needs a storage configuration and a port for to listening http requests.
You can not use multiple storage option, please pick one according to your needs.
-   databaseName: boltDB database name.
-   jsonsFolderPath: rules and responses can be stored as json serialized if you specify the folder path.
-   farxFilesFolderPath: farx files can be served directly if you specify the folder path.

In this example below, FARX folder specified.
```
{
    "databaseName": "",
    "jsonsFolderPath": "",
    "farxFilesFolderPath": "./farx",
    "port": 80
}
```

## ToDo

-   How can we use SSL? or should we?