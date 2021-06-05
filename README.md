# text_to_audio

This is a simple web service that accepts text and sends back an audio file for users to listen
to that text.

## Purpose

It is meant to accompany a PWA app preferably built in django or something lighter like a frontend only
app built in vuejs or quasar to be enable users receive their audios on their phones if they download it.

## Design

- A single service that converts and sends back the audio file location.
  - The server returns the audio immediately as a download and deletes the file after (using a defer) [Or no need to save the file first. It could just be passed over to the user]
  - Authenticates user to know they are among the known users [This is to reduce the money I spend on gcloud]
- A front end application that allows:
  - a user to copy and paste code into it.
  - Click somehwere.
  - Wait for conversion to complete
  - return the downloaded link. I need to use go templates.

### Technology Stack

- The back end is a simple golang app with a web server
- The front end is a simple html file with a PWA enabled. I will need to use go templates.

## Notes

- No need for gofiber. I just need three templates.

  - One displays the textarea where text is to be put
  - The other displays a login with google button if user is not logged in
  - The other displays that sorry, you are not allowed here, if they are not in the list of email addresses allowed.

- For session, I can use google session to track the email of the signed in user.
