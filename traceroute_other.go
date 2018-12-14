// +build !darwin
// +build !linux

package main

import "errors"

func resolveClosestRouteIP(myIP string) (string, error) {
	return "", errors.New("Not support this platform")
}
