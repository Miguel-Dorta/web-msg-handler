const https = require("https");

// printError prints a string to stderr
function printError(err: string) {
    process.stderr.write(err);
}

// Settings is the object that contains the site settings.
class Settings {
    webName: string;
    chatID: string;
    botToken: string;

    constructor(json: string) {
        let sett = JSON.parse(json)
        if (!sett.hasOwnProperty("website_name") || !sett.hasOwnProperty("chat_id") || !sett.hasOwnProperty("bot_token")) {
            throw new Error("invalid object");
        }
        this.webName = sett.website_name;
        this.chatID = sett.chat_id;
        this.botToken = sett.bot_token;
    }
}

// Message is the interface that contains the message itself.
// The object provided will always implement this interface.
interface Message {
    name: string;
    mail: string;
    msg: string;
}

// escapeHTML escapes reserved characters in HTML
function escapeHTML(s: string) : string {
    return s.replace("&", "&amp;")
            .replace("'", "&#39;")
            .replace("<", "&lt;")
            .replace(">", "&gt;")
            .replace("\"", "&#34;");
}

// composeMsg creates the string message that will be sent
function composeMsg(sett: Settings, msg: Message) : string {
    let webName = escapeHTML(sett.webName);
    let name = escapeHTML(msg.name);
    let mail = escapeHTML(msg.mail);
    let escapedMsg = escapeHTML(msg.msg);
    return `Message from ${webName}\n\n<b>Name:</b> ${name}\n<b>Email:</b> ${mail}\n<b>Message:</b> ${escapedMsg}`;
}

// send is the main function of the script.
// It will send the msg provided to the sender specified in sett.
function send(sett: Settings, msg: Message) {
    // compose data to send
    let data = JSON.stringify({
        "chat_id": sett.chatID,
        "text": composeMsg(sett, msg),
        "parse_mode": "HTML",
        "disable_web_page_preview": true,
    });

    // make POST request
    let req = https.request({
        hostname: "api.telegram.org",
        path: `/bot${sett.botToken}/sendMessage`,
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
            'Content-Length': data.length
        },
    }, (resp) => {
        resp.on('data', (d) => {
            let respObj = JSON.parse(d)
            if (!respObj.ok) {
                printError(`server returned error: ${respObj.error_message}`)
            }
        })
    });

    req.on('error', function(e) {
        printError(`request failed: ${e.message}`)
    });

    req.write(data)
    req.end()
}

send(new Settings(process.argv[2]), JSON.parse(process.argv[3]));
