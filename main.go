package main

import (
	"flag"
	"log"
)

func main() {
	// check if there're uncommitted changes in the repo
	uncomm, err := hasUncommittedChanges()
	if err != nil {
		log.Fatal(err)
	}
	if uncomm {
		log.Fatal("uncommitted changes in the repo")
	}

	v := flag.String("v", "patch", "the update to the package version (mayor|minor|patch)")
	fname := flag.String("f", "publish.go.yml", "the config file for the update")
	push := flag.Bool("push", false, "if set to true will push the changes to the git repo")
	flag.Parse()

	vupdater, err := semverUpdater(*v)
	if err != nil {
		log.Fatal(err)
	}

	config, err := loadPublishConfig(*fname)
	if err != nil {
		log.Fatal(err)
	}

	pkgNames, err := config.getPackagesNames(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	err = config.updatePackagesVersion(pkgNames, vupdater)
	if err != nil {
		log.Fatal(err)
	}
	// get the tag versions to use
	tagVersions := config.getPackagesTags(pkgNames)

	workPkgPaths := config.getWorkPkgPaths()
	pkgPaths := config.getPkgPaths(pkgNames)

	// copy the packages to the root of the workspace
	if err := copyPackagesToRoot(config, pkgNames); err != nil {
		log.Fatal(err)
	}
	var cleanUpPackageInRoot = func() {
		log.Println("removing packages from root...")
		// remove the packages from root (cleanup)
		if err := removePackagesFromRoot(pkgNames); err != nil {
			log.Println(err)
		}
	}

	// remove the old packages from the go.work and add the new
	if err := updateWorkspacePackages(pkgPaths, workPkgPaths); err != nil {
		cleanUpPackageInRoot()
		log.Fatal(err)
	}
	var cleanUpWorkspace = func() {
		log.Println("reverting workspace packages...")
		// revert the workspace packages
		if err := updateWorkspacePackages(workPkgPaths, pkgPaths); err != nil {
			log.Println(err)
		}
	}

	// update the go.mod files of the new packages
	if err := updatePackageMods(config, pkgNames); err != nil {
		cleanUpPackageInRoot()
		cleanUpWorkspace()
		log.Fatal(err)
	}

	// commit all changes
	if err := commitChanges(getPublishCommitMessage(tagVersions)); err != nil {
		cleanUpPackageInRoot()
		cleanUpWorkspace()
		log.Fatal(err)
	}
	// TODO: revert the commit if something fails

	// tag the versions
	if err := tagPackagesVersion(tagVersions); err != nil {
		cleanUpPackageInRoot()
		cleanUpWorkspace()
		// TODO: remove commits
		log.Fatal(err)
	}
	// TODO: revert the tags if something fails

	// clean up the changes to the repo
	cleanUpPackageInRoot()
	cleanUpWorkspace()

	// save the version changes
	if err = config.saveConfig(); err != nil {
		log.Fatal(err)
	}

	// TODO: delete the commits if something fails
	if err := commitChanges(cleanupCommit); err != nil {
		log.Fatal(err)
	}

	if *push {
		if err := pushChanges(tagVersions); err != nil {
			log.Fatal(err)
		}
	}
}
