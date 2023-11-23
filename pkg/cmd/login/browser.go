package login

import (
	"context"
	"io"
	"log"
	"net/http"

	msg "github.com/aziontech/azion-cli/messages/login"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/skratchdot/open-golang/open"
)

const (
	urlSsoNext = "https://sso.azion.com/login?next=cli"
)

func browserLogin(f *cmdutil.Factory) error {

	ctx, cancel := context.WithCancel(context.Background())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		paramValue := r.URL.Query().Get("c")
		io.WriteString(w, "You may now close this page and return to your terminal")
		tokenValue = paramValue
		cancel()
	})

	srv := &http.Server{Addr: ":8080"}
	err := openBrowser(f)
	if err != nil {
		return err
	}
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Println(err)
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

func openBrowser(f *cmdutil.Factory) error {
	logger.FInfo(f.IOStreams.Out, msg.VisitMsg)
	err := open.Run(urlSsoNext)
	if err != nil {
		return err
	}
	return nil
}
