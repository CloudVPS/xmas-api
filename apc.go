// Copyright 2017 By Rosco Nap (cloudrkt). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type apc struct {
	com   string // community string
	ip    string // Ip
	oid   string // MIB of powerports
	ports int    // Number of available ports
	loc   string // Location
	state int    // State 1 is ON 2 is OFF
}

// getAPC on location specified. Should be replaced with configuration file.
func getAPC(loc string) (apc, error) {
	apcs := []apc{
		apc{
			com:   "community-string",
			ip:    "0.0.0.0",
			oid:   "1.3.6.1.4.1.318.1.1.12.3.3.1.1.4.",
			ports: 8,
			loc:   "cloud",
		},
	}

	for _, apc := range apcs {
		if apc.loc == loc {
			return apc, nil
		}
	}

	return apc{}, errors.New("No APC found on specified location")
}

func (a *apc) validatePort(port int) (int, error) {

	for i := 1; i < a.ports; i++ {
		if i == port {
			return port, nil
		}
	}

	return port, errors.New("Missing or invalid APC port number")
}

// validateState takes the state string and returns the numerical representation
// and an error.
func validateState(state string) (int, error) {

	switch s := strings.ToUpper(state); s {
	case "ON":
		return 1, nil
	case "OFF":
		return 2, nil
	case "FLIP":
		return 99, nil
	default:
		return 0, errors.New("Missing or invalid APC state")
	}
}

// getCurrentState of port
func (a *apc) getCurrentState(port int) (int, error) {
	state, err := a.snmpGet(port)
	if err != nil {
		return 0, err
	}
	a.state = state
	return a.state, nil
}

// switchOn switch port on
func (a *apc) switchOn(port int) error {
	err := a.snmpSet(port, 1)
	if err != nil {
		return err
	}
	return nil
}

// switchOff switch port off
func (a *apc) switchOff(port int) error {
	err := a.snmpSet(port, 2)
	if err != nil {
		return err
	}
	return nil
}

// switchFlip port on and off
func (a *apc) switchFlip(port int) error {

	c, err := a.getCurrentState(port)
	if err != nil {
		return err
	}

	if c == 1 {
		a.switchOff(port)
	}

	if c == 2 {
		a.switchOn(port)
	}

	return nil
}

// snmpGet get snmp status with snmp commands on system
func (a *apc) snmpGet(port int) (int, error) {

	cmdoutput := []byte{}
	cmdName := "/usr/bin/snmpget"
	cmdArgs := []string{"-v1",
		"-c",
		a.com,
		a.ip,
		a.oid + strconv.Itoa(port),
	}

	cmdoutput, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 0, err
	}

	output := string(cmdoutput)
	cstate, err := strconv.Atoi(output[len(output)-2 : len(output)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 0, err
	}
	return cstate, nil
}

// snmpSet set snmp status with snmp command on system
func (a *apc) snmpSet(port int, state int) error {

	args := []string{"-v1",
		"-c",
		a.com,
		a.ip,
		a.oid + strconv.Itoa(port),
		"integer",
		strconv.Itoa(state),
	}

	cmd := exec.Command("/usr/bin/snmpset", args...)

	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
