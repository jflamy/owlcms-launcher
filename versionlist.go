package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var versionList *widget.List

func getAllInstalledVersions() []string {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil
	}

	versionPattern := regexp.MustCompile(`^\d+\.\d+\.\d+(?:-(?:rc|alpha|beta)(?:\d+)?)?$`)
	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && versionPattern.MatchString(entry.Name()) {
			versions = append(versions, entry.Name())
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	return versions
}

func findLatestInstalled() string {
	entries, err := os.ReadDir(".")
	if err != nil {
		return ""
	}

	versionPattern := regexp.MustCompile(`^\d+\.\d+\.\d+(?:-(?:rc|alpha|beta)(?:\d+)?)?$`)
	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && versionPattern.MatchString(entry.Name()) {
			versions = append(versions, entry.Name())
		}
	}

	if len(versions) == 0 {
		return ""
	}

	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	return versions[0]
}

func createVersionList(w fyne.Window, stopButton *widget.Button, downloadGroup, versionContainer *fyne.Container) *widget.List {
	versions := getAllInstalledVersions()

	versionList = widget.NewList(
		func() int { return len(versions) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			launchButton := widget.NewButton("Launch", nil)
			removeButton := widget.NewButton("Remove", nil)
			launchButton.Resize(fyne.NewSize(80, 25))
			removeButton.Resize(fyne.NewSize(80, 25))
			launchButton.Importance = widget.HighImportance // Make the launch button important
			return container.NewBorder(nil, nil, nil, container.NewHBox(launchButton, removeButton), label)
		},
		func(index widget.ListItemID, item fyne.CanvasObject) {
			cont := item.(*fyne.Container)
			label := cont.Objects[0].(*widget.Label)
			buttons := cont.Objects[1].(*fyne.Container)
			launchButton := buttons.Objects[0].(*widget.Button)
			removeButton := buttons.Objects[1].(*widget.Button)

			version := versions[index]
			label.SetText(version)
			launchButton.SetText("Launch")
			launchButton.OnTapped = func() {
				if currentProcess != nil {
					dialog.ShowError(fmt.Errorf("OWLCMS is already running"), w)
					return
				}

				fmt.Printf("Launching version %s\n", version)
				if err := checkJava(); err != nil {
					dialog.ShowError(fmt.Errorf("java check/installation failed: %w", err), w)
					return
				}

				if err := launchOwlcms(version, launchButton, stopButton, downloadGroup, versionContainer); err != nil {
					dialog.ShowError(err, w)
					return
				}
			}

			removeButton.SetText("Remove")
			removeButton.OnTapped = func() {
				dialog.ShowConfirm("Confirm Remove",
					fmt.Sprintf("Do you want to remove OWLCMS version %s?", version),
					func(ok bool) {
						if !ok {
							return
						}

						err := os.RemoveAll(version)
						if err != nil {
							dialog.ShowError(fmt.Errorf("failed to remove OWLCMS %s: %w", version, err), w)
							return
						}

						versions = getAllInstalledVersions()
						versionList.Length = func() int { return len(versions) }
						versionList.Refresh()
					},
					w)
			}
		},
	)

	versionList.OnSelected = func(id widget.ListItemID) {
		if id < len(versions) {
			fmt.Printf("Selected version: %s\n", versions[id])
		}
	}

	if len(versions) > 0 {
		versionList.Select(0)
	}

	if latest := findLatestInstalled(); latest != "" {
		for i, v := range versions {
			if v == latest {
				versionList.Select(i)
				break
			}
		}
	}

	// Log the versions being added
	fmt.Println("Versions being added to the version list:")
	for _, version := range versions {
		fmt.Println(version)
	}

	return versionList
}