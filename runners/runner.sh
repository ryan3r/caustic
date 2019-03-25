POLLING_RATE=5

# Send and error to the user
error() {
	echo $1 >&2
}

# Compile and run the code
run() {
	local fileName="$1"
	local lang="java"

	case "${lang##*.}" in
		java)
			lang="java"
			;;
		cpp|cxx|cc)
			lang="cpp"
			;;
	esac

	./$lang.sh compile $fileName || error "Compile error"
	./$lang.sh run $fileName || error "Runtime error"
}

# while true; do
# 	# TODO: Check for new submissions
# 	sleep $POLLING_RATE
# done
