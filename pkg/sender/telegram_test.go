package sender

import "testing"

func TestTelegram_createMessage(t *testing.T) {
	testTelegram := Telegram{
		URL:             "mywebsite.com",
		RecaptchaSecret: "",
		ChatId:          "123456",
		BotToken:        "654321abc",
	}

	inputName := "This is an <script src=\"hack.js\"></script>unsafe name"
	inputMail := "-123ABCabc!#$%&'*+/=?^_`{|}.~@-456DEFdef_.~.GHIghi"
	inputMsg := "<b>You've win $2000.</b>\nThis is totally not an scam and unsafe message.\n<a href=\"hack.js\">Click here to get hacked</a>"
	expectedResult := "Message from mywebsite.com\n" +
		"\n" +
		"<b>Name</b>: This is an &lt;script src=&#34;hack.js&#34;&gt;&lt;/script&gt;unsafe name\n" +
		"<b>Email</b>: -123ABCabc!#$%&amp;&#39;*+/=?^_`{|}.~@-456DEFdef_.~.GHIghi\n" +
		"<b>Message</b>: &lt;b&gt;You&#39;ve win $2000.&lt;/b&gt;\n" +
		"This is totally not an scam and unsafe message.\n" +
		"&lt;a href=&#34;hack.js&#34;&gt;Click here to get hacked&lt;/a&gt;"

	result := testTelegram.createMessage(inputName, inputMail, inputMsg)
	if result != expectedResult {
		t.Errorf("Error creating message.\n-> Expected message: \"%s\"\n-> Message found: \"%s\"", expectedResult, result)
	}
}
