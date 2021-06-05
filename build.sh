#!/bin/bash
echo -ne "Building the artifacts...\n"
env GOOS=linux GOARCH=amd64 go build -o ./dist/soptextaudio sopherapps.com/text_to_audio
echo "Done...\n"