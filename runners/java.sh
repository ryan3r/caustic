#!/bin/sh
PATH="$PATH:/usr/lib/jvm/java-1.8-openjdk/bin"

case "$1" in
	compile)
		exec javac $2
		;;

	run)
		exec java "$(basename $2 .java)"
		;;
esac

