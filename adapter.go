package main

type Adapter interface {
	Send(m string)
	Connect()
}
