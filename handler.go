// Copyright 2017 By Rosco Nap (cloudrkt). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// NewMux router
func NewMux() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/", indexHandler)
	mux.HandleFunc("/v1/apc/", portHandler)

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Welcome to the Christmas lights API...")
	return
}

func portHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	apc, err := getAPC(r.URL.Query().Get("loc"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, _ := strconv.Atoi(r.URL.Query().Get("port"))
	port, err := apc.validatePort(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state, err := validateState(r.URL.Query().Get("state"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch state {
	case 1: // ON
		err = apc.switchOn(port)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case 2: // OFF
		err := apc.switchOff(port)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case 99: // FLIP
		err := apc.switchFlip(port)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, http.StatusOK)
	return
}
