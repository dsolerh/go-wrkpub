package utils

import (
	"fmt"
	"strconv"
	"strings"

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
	case "pre-patch":
		return func(s string) (string, error) {
			v, err := semver.NewVersion(s)
			if err != nil {
				return "", err
			}
			pre := v.Prerelease()
			if pre == "" {
				vpatch := v.IncPatch()
				vpre, err := vpatch.SetPrerelease("beta.0")
				if err != nil {
					return "", err
				}
				return "v" + vpre.String(), nil
			}
			parts := strings.Split(pre, ".")
			if len(parts) != 2 || parts[0] != "beta" {
				return "", fmt.Errorf("invalid prerelease appended %s", pre)
			}
			preV, err := strconv.Atoi(parts[1])
			if err != nil {
				return "", fmt.Errorf("error parsing prerelease version %s, %w", parts[1], err)
			}
			vpatch := v.IncPatch()
			vpre, err := vpatch.SetPrerelease(fmt.Sprintf("%s.%d", parts[0], preV+1))
			if err != nil {
				return "", err
			}

			return "v" + vpre.String(), nil
		}, nil
	}
	return nil, fmt.Errorf("invalid version type %s", vtype)
}
