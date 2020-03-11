# Micro-database API
API for the microservice converting REST requests into MongoDB requests

## Home Link [/db]
Simple method to test if the Go API is running correctly  
Do not mix up with Status, which tests the status of the MongoDB daemon on pinky, 
this only confirms that at least one container with this API is running

### [GET]
+ Response 200 (text/plain)
    ~~~
    [MICRO-DATABASE] Home Link Joined
    ~~~

## Retrieving random snippets [/db/retrieve/snippets/{amount}]
This action searches the database for the amount of snippets specified, 
the snippets are selected randomly among the snippets that haven't been annotated or marked as unreadable
+ Parameters
    + amount (number) : Maximum number of snippets to be returned

### [GET]
This action has two negative responses defined :  
It will return a status 400 if an error occurs while converting {amount} to int.   
It will return a status 500 if an error occurs in the go service, this can happen either in the selection, 
the marshalling of the elements or the iteration over the selection results.

+ Response 200 (application/json)
    + Body
        ~~~
        [
            {
                "Id": "5e446bbf1cb5861676336ebd",
                "PiFF": {
                    "Meta": {
                        "Type": "line",
                        "URL": ""
                    },
                    "Location": [
                        {
                            "Type": "line",
                            "Polygon": [[0,0],[261,0],[261,343],[0,343]],
                            "Id": "loc_0"
                        }
                    ],
                    "Data": [
                        {
                            "Type": "line",
                            "LocationId": "loc_0",
                            "Value": "",
                            "Id": "0"
                        }
                    ],
                    "Children": null,
                    "Parent": 0
                },
                "Url": "/snippets/filename.png",
                "Annotated": false,
                "Corrected": false,
                "SentToReco": false,
                "Unreadable": false
            }
        ]
        ~~~

+ Response 400 (text/plain)  
    + Body
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~

+ Response 500 (text/plain)  
    + Body 
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~
      
## Retrieving all the database [/db/retrieve/all]
### [GET]
This action will return a status 500 if an error occurs in the Go service, 
this can happen either in the selection, the marshalling of the elements or the iteration over the selection results.

+ Response 200 (application/json)
    + Body
        ~~~TODO
        I'll write the example body when the microservices will be running again
        ~~~
+ Response 500 (text/plain) 
    + Body 
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~
      
## Database Status [/db/status]
### [GET]
Pings the MongoDB database and sends the result as a boolean.
+ Response 200 (application/json)
    + Body
        ~~~
        { 'isDBUp': true }
        ~~~
+ Response 200 (application/json)
    + Body
        ~~~
        { 'isDBUp': false }
        ~~~
      
## Create database entries [/db/insert]
### [POST]
+ Request (application/json)
    + Body
        ~~~TODO
        I'll write the example body when the microservices will be running again
        ~~~
       
+ Response 201 (text/plain)
    + Body
        ~~~TODO
        I'll write the example body when the microservices will be running again
        ~~~

+ Response 400 (text/plain)  
Error while reading body entry.
    + Body
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~

+ Response 500 (text/plain)  
Error in the Go service : Either while unmarshalling one of the entries or the database insertion.
    + Body 
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~
      
## Update flags [db/update/flags]
### [PUT]
+ Request (application/json)
    + Body
        ~~~TODO
        I'll write the example body when the microservices will be running again
        ~~~
       
+ Response 204

+ Response 400 (text/plain)  
Error while reading body entry.
    + Body
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~

+ Response 500 (text/plain)  
Error in the Go service : Either while unmarshalling one of the entries or the database update.
    + Body 
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~

## Add an annotation [db/update/value]
### [PUT]
+ Request (application/json)
    + Body
        ~~~TODO
        I'll write the example body when the microservices will be running again
        ~~~
       
+ Response 204

+ Response 400 (text/plain)  
Error while reading body entry.
    + Body
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~

+ Response 500 (text/plain)  
Error in the Go service : Either while unmarshalling one of the entries or the database update.
    + Body 
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~
      
## Empty database [db/delete/all]
### [PUT] (This endpoint will be removed in a future release)
+ Response 202 
**More error management will be added, standard response will be modified to 200**
### [DELETE]
An error can occur during the deletion of the database
+ Response 200
+ Response 500 (text/plain)  
    + Body 
        ~~~
        [MICRO-DATABASE] {Go error body}
        ~~~