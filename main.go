package main

import (
	"encoding/json"
	"github.com/getlantern/systray"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	systray.Run(onReady, onExit)
}

type MenuItem struct {
	Title       string     `json:"title"`
	Dir         string     `json:"dir"`
	Description string     `json:"description"`
	Command     string     `json:"command"`
	Children    []MenuItem `json:"children"`
}

func getIconData() []byte {
	data, err := os.ReadFile("ra.ico")
	if err != nil {
		log.Fatal("failed to read icon")
	}
	return data
}

func onReady() {
	iconData, err := Asset("ra.ico")
	if err != nil {
		log.Fatalf("failed to get icon data: %v", err)
	}
	systray.SetIcon(iconData)
	systray.SetTitle("RA 커맨드")

	menuItems := loadMenu("menu.json")
	for _, item := range menuItems {
		addMenu(nil, item)
	}

	addMenu(nil, MenuItem{
		Title:       "Exit",
		Description: "나가기",
		Command:     "Exit",
	})
}

func loadMenu(filename string) []MenuItem {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read menu file: %v", err)
	}

	var menuItems []MenuItem
	err = json.Unmarshal(data, &menuItems)
	if err != nil {
		log.Fatalf("Failed to parse menu file: %v", err)
	}
	return menuItems
}

func addMenu(parent *systray.MenuItem, item MenuItem) *systray.MenuItem {
	var menuItem *systray.MenuItem
	if parent == nil {
		menuItem = systray.AddMenuItem(item.Title, item.Description)
	} else {
		menuItem = parent.AddSubMenuItem(item.Title, item.Description)
	}

	if item.Command != "" {
		go func(command string, dir string) {
			for {
				<-menuItem.ClickedCh
				if command == "Exit" {
					systray.Quit()
				} else {
					//cmd := exec.Command("cmd.exe", "/C", command, "&&", "exit")
					cmd := exec.Command("cmd.exe", "/C", command)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.SysProcAttr = &syscall.SysProcAttr{
						CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
					}
					cmd.Dir = dir
					err := cmd.Run()
					if err != nil {
						log.Printf("executing command error: %v", err)
					}
				}
			}
		}(item.Command, item.Dir)
	}
	for _, child := range item.Children {
		addMenu(menuItem, child)
	}
	return menuItem
}

func onExit() {

}
