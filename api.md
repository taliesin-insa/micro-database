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

## Retrieving snippets with annotation suggestions [/db/retrieve/snippets/{amount}]
This action searches the database for the amount of snippets specified,
The snippets are first selected randomly among the snippets that have been annotated by the recognizer
If not enough annotated snippets are found, the application will complete with unannotated snippets 
+ Parameters
    + amount (number) : Number of snippets desired

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
                            "Value": "The effets of marijuana on the human body",
                            "Id": "0"
                        }
                    ],
                    "Children": null,
                    "Parent": 0
                },
                "Url": "/snippets/filename.png",
                "Annotated": true,
                "Corrected": false,
                "SentToReco": true,
                "Unreadable": false,
                "Annotator": "$taliesin_recognizer"
            },
            {
                "Id": "5e446bbf1cb586167cef156",
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
                "Url": "/snippets/filename2.png",
                "Annotated": false,
                "Corrected": false,
                "SentToReco": false,
                "Unreadable": false,
                "Annotator": ""
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
      
## Retrieving snippets to be sent to the recognizer [/db/retrieve/recognizer/{amount}]
This action searches the database for the amount of snippets specified, 
the snippets are selected randomly among the snippets that haven't been annotated and haven't been sent to the recognizer yet.
The fact that the snippets have been sent to the recognizer will then be recorded in the database
+ Parameters
  + amount (number) : Number of snippets desired

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
              "Unreadable": false,
              "Annotator": ""
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
        {"isDBUp":true,"total":4,"annotated":0,"unreadable":0}
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
        ~~~
      [
          {
              "PiFF": {
                  "Meta": {
                      "Type": "line",
                      "URL": ""
                  },
                  "Location": [
                      {
                          "Type": "line",
                          "Polygon": [
                              [
                                  0,
                                  0
                              ],
                              [
                                  1080,
                                  0
                              ],
                              [
                                  1080,
                                  1920
                              ],
                              [
                                  0,
                                  1920
                              ]
                          ],
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
              "Url": "/snippets/wp1951689.jpg"
          },
          {
              "PiFF": {
                  "Meta": {
                      "Type": "line",
                      "URL": ""
                  },
                  "Location": [
                      {
                          "Type": "line",
                          "Polygon": [
                              [
                                  0,
                                  0
                              ],
                              [
                                  667,
                                  0
                              ],
                              [
                                  667,
                                  1000
                              ],
                              [
                                  0,
                                  1000
                              ]
                          ],
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
              "Url": "/snippets/photo-1444703686981-a3abbc4d4fe3.jpg"
          },
          {
              "PiFF": {
                  "Meta": {
                      "Type": "line",
                      "URL": ""
                  },
                  "Location": [
                      {
                          "Type": "line",
                          "Polygon": [
                              [
                                  0,
                                  0
                              ],
                              [
                                  1080,
                                  0
                              ],
                              [
                                  1080,
                                  1920
                              ],
                              [
                                  0,
                                  1920
                              ]
                          ],
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
              "Url": "/snippets/The-art-of-sales-through-story-telling.jpg"
          },
          {
              "PiFF": {
                  "Meta": {
                      "Type": "line",
                      "URL": ""
                  },
                  "Location": [
                      {
                          "Type": "line",
                          "Polygon": [
                              [
                                  0,
                                  0
                              ],
                              [
                                  1080,
                                  0
                              ],
                              [
                                  1080,
                                  1920
                              ],
                              [
                                  0,
                                  1920
                              ]
                          ],
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
              "Url": "/snippets/wallpaper-1698879.jpg"
          }
      ]
        ~~~
       
+ Response 201 (text/plain)
    + Body
        ~~~
        ["5e81db20c096cc792fff5094","5e81db20c096cc792fff5095","5e81db20c096cc792fff5096","5e81db20c096cc792fff5097"]
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

## Add an annotation specifying the annotator [db/update/value/{annotator}]
+ Parameters
    + annotator (string) : Annotator Identifier 
### [PUT]
+ Request (application/json)
    + Body
        ~~~
        [{"Id":"5e679a2c005e59a282790a76","Value":"Premiere annotation"},{"Id":"5e679a2c005e59a282790a98","Value":"Deuxième annotation"}]
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