/*
* A simple application allowing users to send text to an endpoint and receive back
* an audio they can play
 */
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	_ "github.com/joho/godotenv/autoload"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main()  {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal("error creating client: ", err)
	}
	defer client.Close()

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: "Hello, World"},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender: texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal("error synthesizing speech", err)
	}

	// FIXME: Instead of copying the file, do this https://stackoverflow.com/questions/24116147/how-to-download-file-in-browser-from-go-server#answer-24116517
	filename := "output.mp3"
	err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
	if err != nil {
		log.Fatal("Error writing audio file", err)
	}

	fmt.Printf("Audio content written to file: %v\n", filename)
	
}
