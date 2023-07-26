package slacker

// isMessageInThread determines if a message is in a thread
func isMessageInThread(threadTimestamp string, messageTimestamp string) bool {
	if threadTimestamp == "" || threadTimestamp == messageTimestamp {
		return false
	}
	return true
}
