# Twitch Telegram Announcer

## Problem

Manually posting in Telegram every time a Twitch stream goes live is inconvenient.
Ideally, this process should be automated and as inexpensive as possible.

## Solution

This repository provides **two serverless functions** required to implement the following workflow:

1. **Create a subscription (webhook)** for Twitch channel status updates.
   The webhook URL should point to the second serverless function.
2. **Receive a “stream started” event from Twitch** and automatically post a message
   in the specified Telegram channel.

Each function is configured via environment variables:

### Function: Create Subscription

| Environment Variable      | Description                                                                                   |
|---------------------------|-----------------------------------------------------------------------------------------------|
| `TW_WEBHOOK_CALLBACK_URL` | Webhook URL. This must point to the second serverless function.                               |
| `TW_CLIENT_ID`            | Twitch Client ID. Obtain it from the [Twitch Dev Console](https://dev.twitch.tv/console).     |
| `TW_CLIENT_SECRET`        | Twitch Client Secret. Obtain it from the [Twitch Dev Console](https://dev.twitch.tv/console). |
| `TW_CHANNEL_ID`           | Numeric ID of the Twitch channel to monitor (not the username).                               |
| `TW_WEBHOOK_SECRET`       | Secret string used to verify events received from Twitch. It can be any value.                |

### Function: Handle Twitch Events

| Environment Variable | Description                                                                                   |
|----------------------|-----------------------------------------------------------------------------------------------|
| `TG_CHAT_ID`         | Telegram chat/channel ID where the “stream started” message should be posted.                 |
| `TG_BOT_TOKEN`       | Telegram bot token for posting messages. **The bot must be present in the chat/channel.**     |
| `TW_CLIENT_ID`       | Twitch Client ID. Obtain it from the [Twitch Dev Console](https://dev.twitch.tv/console).     |
| `TW_CLIENT_SECRET`   | Twitch Client Secret. Obtain it from the [Twitch Dev Console](https://dev.twitch.tv/console). |
| `TW_WEBHOOK_SECRET`  | Secret string included with Twitch events, used to verify their authenticity.                 |
