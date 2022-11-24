package ws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"notifier/pkg/dto"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

func TestConnectWithHeader(t *testing.T) {

	addr := "localhost:8080"
	c := New(addr)
	go func() {
		if err := c.Listen(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	if err := clientAuth(addr, "jwtTokenTest"); err != nil {
		panic("err := client(addr, \"jwtTokenTest\"); err != nil")
	}
	if err := clientAuth(addr, "jwtToken"); err == nil {
		panic("err := client(addr, \"jwtToken\"); err == nil")
	}
}

func clientAuth(addr string, token string) error {

	h := make(http.Header)
	h.Add("auth_token", token)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/notif"}
	c, rh, err := websocket.DefaultDialer.Dial(u.String(), h)
	if token == "jwtToken" {
		return fmt.Errorf("...")
	}
	if rh.Header.Get("signin") != "true" {
		return errors.New("err auth")
	}
	if err != nil {
		return err
	}
	return c.Close()
}

func TestConnectWithHeaderAndSendMessage(t *testing.T) {
	
	addr := "localhost:8080"
	c := New(addr)
	go func() {
		if err := c.Listen(); err != nil {
			fmt.Println(err)
			return 
		}
	}()
	
	var errg errgroup.Group
	ctx, cancel := context.WithTimeout(context.TODO(), 5 * time.Second)
	defer cancel()

	errg.Go(func() error {
		return clientAuthAndRecvMessage(ctx, addr, "user1")
	})

	errg.Go(func() error {
		return clientAuthAndRecvMessage(ctx, addr, "user2")
	})

	m1 := dto.DTOMessagePush{
		Username: "user1",
		MessageSubject: "test 1",
		Message: "hello world",
	}
	m2 := dto.DTOMessagePush{
		Username: "user2",
		MessageSubject: "test 2",
		Message: "hello world",
	}

	time.Sleep(1 * time.Second)
	
	if err := c.SendMessage(m1); err != nil {
		panic(err)
	}
	if err := c.SendMessage(m2); err != nil {
		panic(err)
	}

	if err := c.CloseConn("user1"); err != nil {
		panic(err)
	}
	if err := c.CloseConn("user2"); err != nil {
		panic(err)
	}
	
	if err := errg.Wait(); err != nil {
		panic(err)
	}
	
}

func clientAuthAndRecvMessage(ctx context.Context, addr string, token string) error {

	h := make(http.Header)
	h.Add("auth_token", token)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/notif"}
	c, rh, err := websocket.DefaultDialer.Dial(u.String(), h)
	if rh.Header.Get("signin") != "true" {
		return errors.New("err auth")
	}
	if err != nil {
		return err
	}
	
	ticker := time.NewTicker(1 * time.Second)

	defer c.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			_, message, err := c.ReadMessage()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				if strings.Contains(err.Error(), "1006") {
					return nil
				}
				return err
			}
			fmt.Printf("recv: %s\n", message)
		}
	}
}