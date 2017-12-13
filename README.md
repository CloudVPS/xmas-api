# xmas-api
This is the code used to expose the xmas-api to control the APC units in the CloudVPS office. 

# Features

* Multiple APC location support.
* Power on, off or flip the current state.

# Install

Adjust the apc.go file with your APC unit and compile with go build.

# Usage

Run the server and use the API on http://\<ip>:8000/v1/apc/

## Switch on:

 GET http://\<ip>\:8000/v1/apc/?loc=cloud&port=4&state=ON

## Switch off:

 GET http://\<ip\>:8000/v1/apc/?loc=cloud&port=4&state=OFF

## Switch flip:

 GET http://\<ip>\:8000/v1/apc/?loc=cloud&port=4&state=FLIP

# Todo

* Read APC configuration from file.
* Use snmp library instead of local snmp commands.
* TLS


