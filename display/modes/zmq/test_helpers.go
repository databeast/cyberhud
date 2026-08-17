package zmq

import "github.com/databeast/cyberhud/display/modes/zmq/content"

var ZMQRegistryExported = zmqRegistry

func SeedTestMessages() {
	content.ClearMessagesForTest()
	for _, msg := range []string{
		`{"temp":42,"status":"ok"}`,
		"alerts: link up",
		"metrics: cpu=17 mem=42",
	} {
		content.PushMessageForTest(msg)
	}
}

func ResetTestState() {
	content.ClearMessagesForTest()
	content.SetPolicy(content.DefaultPolicy())
}
