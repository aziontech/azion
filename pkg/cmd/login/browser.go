package login

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/login"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

const (
	urlSsoNext     = "https://console.azion.com/login?next=cli"
	defaultPort    = 8080
	maxPortRetries = 10 // Try up to 10 different ports
)

// when it's a single test set true
var enableHandlerRouter = true

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

var (
	globalCtx    context.Context
	globalCancel context.CancelFunc
)

func initializeContext() {
	globalCtx, globalCancel = context.WithCancel(context.Background())
}

func shutdownContext() {
	if globalCancel != nil {
		globalCancel()
	}
}

func (l *login) browserLogin(srv Server) error {
	initializeContext()
	defer shutdownContext()

	if enableHandlerRouter {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			paramValue := r.URL.Query().Get("c")
			_, _ = io.WriteString(w, msg.BrowserMsg)
			if paramValue != "" {
				tokenValue = paramValue
			}
			globalCancel()
		})
	}

	// Find an available port starting from the default one
	port, server, err := l.findAvailablePort()
	if err != nil {
		return err
	}

	// Set the server to use
	l.server = server

	// Open browser with the selected port
	err = l.openBrowserWithPort(port)
	if err != nil {
		return err
	}

	go func() {
		err := l.server.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Error(msg.ErrorServerClosed.Error(), zap.Error(err))
		}
	}()

	<-globalCtx.Done() // wait for the signal to gracefully shutdown the server

	// gracefully shutdown the server:
	// waiting indefinitely for connections to return to idle and then shut down.
	err = l.server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (l *login) openBrowserWithPort(port int) error {
	logger.FInfo(l.factory.IOStreams.Out, fmt.Sprintf(msg.VisitMsg, port))

	// Append the callback port to the URL
	callbackURL := fmt.Sprintf("%s&callback_port=%d", urlSsoNext, port)

	err := l.run(callbackURL)
	if err != nil {
		return err
	}
	return nil
}

// findAvailablePort tries to find an available port starting from the default port
func (l *login) findAvailablePort() (int, Server, error) {
	port := defaultPort

	for i := 0; i < maxPortRetries; i++ {
		// Check if the port is available
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			// Port is available, close the listener and use this port
			_ = listener.Close()
			server := &http.Server{Addr: ":" + strconv.Itoa(port)}
			return port, server, nil
		}

		// Log that we're trying another port
		logger.Debug(fmt.Sprintf("Port %d is in use, trying port %d", port, port+1), zap.Error(err))

		// Try the next port
		port++
	}

	return 0, nil, fmt.Errorf("could not find an available port after %d attempts", maxPortRetries)
}
