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
        ```
        curl http://localhost/http-auto-responder/reload?path=farx_test.farx
        ```

## Installation

```
go get -u github.com/asalih/http-auto-responder
```
Then edit config.json

## ToDo

-   How can we use SSL? or should we?