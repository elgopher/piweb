// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, port int) {
	addr := fmt.Sprintf("localhost:%d", port)

	log.Printf("Starting web server on http://%s", addr)

	server := &http.Server{Addr: addr, Handler: newHandler()}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Print("Server stopped")
				return
			}
			log.Fatalf("Problem starting web server: %v", err)
		}
	}()

	<-ctx.Done()
	_ = server.Close()
}
