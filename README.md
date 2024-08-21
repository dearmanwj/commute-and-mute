# Commute and Mute

Monorepo for Commute and Mute backend. This projects includes all Lambda functions with the infrastructure as code required to deploy the functions to AWS.

The backend uses the Strava Webhooks API to detect when a user's activity is between their home and work and if so to set it as a commute activity type and mute it to hide it from their feed.

## Architecture

This project users a serverless architecture to handle events from the Strava API and to manage user configuration. This is split up into the following parts:

### Onboard

The onboard funcion responds to a user's request to grant permission to this application to receive events based on their activity. The users access tokens are stored in the database with the athelete's ID

### Activity

This function subscribes to the activities emitted by Strava, filters only the new activity events and retrieves the new activity ID. A request is then made to the Strava REST API to get the activity and then matches the start and finish locations to the user's home and work locations for cycling activities. If the activity passes all the checks, a PUT request is made to update the acitivty with the commute and mute options.

This function also handles OAuth2 token refresh if needed.

### Users API Gateway

...
