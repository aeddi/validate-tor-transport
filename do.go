package tor

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/cretz/bine/tor"
	"github.com/ipsn/go-libtor"
)

var (
	// https://github.com/ipsn/go-libtor/issues/22
	running = false
	runLock sync.Mutex
)

// Do Perform the dial, it returns a message to print to the user.
// If debug != nil debug will be used as debug chan else it will be ignored.
func Do(debug io.Writer) string {
	runLock.Lock()
	if running {
		runLock.Unlock()
		return "An other instance is running"
	}
	running = true
	runLock.Unlock()
	defer func() {
		runLock.Lock()
		running = false
		runLock.Unlock()
	}()
	// Check for internet (maybe trigger an autorisation request ?)
	_, err := http.Get("http://golang.org/")
	if err != nil {
		return fmt.Sprintf("Can't dial golang : %q", err)
	}
	// Starting tor node
	t, err := tor.Start(nil, &tor.StartConf{ProcessCreator: libtor.Creator, DebugWriter: debug})
	if err != nil {
		return fmt.Sprintf("Failed to start tor : %q", err)
	}
	defer t.Close()
	dialer, err := t.Dialer(context.Background(), &tor.DialConf{})
	if err != nil {
		return fmt.Sprintf("Failed to create dialler : %q", err)
	}
	// Creating transport
	h := http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}
	r, err := h.Get("http://2vgsljkvmsyrxryc4fvb2gl5srazxfbrpszkimpej23yc37jsntp3did.onion:80/tor.txt")
	if err != nil {
		return fmt.Sprintf("Failed to do request : %q", err)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Sprintf("Failed to read body : %q", err)
	}
	r.Body.Close()
	return fmt.Sprintf("Status: %s, Proto: %s, Body: %s", r.Status, r.Proto, string(body))
}
