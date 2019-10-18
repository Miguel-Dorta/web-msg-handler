# web-msg-handler
An API for handling messages from multiple website contact pages

## Objective
Unify multiple web contact forms backends in a single instance, with a simple and modular configuration.

## Supported protocols
* Email
* Telegram Bot

## Set up
#### Getting the software
You can either download a [built release](https://github.com/Miguel-Dorta/web-msg-handler/releases) or compile it yourself (it requires Go 1.13 or above for compilation).

### Linux manual installation (recommended)
* Extract the .tar.gz
* Move the binary (web-msg-handler) to its working directory ("/var/www/web-msg-handler" with owner and group "www-data" recommended).
* Create a config file. You can see an example in the extracted folder called "examples", or see the specification below.
* (Recommended) Don't expose it directly to the internet. Set up an nginx-reverse-proxy or similar.
* (Optional) Create a systemd service config or init config to start it up automatically.

### Linux automatic installation
IMPORTANT: the automatic installation will assume you use systemd and will create a working directory in "/var/www/web-msg-handler".
* Extract the .tar.gz
* Run install.sh
* Create a config file in /var/www/web-msg-handler/config.json. You can see an example in "/var/www/web-msg-handler/config.json.example", or see the specification below.

### Windows
* Extract the .zip
* Create a config.json file in the extracted directory (you can find an example in "examples/config.json" or see the specification below).
* Double click in web-msg-handler.exe to execute it.

### MacOS
* Download the software
* Create a config file (you can find an example in "examples/config.json" or see the specification below).
* Open a terminal
* Execute it (something like `Downloads/web-msg-handler_macOS/web-msg-handler --config <path>`)

#### Run parameters
```
--config <path>      sets the config.json path         Default: config.json
-h, --help           shows a help message and exits
--log-file <path>    sets a log file path
--port <port>        sets the port                     Default: 8080
--verbose <level>    sets verbose level (see below)    Default: 3
--version            prints version and exits
```

##### Verbose levels
* 0 = no log
* 1 = critical errors only
* 2 = errors and critical errors
* 3 = info, errors and critical errors
* 4 = debug, info, errors and critical errors

## Configuration file
The configuration file must be a JSON that contains a key called "sites" which is an array of site objects. The sites objects have:
* An ID (json key: "id"): a number from 0 to 18446744073709551615 (2^64 - 1) unique for every site object. It's recommended for them to be random.
* An URL (json key: "url"): the string of the URL of the form.
* A Google's reCAPTCHA v2 Secret (json key: "recaptchaSecret"): a string that contains the reCAPTCHA secret (more information [here](https://developers.google.com/recaptcha/intro)).
* A Sender object (json key: "sender"): this object contains the information of the type of sender for this site. It has:
    * A type (json key: "type"): a string that indicates which kind of sender is required ("mail" or "telegram").
    * A settings object (json key: "settings"): this object contains the settings unique for every type of sender. See below for more info.

### Mail sender settings
The settings object in the sender "mail" contains the following field:
* A send address (json key: "mailto"): the address where the contact information will be send.
* An username (json key: "username"): the username of the sender for logging.
* A password (json key: "password"): the password of the sender for logging.
* A SMTP hostname (json key: "hostname"): the address where the sender will log in for sending the mail.
* The port of the SMTP hostname (json key: "port"): must be a string!

### Telegram sender settings
The settings object in the sender "telegram" contains the following fields:
* A chat ID (json key: "chatID").
* A bot token (json key: "botToken").

More information about this [here](https://core.telegram.org/bots)

## Public API
The API of web-msg-handler tries to be minimal. It consists only in a request and a response.

### Request
The request must be made to the URL `/<ID>` where `<ID>` is the site ID of the config.json. This request must:
* Have a valid ID
* Be a POST request
* Have a header with key "Content-Type" and value "application/json"

The request must be a JSON that contains the following fields:
* "name"
* "mail"
* "msg"
* "g-recaptcha-response"

### Response
The response is a JSON that contains the following fields:
* "success": a boolean that indicates if the message was successfully send.
* "error" (only when success==false): a string that indicates why it failed.

## To do
* Add an example of a NGINX reverse proxy configuration.
* Implement a modular system of plugins for more types of senders.

## License
This software is licensed under MIT License. See [LICENSE](https://github.com/Miguel-Dorta/web-msg-handler/blob/master/LICENSE) for more information.
