Natyla
======

Tarvis CI
[![Build Status](https://travis-ci.org/Fersca/natyla.svg?branch=master)](https://travis-ci.org/Fersca/natyla)
Drone CI
[![Build Status](https://drone.io/github.com/Fersca/natyla/status.png)](https://drone.io/github.com/Fersca/natyla/latest)
~~~
Current Code Coverage: 92.1%
~~~

![(https://github.com/Fersca/natyla/blob/master/logoNatyla.png?raw=true)
]
Natyla is a Full-Stack REST-API/Cache/Key-Value-Store application to configure and run simple APIs in minutes. Written in Golang, it provides the same functionality as a multithreaded application running with Memcached and MongoDB.

Install and Run
===============

Just clone this repository, install [golang](http://golang.org/) and run:

~~~
go get
go run natyla.go
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

Custom Config
=============

You can create (or download the example) a config file called config.json, where you can setup a custom configuration.

~~~
{
  "token":"customToken",
  "cache":false,
  "memory":10485760,
  "data_dir":"myDir"
}
~~~


Using Natyla RESTful API
========================

Natyla provides a RESTful API to read, update and store JSON resources.

To create a resource (a Person), just POST or PUT the JSON object to the specific resource:
~~~
curl -X POST localhost:8080/Person -d '{"id":123456,"name":"Ferdinand", "age":32,"profession":"engineer"}'
~~~

**You Always have to provide an "id" field**

If you want to read a resource, just call the API with the resource ID:

~~~
curl localhost:8080/Person/123456
~~~

You will get the Stored JSON:

~~~
{"id":123456,"name":"Ferdinand", "age":32,"profession":"engineer"}
~~~

If you are curious, you will notice that Natyla stored the JSON resource under you "data" directory.
The previous example will save the JSON (in plain text) in the following file: 

~~~
data/Person/123456.json
~~~

To delete an Object, just delete it indicating the Object ID:

~~~
curl -X DELETE localhost:8080/Person/123456
~~~

**TODO:** 
Multiget: In the near future you will be able to request several recources at the same time. Eg:

~~~
curl localhost:8080/Person?ids=123456,789101
~~~
You will receive:
~~~
[{"id":123456,"name":"Ferdinand", "age":32,"profession":"engineer"},{"id":789101,"name":"Norbert", "age":57,"profession":"engineer"}]
~~~


Searching
=========

If you want to search for a particular value in a resource field, you should use the "search" feature from Natyla.

E.g: If you want to search for all the "engineers" in the "Person" resource, just call:

~~~
curl localhost:8080/Person/search?field=profession&equal=engineer
~~~

And you will get an array of resources that satisfy the query

~~~
[{"id":123456,"name":"Ferdinand", "age":32,"profession":"engineer"}]
~~~

**TODO:**
In the near future you will be able to do "like", "or", "greater than", etc operations on several fields

Caching
=======

If you keep the caching enabled (default) Natyla will use a 10MB (default) memory cache to store the most used Objects. If you reach the max defined amount of memory, Natyla will only cache the object metadata (but not the object content) it prevents for example invalid disk access for not previously cached DELETES. To disable cache, just add "cache":false in the config file.

Notifications
=============

TODO: 
Natyla is not providing a notification system by now, but it will. 
The idea is the following: If you want to be notified for each resource change, you will be able to register your callback URL and Natyla will notify you with a JSON POST. 
The second notification system will able you to listen to a socket stream and receive the same notifications while you are connected to the channel.

Telnet Administration
=====================

You can manage Natyla by doing a telnet to port 8081. 
Just type "help" and you will see how to check the memory usage, interact with the resources, etc.
