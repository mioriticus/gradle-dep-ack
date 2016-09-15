#!/bin/sh

rm -f gradle-dep-list.txt
./gradlew jwl:dependencies --configuration compile > gradle-dep-list.txt
