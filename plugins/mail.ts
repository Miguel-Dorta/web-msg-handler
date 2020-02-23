const nodemailer = require("nodemailer");

// Settings is the object that contains the site settings.
class Settings {
    webName: string;
    mailto: string;
    username: string;
    password: string;
    hostname: string;
    port: number;

    constructor(json: string) {
        let sett = JSON.parse(json);
        if (
            !sett.hasOwnProperty("website_name") ||
            !sett.hasOwnProperty("mailto") ||
            !sett.hasOwnProperty("username") ||
            !sett.hasOwnProperty("password") ||
            !sett.hasOwnProperty("hostname") ||
            !sett.hasOwnProperty("port")
        ) {
            throw new Error("invalid object")
        }
        this.webName = sett.website_name;
        this.mailto = sett.mailto;
        this.username = sett.username;
        this.password = sett.password;
        this.hostname = sett.hostname;
        this.port = sett.port;
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
        .replace("\"", "&#34;")
}

// composeMsg creates the string message that will be sent
function composeMsg(sett: Settings, msg: Message) : string {
    let webName = escapeHTML(sett.webName);
    let name = escapeHTML(msg.name);
    let mail = escapeHTML(msg.mail);
    let escapedMsg = escapeHTML(msg.msg).replace("\n", "<br>");
    return `<html><body>Message from ${webName}<br><br><b>Name:</b> ${name}<br><b>Email:</b> ${mail}<br><b>Message:</b> ${escapedMsg}</body></html>`
}

function send(sett: Settings, msg: Message) {
    let transporter = nodemailer.createTransport({
        host: sett.hostname,
        port: sett.port,
        secure: true,
        auth: {
            user: sett.username,
            pass: sett.password
        }
    });

    return transporter.sendMail({
        from: sett.username,
        to: sett.mailto,
        subject: "Message from " + sett.webName,
        html: composeMsg(sett, msg)
    });
}

send(new Settings(process.argv[2]), JSON.parse(process.argv[3])).catch((err) => process.stderr.write(err))
