// Copyright 2017 By Rosco Nap (cloudrkt). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"time"
)

func main() {

	log.Println("Starting APC api on: http://<external-ip>:8000/v1/apc/")

	srv := http.Server{
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        NewMux(),
		Addr:           "0.0.0.0:8000",
	}

	log.Fatal(srv.ListenAndServe())
}
