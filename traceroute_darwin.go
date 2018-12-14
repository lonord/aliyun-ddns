// +build darwin

package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func resolveClosestRouteIP(myIP string) (string, error) {
	_, err := exec.LookPath("traceroute")
	if err != nil {
		return "", errors.New("Could not find traceroute in PATH")
	}
	cmd := exec.Command("traceroute", "-m", "1", "-n", myIP)
	b, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Exec traceroute failed: %s", err.Error())
	}
	lines := strings.Split(string(b), "\n")
	if len(lines) < 1 {
		return "", errors.New("Exec traceroute failed: format error")
	}
	cols := strings.Split(lines[0], " ")
	if len(cols) < 4 {
		return "", errors.New("Exec traceroute failed: format error")
	}
	return cols[3], nil
}
