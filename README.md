# Redash Slackbot

Simple slackbot for redash

TODO: images

## Setting up redash-slackbot

1. [Create a bot in Slack](https://slack.com/apps/new/A0F7YS25R)

2. Create a channel called "redash-images". This is used to succinctly store and display visualizations to users.

## Deploy with Docker

docker run -e REDASH_API_KEY -e REDASH_URL -e SLACK_TOKEN justmiles/redash-slackbot

## Run locally

Ensure you have Chrome installed. Redash's API, unfortunatly, doesn't expose a way to extract the visualization image. Instead, we're using chromedp to capture a screenshot of the embedded URL for a visualization. Hacky, but it works.

export REDASH_API_KEY=
export REDASH_URL=
export SLACK_TOKEN=

redash-slackbot
