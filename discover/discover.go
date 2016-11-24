package main

import "fmt"
import "github.com/pjediny/mndp/mndplib"

func printMNDPs(ch chan *mndplib.MNDPMessage) {
	for msg := range ch {
		fmt.Println(msg.String())
	}
}
func main() {
	ch := make(chan *mndplib.MNDPMessage)

	listener := mndplib.NewMNDPListener()
	listener.Listen(ch)

	go printMNDPs(ch)

	var response int
	for {
		fmt.Scanf("%c", &response) //<--- here
		switch response {
		case 'q':
			fmt.Println("Quiting")
			return
		case 'Q':
			fmt.Println("Quiting")
			return
		default:
			listener.RequestRefresh()
			fmt.Println("Requesting Refresh")
		}
	}

}
