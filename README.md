<div style="text-align:center">
  <img src ="https://github.com/Fersca/natyla/blob/master/logoNatyla.png?raw=true" />
</div>

~~~
Current Code Coverage: 85.9% 
~~~

# Natyla - Startup's Best Friend :)

[![Build Status](https://travis-ci.org/Fersca/natyla.svg?branch=master)](https://travis-ci.org/Fersca/natyla)

Natyla is a Full-Stack REST-API/Cache/Key-Value-Store application to configure and run simple APIs in minutes. Written in Golang, it provides the same functionality as a multithreaded application running with Memcached and MongoDB.

### Need to create an API in 3 minutes?

  - Download natyla.
  - Start it.
  - POST your information in JSON format.
  - Share the resource URL.
  - Done. Be Happy :)

## Install and Run

Just clone this repository, install [golang](http://golang.org/) and run:

~~~
$ go get github.com/Fersca/natyla 
$ cd $GOPATH/src/github.com/Fersca/natyla
$ go run natyla.go
~~~

It will start Natyla with the default configuration:
  - 10MB memory cache.
  - "adminToken" as default token for admin access to PUT/POST/DELETE.
  - $CURRENT_DIR/data as default directory for storing resources and JSON Objects. 

~~~
Starting Natyla...
Core numbers:  4
Can't found 'config.json' using default parameters
Using Config: map[token:adminToken memory:10485760]
Max memory defined as:  10  Mbytes
Ready.
-------------------------------
~~~

## Custom Config

You can create (or use the example) a config file called config.json, where you can setup a custom configuration.
~~~
{
    "token":"customToken",
    "cache":false,
    "memory":10485760,
    "data_dir":"myDir",
    "api_port":"8080",
    "telnet_port":"8081",
    "print_log":true
}
~~~

## Using Natyla RESTful API

Natyla provides a RESTful API to read, update and store JSON resources.

### Create a resource (POST, PUT)

To create a resource (a Person), just POST or PUT the JSON object to the specific resource:

~~~
curl -X POST localhost:8080/person -d '{"id":123456,"name":"Ferdinand", "age":32,"profession":"engineer"}'
~~~

**You Always have to provide an "id" field**

### Read a resource (GET)

If you want to read a resource, just call the API with the resource ID:

~~~
curl localhost:8080/person/123456
~~~

You will get the Stored JSON:

~~~
{"id":123456,"name":"Ferdinand", "age":32,"profession":"engineer"}
~~~

If you are curious, you will notice that Natyla stored the JSON resource under you "data" directory.
The previous example will save the JSON (in plain text) in the following file: 

~~~
data/person/123456.json
~~~

### Delete a resource (DELETE)

To delete a resource, just DELETE it indicating the Object ID:

~~~
curl -X DELETE localhost:8080/person/123456
~~~

### Multiget **(Future)**

In the near future you will be able to request several resources at the same time. Eg:

~~~
curl localhost:8080/person?ids=123456,789101
~~~
You will receive:
~~~
[
  {"id":123456, "name":"Ferdinand", "age":32, "profession":"engineer"},
  {"id":789101, "name":"Norbert",   "age":57, "profession":"engineer"}
]
~~~

## Searching

If you want to search for a particular value in a resource field, you should use the "search" feature from Natyla.

E.g. Assume you have the following "Person" resources:

~~~
[
  {"id":123456, "name":"Ferdinand", "age":32, "profession":"engineer"},
  {"id":123456, "name":"Andrea",    "age":26, "profession":"designer"},
  {"id":789012, "name":"Argandas",  "age":25, "profession":"engineer"}
]
~~~
Wich are located on:
~~~
localhost:8080/person
~~~
And you want to search for all the "engineers" in the "Person" resource, just call:
~~~
curl localhost:8080/person/search?field=profession&equal=engineer
~~~
or
~~~
curl localhost:8080/person?profession=engineer
~~~
And you will get an array of resources that satisfies the query:
~~~
[
  {"id":123456,"name":"Ferdinand", "age":32,"profession":"engineer"},
  {"id":789012,"name":"Argandas",  "age":25,"profession":"engineer"}
]
~~~
Even you can search multiple fields to refine your search, doing as follow:
~~~
curl localhost:8080/person?profession=engineer&age=32
~~~
Will return an array of resources that satisfies the query:
~~~
[
  {"id":123456,"name":"Ferdinand", "age":32,"profession":"engineer"}
]
~~~

**TODO:**
In the near future you will be able to do "like", "or", "greater than", etc. operations on several fields

## Caching

If you keep the caching enabled (default) Natyla will use a 10MB (default) memory cache to store the most used recources. If you reach the max defined amount of memory, Natyla will only cache the resource metadata (but not the resource content) it prevents (for example) invalid disk access for not previously cached DELETES. To disable cache, just add "cache":false in the config file.

## Formatting

Natyla provides a pretty printing format if you GET any resource from a browser. It allows you to interact with your API resources in a friendly way :)

## Notifications

**TODO:**
Natyla is not providing a notification system by now, but it will. The idea is the following: If you want to be notified when a resource changes, you will be able to register your callback URL and Natyla will notify you with a JSON POST. The second notification system will able you to listen to a socket stream and receive the same notifications while you are connected to the channel.

## Telnet Administration

You can manage Natyla by doing a telnet to port 8081. Just type "help" and you will see how to check the memory usage, interact with the resources, etc.

## Want to Help?

![alt tag](https://github.com/Fersca/natyla/blob/gh-pages/images/go.png?raw=true)

Natyla API needs the help of the world best Golang Developers to improve its functionality!!! if you are interested just send a twitter message or start adding code to the repo! just try to follow the project backlog :) Thanks!!

[Backlog](https://github.com/Fersca/natyla/blob/master/Backlog.txt)
