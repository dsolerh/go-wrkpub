package main

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type PublishConfig struct {
	configPath string
	Repo       string                  `yaml:"repo"`
	Root       string                  `yaml:"root"`
	Packages   map[string]*PackageInfo `yaml:"packages"`
}

type PackageInfo struct {
	WorkName string `yaml:"work_name"`
	PkgName  string `yaml:"pkg_name"`
	Version  string `yaml:"version"`
}

func loadPublishConfig(fname string) (*PublishConfig, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}

	config := new(PublishConfig)
	config.configPath = fname

	err = yaml.NewDecoder(f).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("error decoding the config file: %w", err)
	}

	return config, nil
}

func (c *PublishConfig) getPackagesNames(args []string) ([]string, error) {
	if len(args) == 0 {
		packageNames := make([]string, 0, len(c.Packages))
		for pkgName := range c.Packages {
			packageNames = append(packageNames, pkgName)
		}
		return packageNames, nil
	}

	for _, arg := range args {
		if _, exist := c.Packages[arg]; !exist {
			return nil, fmt.Errorf("package '%s' is not defined in %s", arg, c.configPath)
		}
	}
	return args, nil
}

func (c *PublishConfig) updatePackagesVersion(packages []string, updater func(string) (string, error)) (err error) {
	for _, pkgName := range packages {
		pkg := c.Packages[pkgName]
		pkg.Version, err = updater(pkg.Version)
		if err != nil {
			return
		}
	}
	return
}

func (c *PublishConfig) getPackagesTags(packages []string) []string {
	tags := make([]string, 0, len(packages))
	for _, pkgName := range packages {
		tags = append(tags, fmt.Sprintf("%s/%s", pkgName, c.Packages[pkgName].Version))
	}
	return tags
}

func (c *PublishConfig) getWorkPkgPaths() []string {
	packages := make([]string, 0, len(c.Packages))
	for _, pkg := range c.Packages {
		packages = append(packages, path.Join(c.Root, pkg.WorkName))
	}
	return packages
}

func (c *PublishConfig) getPkgPaths(packages []string) []string {
	pkgs := make([]string, 0, len(c.Packages))
	for _, pkg := range packages {
		pkgs = append(pkgs, path.Join(c.Root, pkg))
	}
	return pkgs
}

func (c *PublishConfig) saveConfig() error {
	f, err := os.Create(c.configPath)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(f).Encode(c)
}
