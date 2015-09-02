package socket

// Message is the normal data for messages passed on the console socket.
type Message struct {
	// Type of message being passed
	Type string `json:"type"`
}

// TerminalRequest is the normal data for messages passing a pseudoterminal master.
type TerminalRequest struct {
	Message

	// Container ID for the container whose pseudoterminal master is being set.
	Container string `json:"container"`
}

// Response is the normal data for response messages.
type Response struct {
	Message

	// Message is a phrase describing the response.
	Message string `json:"message,omitempty"`
}
