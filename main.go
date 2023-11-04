package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/kevpar/repl-go"
	"golang.org/x/sys/windows"
)

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	return repl.Run(&state{systems: make(map[string]*cs)}, allCommands(), func(state *state) string { return state.def })
}

func pumpNotificationsUntil(ctx context.Context, ch <-chan *hcsNotification, desiredType uint32) (*hcsNotification, error) {
	for {
		select {
		case n := <-ch:
			if n.notificationType == desiredType {
				return n, nil
			}
			fmt.Println(stringNotification(n))
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func stringNotification(n *hcsNotification) string {
	b, err := json.MarshalIndent(n.notificationData, "\t", "\t")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("notification:\n\tnotificationType: %d\n\tnotificationStatus: %d\n\tnotificationData: %s\n", n.notificationType, n.notificationStatus, b)
}

type hcsNotification struct {
	notificationType   uint32
	notificationStatus uint32
	notificationData   *struct {
		Error       int32
		ErrorEvents []struct {
			Data []struct {
				Type  string
				Value string
			}
			EventId  uint
			Message  string
			Provider string
		}
		ErrorMessage string
	}
}

var notificationWatcher struct {
	m     sync.Mutex
	chans map[uintptr]chan<- *hcsNotification
}

func registerNotification(context uintptr, ch chan<- *hcsNotification) {
	notificationWatcher.m.Lock()
	defer notificationWatcher.m.Unlock()
	if notificationWatcher.chans == nil {
		notificationWatcher.chans = make(map[uintptr]chan<- *hcsNotification)
	}
	notificationWatcher.chans[context] = ch
}

func computeSystemCallback(notificationType uint32, context uintptr, notificationStatus uint32, notificationData *uint16) uintptr {
	n := &hcsNotification{
		notificationType:   notificationType,
		notificationStatus: notificationStatus,
	}
	if notificationData != nil {
		if err := json.Unmarshal([]byte(windows.UTF16PtrToString(notificationData)), &n.notificationData); err != nil {
			panic(err)
		}
	}
	notificationWatcher.m.Lock()
	defer notificationWatcher.m.Unlock()
	notificationWatcher.chans[context] <- n
	return 0
}
