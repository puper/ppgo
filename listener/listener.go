package listener

import (
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	currentOrder = uintptr(2)
	listeners    = make(map[string]*orderedListener)
	mu           sync.Mutex
)

type orderedListener struct {
	listener net.Listener
	order    uintptr
}

func (this *orderedListener) File() *os.File {
	tl := this.listener.(*net.TCPListener)
	fl, _ := tl.File()
	return fl
}

func SetConfig(data string) {
	mu.Lock()
	defer mu.Unlock()
	for _, v := range strings.Split(data, ",") {
		currentOrder += 1
		log.Println(v)
		listeners[v] = &orderedListener{
			order: currentOrder,
		}
	}
}

func GetConfig() string {
	mu.Lock()
	defer mu.Unlock()
	result := make([]string, currentOrder-2)
	for k, v := range listeners {
		result[int(v.order)-3] = k
	}
	return strings.Join(result, ",")
}

func GetFiles(addrs []string) ([]*os.File, error) {
	mu.Lock()
	defer mu.Unlock()
	if addrs == nil {
		addrs = make([]string, len(listeners))
		i := 0
		for k, _ := range listeners {
			addrs[i] = k
			i += 1
		}
	}
	result := make([]*os.File, len(addrs))
	for i, addr := range addrs {
		l, err := getListener(addr)
		if err != nil {
			return nil, err
		}
		result[i], _ = l.(*net.TCPListener).File()
	}
	return result, nil
}

func GetListener(addr string) (net.Listener, error) {
	mu.Lock()
	defer mu.Unlock()
	return getListener(addr)
}

func getListener(addr string) (net.Listener, error) {
	var err error
	if listener, ok := listeners[addr]; ok {
		if listener.listener == nil {
			f := os.NewFile(uintptr(listener.order), "")
			listener.listener, err = net.FileListener(f)
			if err != nil {
				l, err := newListener(addr)
				if err == nil {
					listeners[addr].listener = l
				}
				return l, err
			}
		}
		return listener.listener, nil
	}
	return newListener(addr)
}

func newListener(addr string) (net.Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	currentOrder += 1
	listeners[addr] = &orderedListener{
		listener: l,
		order:    currentOrder,
	}
	return l, nil
}
