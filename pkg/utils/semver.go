package utils

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

func SemverUpdater(vtype string) (func(string) (string, error), error) {
	switch vtype {
	case "mayor":
		return func(s string) (string, error) {
			v, err := semver.NewVersion(s)
			if err != nil {
				return "", err
			}
			return "v" + v.IncMajor().String(), nil
		}, nil
	case "minor":
		return func(s string) (string, error) {
			v, err := semver.NewVersion(s)
			if err != nil {
				return "", err
			}
			return "v" + v.IncMinor().String(), nil
		}, nil
	case "patch":
		return func(s string) (string, error) {
			v, err := semver.NewVersion(s)
			if err != nil {
				return "", err
			}
			return "v" + v.IncPatch().String(), nil
		}, nil
	}
	return nil, fmt.Errorf("invalid version type %s", vtype)
}
