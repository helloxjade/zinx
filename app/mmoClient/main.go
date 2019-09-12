package main

func main() {
	client := NewClient("127.0.0.1", 8999)
	client.Start()
}
