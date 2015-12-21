// Steven Phillips / elimisteve
// 2015.12.21

package main

import "time"

var (
	I_AM_UP = "I am up"
)

type Status struct {
	// Set by sender
	Status string    `json:"status"`
	SentAt time.Time `json:"sent_at"`

	// Set by this server
	Token      string    `json:"token"`
	ReceivedAt time.Time `json:"received_at"`
}
