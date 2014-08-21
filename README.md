Natyla
======

Tarvis CI
[![Build Status](https://travis-ci.org/Fersca/natyla.svg?branch=master)](https://travis-ci.org/Fersca/natyla)
Drone CI
[![Build Status](https://drone.io/github.com/Fersca/natyla/status.png)](https://drone.io/github.com/Fersca/natyla/latest)
~~~
Current Code Coverage: 44% and increasing...
~~~

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

Caching
=======

If you keep the caching enabled (default) Natyla will use a 10MB (default) memory cache to store the most used Objects. If you reach the max defined amount of memory, Natyla will only cache the object metadata (but not the object content) it prevents for example invalid disk access for not previously cached DELETES. To disable cache, just add "cache":false in the config file.




