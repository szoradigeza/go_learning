package ws

import (
	"log"
)

type ClientConnection struct {
	client *Client
	id     int
}

type FormPoolController struct {
	Register   chan *ClientConnection
	Unregister chan *ClientConnection
	Broadcast  chan Message
	FormPools  []map[int]Pool
}

func NewFormPoolController() *FormPoolController {
	return &FormPoolController{
		Register:   make(chan *ClientConnection),
		Unregister: make(chan *ClientConnection),
		FormPools:  make([]map[int]Pool, 0),
		Broadcast:  make(chan Message),
	}
}

func (controller *FormPoolController) KeyExists(key int) bool {
	for _, poolMap := range controller.FormPools {
		if _, ok := poolMap[key]; ok {
			return true
		}
	}
	return false
}

func (controller *FormPoolController) FindPool(key int) (*Pool, bool) {
	for _, poolMap := range controller.FormPools {
		if val, ok := poolMap[key]; ok {
			return &val, true
		}
	}
	return &Pool{}, false
}

func (controller *FormPoolController) RemovePool(formId int) {
	for i, formPool := range controller.FormPools {
		for key := range formPool {
			if key == formId {
				controller.FormPools = append(controller.FormPools[:i], controller.FormPools[i+1:]...)
				return
			}
		}
	}
}

func (fpc *FormPoolController) Start() {
	for {
		select {
		case cc := <-fpc.Register:
			pool, found := fpc.FindPool(cc.id)
			if found {
				cc.client.Pool = pool
				pool.Register <- cc.client
			} else {
				log.Printf("New Pool")
				pool = NewPool()
				go pool.Start()
				cc.client.Pool = pool

				newMap := make(map[int]Pool)
				newMap[cc.id] = *pool
				fpc.FormPools = append(fpc.FormPools, newMap)
				pool.Register <- cc.client
			}
		case cc := <-fpc.Unregister:
			log.Println(cc)
			return
		}
	}
}
