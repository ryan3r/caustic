while ! nc -z db 3306; do
	sleep 1
done

sleep 5
java -jar caustic.jar
