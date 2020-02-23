package main

import "fmt"
import "github.com/pjediny/mndp/pkg/mndp"

func printMNDPs(ch chan *mndp.Message) {
	for msg := range ch {
		fmt.Println(msg.String())
	}
}
func main() {
	ch := make(chan *mndp.Message)

	listener := mndp.NewListener()
	listener.Listen(ch)

	go printMNDPs(ch)

	var response int
	for {
		_, _ = fmt.Scanf("%c", &response) //<--- here
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
