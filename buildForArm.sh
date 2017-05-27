#!/bin/bash

# This scripts builds and prepares goScience software
# for distributing on RPI devices (arm processors)

DIR="builds/armBuild"

echo ""
echo "**********************************************"
echo "* Start building goScience for arm processor *"
echo "**********************************************"

# build for arm computer main.go
env GOOS=linux GOARCH=arm go build main.go
echo "Code compilation completed"

# check if directory $DIR exists otherwise create it
if [ -d "$DIR" ]; then
    echo "Removing old files from $DIR"
    rm -r "$DIR"
    mkdir "$DIR"
else
    echo "$DIR does not exist, creating new $DIR directory"
    mkdir "$DIR"
fi

# move main.go to armBuild folder
mv main "$DIR"

# copy dependency folders (templates, js and css) to $DIR
echo "Copying goScience dependencies to $DIR directory"
cp -r public "$DIR"
cp -r templates "$DIR"

echo "Build complete"
