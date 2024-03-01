package main

import (
	"fmt"
	"sync"
	"time"
)

type Barber struct {
	id          int
	busy        bool
	waitingRoom *WaitingRoom
}

type Client struct {
	id int
}

type WaitingRoom struct {
	capacity int
	clients  []*Client
	mutex    sync.Mutex
}

type Barbershop struct {
	closingTime time.Time
	barbers     []*Barber
	waitingRoom *WaitingRoom
	mutex       sync.Mutex
}

func (b *Barber) cutHair(client *Client) {
	fmt.Printf("Barber %d is cutting hair for client %d\n", b.id, client.id)
	time.Sleep(time.Second * 3) // haircut duration
}

func (wr *WaitingRoom) addClient(client *Client) bool {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	if len(wr.clients) < wr.capacity {
		wr.clients = append(wr.clients, client)
		return true
	}
	return false
}

func (wr *WaitingRoom) removeClient() *Client {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	if len(wr.clients) > 0 {
		client := wr.clients[0]
		wr.clients = wr.clients[1:]
		return client
	}
	return nil
}

func (bs *Barbershop) open() {
	for _, barber := range bs.barbers {
		go func(b *Barber) {
			for {
				if client := bs.waitingRoom.removeClient(); client != nil {
					b.busy = true
					b.cutHair(client)
					b.busy = false
				} else {
					if time.Now().After(bs.closingTime) {
						break
					}
					time.Sleep(time.Second)
				}
			}
		}(barber)
	}
}
func main() {
	closingTime := time.Now().Add(time.Second * 20) // 20 seconds of barbershop operation
	wr := &WaitingRoom{capacity: 3}
	bs := &Barbershop{
		closingTime: closingTime,
		waitingRoom: wr,
	}

	for i := 0; i < 2; i++ { // Two barbers in the shop
		bs.barbers = append(bs.barbers, &Barber{id: i + 1, waitingRoom: wr})
	}

	bs.open()

	//  clients arriving
	for i := 1; time.Now().Before(closingTime); i++ {
		client := &Client{id: i}
		if !wr.addClient(client) {
			fmt.Printf("Client %d leaves because the waiting room is full\n", i)
		} else {
			fmt.Printf("Client %d enters the waiting room\n", i)
		}
		time.Sleep(time.Second) // New client reach every second
	}

	fmt.Println("Barbershop is closed")
}
