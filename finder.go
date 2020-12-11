package sfinder

import (
	"context"
	"strings"
	"time"
)

import "go.bug.st/serial/enumerator"



type Uart struct {
	Name string
	SerialNumber string
	Vendor string
	Product string
}

type MonitorConfig struct {
	Ctx context.Context
	Delay time.Duration
	Build func(name string) (interface{}, error)
}

type Finder interface {
	Find()(ports []Uart, err error)
	Monitor(config MonitorConfig)
	Remove(name string)
}

type finder struct {
	serialType SerialType
	config *MonitorConfig
	working bool

	store map[string]interface{}
	storeStatus map[string]bool

}

func New(vendorName SerialType) Finder {

	s := finder{
		serialType: vendorName,
	}

	var i Finder = &s
	return i
}

func (f *finder) Monitor(config MonitorConfig) {
	if f.working {
		return
	}
	if config.Ctx != nil && config.Build != nil {
		c := &MonitorConfig{
			Ctx:   config.Ctx,
			Delay: config.Delay,
			Build: config.Build,
		}
		if c.Delay < time.Second {
			c.Delay = time.Second
		}
		f.config = c
	}

	f.working = true
	f.store = make(map[string]interface{})
	f.storeStatus = make(map[string]bool)
	go func() {
		for {
			select {
			case <-f.config.Ctx.Done():
				break
			default:
				f.monitor()
				time.Sleep(f.config.Delay)
			}
		}
	}()
}

func (f *finder) Remove(name string) {
	if f.storeStatus != nil {
		f.storeStatus[name] = false
	}
}

func (f *finder) monitor() {

	for name, valid :=range f.storeStatus {
		if !valid {
			delete(f.store, name)
		}
	}

	ports, err := f.Find()
	if err != nil {
		return
	}

	for _, port := range ports {
		if _, ok := f.store[port.Name]; ok { // skip
			continue
		}
		elem, exp := f.config.Build(port.Name)
		if exp == nil {
			f.store[port.Name] = elem
			f.storeStatus[port.Name] = true
		}
	}
}

func (f *finder) Find() (ports []Uart, err error) {

	list, err := serialPorts()
	if err != nil {
		return
	}
	for _, uart := range list {
		if f.hit(uart) {
			ports = append(ports, uart)
		}
	}
	return
}

func (f *finder) hit(uart Uart) bool {

	if list, ok := vendor[f.serialType]; ok {
		for _, ven := range list {
			if uart.Vendor == ven {
				return true
			}
		}
	} else {
		return false
	}
	return false
}

func serialPorts() (list []Uart, err error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return
	}

	for _, port := range ports {
		if port.IsUSB {
			uart := Uart{
				Name:         port.Name,
				SerialNumber: port.SerialNumber,
				Vendor:       strings.ToLower(port.VID),
				Product:      strings.ToLower(port.PID),
			}
			list = append(list, uart)
		}
	}
	return
}