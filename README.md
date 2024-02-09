# One Time Service - Home Assistant Add-on

_Call a service with a one time token from outside Home Assistant!_

"One Time Service" lets you create publically exposed tokens for one time usage, that will execute a Home Assistant service without any authentication. This is great, e.g. for visitors who don't have an account on your Home Assistant instance.

**Possible use cases:**
- Allow a visitor to unlock the door themselves if no one else is at home
- Allow neighbors to notifiy you of emergencies by turning a smart RGB light red
- Let your neighbors activate your vacuum cleaner when you are on vacation

## Installation

[![Open your Home Assistant instance and show the add add-on repository dialog with a specific repository URL pre-filled.](https://my.home-assistant.io/badges/supervisor_add_addon_repository.svg)](https://my.home-assistant.io/redirect/supervisor_add_addon_repository/?repository_url=https%3A%2F%2Fgit.leon.wtf%2Fleon%2Fleon.wtf-home-assistant-addons)

This button will add the add-on repository at `https://git.leon.wtf/leon/leon.wtf-home-assistant-addons` to your Home Assistant. You can then install and update One Time Service via the Add-on Store.

## How does it work?

The add-on exposes a separate port (`1337` by default) which can be used by a reverse proxy exposed to the internet, while Home Assistant itself can reside behind your firewall (see examples below). This endpoint is used to submit the tokens from the "outside world". The add-on's admin dashboard is used to define service calls and to generate tokens. You can add the dashboard to your Home Assistant sidebar from the add-on settings page.

A **service call** (previously called "macros") is a possibly complex Home Assistant service call used for better reusability inside the add-on, e.g. if used in multiple tokens at once or if a token must be recreated on a regular basis.

A **token** is configured to execute a previously created service call and can optionally have a comment visible on the submission page and/or an expiry time.

The **differences to just using vanilla webhook automations** are:
- A webhook URL can neither expire nor be invalidated once used
- The accidential use of a webhook URL can occur. E.g. when you call the webhook URL from a mobile browser, leave the page open, and reopen the browser later, the webhook is called again.

Tokens an service call definitions are stored in a JSON file in `/share/one-time-service/db.json`.

#### Define a new service call

![Service call creation](https://git.leon.wtf/leon/one-time-service/-/raw/main/screenshots/macro-creation.png)

#### Create a token from a service call

![Token creation](https://git.leon.wtf/leon/one-time-service/-/raw/main/screenshots/token-creation.png)

#### Public token submission page

![Public token submission](https://git.leon.wtf/leon/one-time-service/-/raw/main/screenshots/token-submission.png)

## Add-on options

| Option                       | Description |
|------------------------------|-------------|
| `public_token_base_url`      | This is the URL (or IP) exposed by the reverse proxy, that proxies requests to the submission port (default `1337`). Tokens are generated for this URL like this: `<base-url>/?token=<token>`. |
| `notify_on_token_submission` | If enabled token submissions are logged to the notification target. |
| `notification_target`        | The Home Assistant service to send notifications with. E.g. use `notify.notify` to send push notifications to all pushers. |  

## Reverse proxy examples

In these examples `10.0.30.1` is the IP of Home Assistant inside our network and `https://smart-token.example.com` is set as `public_token_base_url`. Tokens are then generated like this: `https://smart-token.example.com/?token=<token>`.

**Caddy 2**

```Caddyfile
https://smart-token.example.org {
    reverse_proxy * 10.0.30.1:1337
}
```

## Tech stack

The backend is written in Go 1.20 using the [Gin Web Framework](https://github.com/gin-gonic).

The public frontend for token submission is written in plain HTML and JavaScript.

The internal ingress frontend for service call and token management is written in [React](https://react.dev/) (via CRA).

Both frontends use [Pico CSS](https://picocss.com/) as their CSS framework.
