# web-msg-handler
An API for handling messages from multiple website contact pages

## Objective
Unify multiple web contact forms backends in a single instance, with a simple and modular configuration.

## Senders available by default
* Email
* Telegram Bot

## Dependencies
* Node.js
* \[Optional\] nodemailer (Node.js package): required for sending emails.
#### Additional build dependencies
* Bash
* Go
* GoReleaser
* TypeScript

## Set up
#### Getting the software
You can either download a [built release](https://github.com/Miguel-Dorta/web-msg-handler/releases) or build it yourself.

### Installation
* Extract the .tar.gz (recommended to extract it in a new directory).
* Execute install.sh (working directory must be where the .tar.gz contents were extracted).

## Public API
The API of web-msg-handler tries to be minimal. It consists only in a request and a response.

### Request
The request must be made to the URL `/<ID>` where `<ID>` is the site ID of the config.json. This request must:
* Have a valid ID
* Be a POST request
* Have a header with key "Content-Type" and value that contains "application/json"

The request must be a JSON that contains the following fields:
* "name"
* "mail"
* "msg"
* "g-recaptcha-response"

### Response
The response is a JSON that contains the following fields:
* "success": a boolean that indicates if the message was successfully send.
* "error" (only when success==false): a string that indicates why it failed.

## License
This software is licensed under MIT License. See [LICENSE](https://github.com/Miguel-Dorta/web-msg-handler/blob/master/LICENSE) for more information.
