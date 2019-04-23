while ! nc -z db 3306; do
	sleep 1
done

java -jar caustic.jar
