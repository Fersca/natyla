#!/bin/bash
c=1
while [ $c -le 19000 ]
do
	echo "Welcone $c times"
	curl -H "Content-Type: application/json" -X POST -d '{"name":"fer","surname":"scasserra","age":31, "location":"buenos aires de america", "scuadra":"river plata varias veces campeon del mundo","pais":"argentina, un gran pais para vivir"}' http://localhost:8080/fer/$1A$c
	(( c++ ))
done
echo ""
echo "Fin cache test"
