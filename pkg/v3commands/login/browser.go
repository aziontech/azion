package login

import (
	"context"
	"net/http"

	msg "github.com/aziontech/azion-cli/messages/login"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

const (
	urlSsoNext = "https://sso.azion.com/login?next=cli"
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

			htmlResponse := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<link rel="icon" type="image/png" sizes="32x32" href="https://avatars.githubusercontent.com/u/6660972?s=200&v=4">
			<title>Azion</title>
			<style>
				body {
					display: flex;
					flex-direction: column;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
				}
				.container {
					text-align: center;
				}
				.logo {
					width: 100px;
				}
				.text {
					color: #000;
					font-size: 12px;
					margin-top: 20px;
				}
				.footer {
					position: fixed;
					bottom: 10px;
					text-align: center;
					width: 100%;
					font-size: 14px;
					color: #888;
				}
				.footer a {
					color: #000;
					text-decoration: none;
				}
			</style>
		</head>
		<body style="background: #ffffff;"> 
			<div class="container">
				<img src="https://avatars.githubusercontent.com/u/6660972?s=200&v=4" alt="Logo" class="logo">
				<div class="text">Authenticated, you can now close this page and return to your terminal</div>
			</div>
			<div class="footer">
				<p>&copy; 2024 <a href="https://github.com/aziontech", >Azion</a>. Licensed under the <a href="https://opensource.org/license/mit" target="_blank" style="color: #000;">Mit License</a>. All rights reserved.</p>
			</div>
		</body>
		</html>`

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(htmlResponse))
			if err != nil {
				logger.Error("Error render html", zap.Error(err))
			}

			if paramValue != "" {
				tokenValue = paramValue
			}
			globalCancel()
		})
	}

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

	<-globalCtx.Done() // wait for the signal to gracefully shutdown the server

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
