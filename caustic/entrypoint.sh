while ! nc -z $MYSQL_HOST 3306; do
	sleep 1
done

java -jar caustic.jar
