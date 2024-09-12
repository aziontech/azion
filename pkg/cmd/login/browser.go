package login

import (
	"context"
	"io"
	"net/http"

	msg "github.com/aziontech/azion-cli/messages/login"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

const (
	urlSsoNext = "https://sso.azion.com/login?next=cli"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func (l *login) browserLogin(srv Server) error {
	ctx, cancel := context.WithCancel(context.Background())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		paramValue := r.URL.Query().Get("c")
		_, _ = io.WriteString(w, msg.BrowserMsg)
		if paramValue != "" {
			tokenValue = paramValue
		}
		cancel()
	})

	err := l.openBrowser()
	if err != nil {
		return err
	}
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Error(msg.ErrorServerClosed.Error(), zap.Error(err))
		}
	}()

	<-ctx.Done() // wait for the signal to gracefully shutdown the server

	// gracefully shutdown the server:
	// waiting indefinitely for connections to return to idle and then shut down.
	err = srv.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (l *login) openBrowser() error {
	logger.FInfo(l.factory.IOStreams.Out, msg.VisitMsg)
	err := l.run(urlSsoNext)
	if err != nil {
		return err
	}
	return nil
}
