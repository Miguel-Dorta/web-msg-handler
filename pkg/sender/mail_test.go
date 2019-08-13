package sender

import "testing"

func TestMail_lfToBr(t *testing.T) {
	input := "This is a test with \nUnix line breaks and \r\nWindows line breaks"
	expectedResult := "This is a test with <br>Unix line breaks and <br>Windows line breaks"
	actualResult := lfToBr(input)
	if actualResult != expectedResult {
		t.Errorf("Error converting line finalizers to <br>\n-> Expected result: %s\n-> Result found: %s", expectedResult, actualResult)
	}
}

func TestMail_createMessage(t *testing.T) {
	testMail := Mail{
		URL:             "mywebsite.com",
		RecaptchaSecret: "",
		Mailto:          "test@mywebsite.com",
		Username:        "test@mywebsite.com",
		Password:        "1234567890",
		Hostname:        "mail.mywebsite.com",
		Port:            "587",
	}

	inputName := "This is an <script src=\"hack.js\"></script>unsafe name"
	inputMail := "-123ABCabc!#$%&'*+/=?^_`{|}.~@-456DEFdef_.~.GHIghi"
	inputMsg := "<b>You've win $2000.</b>\nThis is totally not an scam and unsafe message.\n<a href=\"hack.js\">Click here to get hacked</a>"
	expectedResult := "From: test@mywebsite.com\r\n" +
		"To: test@mywebsite.com\r\n" +
		"Subject: Message from mywebsite.com\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		"<html><body><b>Name</b>: This is an &lt;script src=&#34;hack.js&#34;&gt;&lt;/script&gt;unsafe name<br><b>Email</b>: -123ABCabc!#$%&amp;&#39;*+/=?^_`{|}.~@-456DEFdef_.~.GHIghi<br><b>Message</b>: &lt;b&gt;You&#39;ve win $2000.&lt;/b&gt;<br>This is totally not an scam and unsafe message.<br>&lt;a href=&#34;hack.js&#34;&gt;Click here to get hacked&lt;/a&gt;</body></html>\r\n"

	result := string(testMail.createMessage(inputName, inputMail, inputMsg))
	if result != expectedResult {
		t.Errorf("Error creating message.\n-> Expected message: \"%s\"\n-> Message found: \"%s\"", expectedResult, result)
	}
}
