package monitoring

import "github.com/prometheus/client_golang/prometheus"

var (
	SentMessagesCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sender_messages_sent_total",
			Help: "Total number of messages sent by sender service.",
		},
	)
)

func init() {
	prometheus.MustRegister(SentMessagesCounter)
}
