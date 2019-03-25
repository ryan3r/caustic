#!/bin/sh
case "$1" in
	compile)
		exec javac $2
		;;

	run)
		exec java "$(basename $2 .java)"
		;;
esac

