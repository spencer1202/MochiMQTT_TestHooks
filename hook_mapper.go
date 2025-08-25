package hookmap

import (
	"bytes"
	"fmt"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type TestHook struct {
	mqtt.HookBase
}

// List of functions this hook implements
func (h TestHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate,
		mqtt.OnACLCheck,
		mqtt.OnSubscribe,
		mqtt.OnSubscribed,
		mqtt.OnPublish,
		mqtt.OnSelectSubscribers,
	}, []byte{b})
}

// Show Packet
func showPacket(pk packets.Packet) {
	packetType := packets.PacketNames[pk.FixedHeader.Type]
	fmt.Printf("Show Packet: type %v\n", packetType)

	switch packetType {
	case "Connect":
		fmt.Printf(" UsernameFlag %v\tUsername: %s\tUserID: %s\n",
			pk.Connect.UsernameFlag, string(pk.Connect.Username), pk.Connect.ClientIdentifier)
		fmt.Printf(" PasswordFlag %v\tPassword: %s",
			pk.Connect.PasswordFlag, string(pk.Connect.Password))
	case "Publish":
		fmt.Printf(" Topic %s\n", pk.TopicName)
	}
}

// Show Client
func showClient(cl *mqtt.Client) {
	fmt.Println("Show Client")
	if cl == nil {
		fmt.Println(" no client")
	} else {
		fmt.Printf(" Username: %s\tUserID: %s\n",
			string(cl.Properties.Username), cl.ID)
	}
}

// OnConnectAuthenticate
func (h TestHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	fmt.Println("\nHOOK: OnConnectAuthenticate")
	showClient(cl)
	showPacket(pk)
	return true
}

// OnSubscribe
func (h TestHook) OnSubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	fmt.Println("\nHOOK: OnSubscribe")
	showClient(cl)
	showPacket(pk)
	return pk
}

// OnACLCheck
func (h TestHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	fmt.Println("\nHOOK: OnACLCheck")
	fmt.Printf(" topic: %v  write: %v\n", topic, write)
	showClient(cl)

	return true
}

// OnSubscribed
func (h TestHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	fmt.Println("\nHOOK: OnSubscribed")
	showClient(cl)
	showPacket(pk)
}

// OnPublish
func (h TestHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	fmt.Println("\nHOOK: OnSubscribed")
	showClient(cl)
	showPacket(pk)
	return pk, nil
}

// OnSelectSubscribers intercepts an incoming publish message and allows us to modify the list
// of subscribers it will be sent out to.
func (h TestHook) OnSelectSubscribers(subs *mqtt.Subscribers, pk packets.Packet) *mqtt.Subscribers {
	fmt.Println("\nHOOK: OnSelectSubscribers")
	showPacket(pk)
	fmt.Printf("Subscribers: ")
	for uID := range subs.Subscriptions {
		fmt.Printf("%v ", uID)
	}
	fmt.Printf("\n")
	return subs
}
