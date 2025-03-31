package main

func main() {
	server := initWebServer()
	server.Run(":8080")
}
