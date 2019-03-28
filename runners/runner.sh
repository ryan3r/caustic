POLLING_RATE=5

# Send and error to the user
error() {
	echo $1 >&2
}

timedKiller() {
	sleep 3
	kill -9 $1
}

# Run a program for a limited amount of time
limitTime() {
	timedKiller $1 &
	wait $1
	status=$?
	kill $!
	return $status
}

# Compile and run the code
run() {
	local fileName="$1"
	local lang="java"

	case "${fileName##*.}" in
		java)
			lang="java"
			;;
		cpp|cxx|cc)
			lang="cpp"
			;;
	esac

	/app/$lang.sh compile $fileName || {
		error "Compile error"
		return
	}

	/app/$lang.sh run $fileName &
	limitTime $! || error "Runtime error"
}

cd /mnt/submissions
while true; do
	[ -f "/mnt/submissions/foo.java" ] && run /mnt/submissions/foo.java
	sleep $POLLING_RATE
done
