/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/dsolerh/go-wrkpub/pkg/config"
	"github.com/dsolerh/go-wrkpub/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	shouldPush     bool
	updateStrategy string
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Prepares and publishes packages on the workspace",
	Long: `Prepares packages for publish in a worspace context,
the packages to publish can be specified as well
as if should push the changes to the scv.

Note: only git for the moment`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// check if there're uncommitted changes in the repo
		uncomm, err := utils.HasUncommittedChanges()
		if err != nil {
			log.Fatal(err)
		}
		if uncomm {
			log.Fatal("uncommitted changes in the repo")
		}

		vupdater, err := utils.SemverUpdater(updateStrategy)
		if err != nil {
			log.Fatal(err)
		}

		cfg, err := config.LoadPublishConfig(cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		pkgNames, err := cfg.GetPackagesNames(args)
		if err != nil {
			log.Fatal(err)
		}

		err = cfg.UpdatePackagesVersion(pkgNames, vupdater)
		if err != nil {
			log.Fatal(err)
		}
		// get the tag versions to use
		tagVersions := cfg.GetPackagesTags(pkgNames)

		workPkgPaths := cfg.GetWorkPkgPaths()
		pkgPaths := cfg.GetPkgPaths(pkgNames)

		// copy the packages to the root of the workspace
		if err := utils.CopyPackagesToRoot(cfg, pkgNames); err != nil {
			log.Fatal(err)
		}
		var cleanUpPackageInRoot = func() {
			log.Println("removing packages from root...")
			// remove the packages from root (cleanup)
			if err := utils.RemovePackagesFromRoot(pkgNames); err != nil {
				log.Println(err)
			}
		}

		// remove the old packages from the go.work and add the new
		if err := utils.UpdateWorkspacePackages(pkgPaths, workPkgPaths); err != nil {
			cleanUpPackageInRoot()
			log.Fatal(err)
		}
		var cleanUpWorkspace = func() {
			log.Println("reverting workspace packages...")
			// revert the workspace packages
			if err := utils.UpdateWorkspacePackages(workPkgPaths, pkgPaths); err != nil {
				log.Println(err)
			}
		}

		// update the go.mod files of the new packages
		if err := utils.UpdatePackageMods(cfg, pkgNames); err != nil {
			cleanUpPackageInRoot()
			cleanUpWorkspace()
			log.Fatal(err)
		}

		// commit all changes
		if err := utils.CommitChanges(utils.GetPublishCommitMessage(tagVersions)); err != nil {
			cleanUpPackageInRoot()
			cleanUpWorkspace()
			log.Fatal(err)
		}
		// TODO: revert the commit if something fails

		// tag the versions
		if err := utils.TagPackagesVersion(tagVersions); err != nil {
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
		if err = cfg.SaveConfig(); err != nil {
			log.Fatal(err)
		}

		// TODO: delete the commits if something fails
		if err := utils.CommitChanges(utils.CleanupCommit); err != nil {
			log.Fatal(err)
		}

		if shouldPush {
			if err := utils.PushChanges(tagVersions); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().BoolVarP(&shouldPush, "push", "p", false, "if set to true will push the changes to the scv repo")
	publishCmd.Flags().StringVarP(&updateStrategy, "update", "u", "patch", "the update strategy for the packages version (mayor|minor|patch)")
}
