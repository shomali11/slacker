# slacker [![Slack](https://img.shields.io/badge/slack-%23slacker--framework-orange)](https://gophers.slack.com/archives/C051MGM3GFL) [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/slacker)](https://goreportcard.com/report/github.com/shomali11/slacker) [![GoDoc](https://godoc.org/github.com/shomali11/slacker?status.svg)](https://godoc.org/github.com/shomali11/slacker) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) 

Built on top of the Slack API [github.com/slack-go/slack](https://github.com/slack-go/slack), Slacker is a low-friction framework for creating Slack Bots.

## Features

- Supports Slack Apps using [Socket Mode](https://api.slack.com/apis/connections/socket)
- Easy definitions of commands and their input
- Built-in `help` command
- Bot responds to mentions and direct messages
- Simple parsing of String, Integer, Float and Boolean parameters
- Customizable, intuitive and with many examples to follow
- Replies can be new messages or in threads
- Replies can be ephemeral, scheduled, updated or deleted
- Supports Slash Commands and Interactive Messages
- Supports `context.Context`
- Supports middlewares & grouping of commands
- Supports Cron Jobs using [https://github.com/robfig/cron](https://github.com/robfig/cron)
- Handlers run concurrently via goroutines
- Full access to the Slack API [github.com/slack-go/slack](https://github.com/slack-go/slack)

# Install

```
go get github.com/shomali11/slacker/v2
```

# Examples

We wrote extensive [examples](./examples) to help you familiarize yourself with Slacker!

# Preparing your Slack App

To use Slacker you'll need to create a Slack App, either [manually](#manual-steps) or with an [app manifest](#app-manifest). The app manifest feature is easier, but is a beta feature from Slack and thus may break/change without much notice.

## Manual Steps

Slacker works by communicating with the Slack [Events API](https://api.slack.com/apis/connections/events-api) using the [Socket Mode](https://api.slack.com/apis/connections/socket) connection protocol.

To get started, you must have or create a [Slack App](https://api.slack.com/apps?new_app=1) and enable `Socket Mode`, which will generate your app token (`SLACK_APP_TOKEN` in the examples) that will be needed to authenticate.

Additionally, you need to subscribe to events for your bot to respond to under the `Event Subscriptions` section. Common event subscriptions for bots include `app_mention` or `message.im`.

After setting up your subscriptions, add scopes necessary to your bot in the `OAuth & Permissions`. The following scopes are recommended for getting started, though you may need to add/remove scopes depending on your bots purpose:

* `app_mentions:read`
* `channels:history`
* `chat:write`
* `groups:history`
* `im:history`
* `mpim:history`

Once you've selected your scopes install your app to the workspace and navigate back to the `OAuth & Permissions` section. Here you can retrieve yor bot's OAuth token (`SLACK_BOT_TOKEN` in the examples) from the top of the page.

With both tokens in hand, you can now proceed with the examples below.

## App Manifest

Slack [App Manifests](https://api.slack.com/reference/manifests) make it easy to share a app configurations. We provide a [simple manifest](./app_manifest/manifest.yml) that should work with all the examples provided below.

The manifest provided will send all messages in channels your bot is in to the bot (including DMs) and not just ones that actually mention them in the message.

If you wish to only have your bot respond to messages they are directly messaged in, you will need to add the `app_mentions:read` scope, and remove:

- `im:history`       # single-person dm
- `mpim:history`     # multi-person dm
- `channels:history` # public channels
- `groups:history`   # private channels

You'll also need to adjust the event subscriptions, adding `app_mention` and removing:

- `message.channels`
- `message.groups`
- `message.im`
- `message.mpim`

# Contributing / Submitting an Issue

Please review our [Contribution Guidelines](CONTRIBUTING.md) if you have found
an issue with Slacker or wish to contribute to the project.

# Troubleshooting

## My bot is not responding to events

There are a few common issues that can cause this:

* The OAuth (bot) Token may be incorrect. In this case authentication does not fail like it does if the App Token is incorrect, and the bot will simply have no scopes and be unable to respond.
* Required scopes are missing from the OAuth (bot) Token. Similar to the incorrect OAuth Token, without the necessary scopes, the bot cannot respond.
* The bot does not have the correct event subscriptions setup, and is not receiving events to respond to.
