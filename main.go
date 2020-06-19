package main

func main(){
	go watch()

	client := &Client{
		Sip:   "127.0.0.1",
		Sport: 8881,
		ID:    "158E4B6005D5639552EF",
	}
	client.Run()
}
