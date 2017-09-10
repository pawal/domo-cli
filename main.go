package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/galdor/go-cmdline"
	domoto "github.com/pawal/go-domoto"
)

var verbose bool
var c *domoto.Config

func main() {
	cmdline := cmdline.New()
	cmdline.AddFlag("v", "verbose", "log more information")
	cmdline.AddOption("", "url", "url", "Domoticz API URL")
	cmdline.SetOptionDefault("url", "http://localhost:8080")
	cmdline.AddOption("p", "password", "password", "Domoticz User password")
	cmdline.AddOption("u", "user", "user", "Domoticz API username")
	cmdline.AddCommand("device", "info on device <id>")
	cmdline.AddCommand("device-toggle", "toggle device <id>")
	cmdline.AddCommand("device-on", "switch device on <id>")
	cmdline.AddCommand("device-off", "switch device off <id>")
	cmdline.AddCommand("scene-run", "execute scene <id>")
	cmdline.AddCommand("group-on", "switch group on <id>")
	cmdline.AddCommand("group-off", "switch group off <id>")
	cmdline.AddCommand("list-devices", "list all devices")
	cmdline.AddCommand("list-scenes", "list all scenes and groups")
	cmdline.AddCommand("scene-info", "list devices in scene/group")
	cmdline.Parse(os.Args)

	verbose = cmdline.IsOptionSet("verbose")

	var cmdFunc func([]string)
	switch cmdline.CommandName() {
	case "device":
		cmdFunc = CmdDevice
	case "device-toggle":
		cmdFunc = CmdToogleDevice
	case "device-on":
		cmdFunc = CmdDeviceOn
	case "device-off":
		cmdFunc = CmdDeviceOff
	case "list-devices":
		cmdFunc = CmdListDevices
	case "list-scenes":
		cmdFunc = CmdListScenes
	case "scene-info":
		cmdFunc = CmdScene
	case "scene-run":
		cmdFunc = CmdSceneRun
	case "group-on":
		cmdFunc = CmdSceneRun //same as "on" for a scene
	case "group-off":
		cmdFunc = CmdGroupOff
	default:
		cmdline.PrintUsage(os.Stdout)
	}

	c = domoto.New(cmdline.OptionValue("url"),
		cmdline.OptionValue("user"), cmdline.OptionValue("password"))

	// execute command
	cmdFunc(cmdline.CommandArgumentsValues())
}

// CmdDevice shows info on a device
func CmdDevice(args []string) {
	if len(args) == 0 {
		fmt.Println("Missing device id")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if verbose {
		fmt.Println("Info on device id " + args[0])
	}
	resp, err := c.Device(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// pretty print JSON
	j, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(j))
}

// CmdListDevices calls the domoto AllDevices
func CmdListDevices(args []string) {
	var filter string
	if len(args) != 0 {
		filter = args[0]
	} else {
		filter = ""
	}
	fmt.Println(filter)
	res, err := c.AllDevices(filter)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, d := range res.Devices {
		if verbose {
			fmt.Printf("%s: %s (%s/%s)\n", d.Idx, d.Name, d.SubType, d.SwitchType)
		} else {
			fmt.Printf("%s: %s\n", d.Idx, d.Name)
		}
	}
	if verbose {
		fmt.Println(res.Status)
	}
}

// CmdToogleDevice toggles a domoto device
func CmdToogleDevice(args []string) {
	if len(args) == 0 {
		fmt.Println("Missing device id")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if verbose {
		fmt.Println("Toggle device id " + args[0])
	}
	res, err := c.DeviceToggle(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if verbose {
		fmt.Println(res.Status)
	}
}

func deviceRun(args []string, cmd string) {
	if len(args) == 0 {
		fmt.Println("Missing device id")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if verbose {
		fmt.Printf("Command \"%s\"for device id %s", cmd, args[0])
	}
	res, err := c.DeviceSwitch(id, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if verbose {
		fmt.Println(res.Status)
	}
}

// CmdDeviceOn sets a device to "On"
func CmdDeviceOn(args []string) {
	deviceRun(args, "On")
}

// CmdDeviceOff sets a device to "Off"
func CmdDeviceOff(args []string) {
	deviceRun(args, "Off")
}

// CmdScene lists devices in a scene/group
func CmdScene(args []string) {
	if len(args) == 0 {
		fmt.Println("Missing scene/group id")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if verbose {
		fmt.Println("Info on device id " + args[0])
	}
	resp, err := c.SceneDevices(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, d := range resp.Result {
		if verbose {
			fmt.Printf("%s: %s (%s/%s)\n", d.DevRealIdx, d.Name, d.Type, d.SubType)
		} else {
			fmt.Printf("%s: %s\n", d.DevRealIdx, d.Name)
		}
	}
}

// CmdListScenes calls the domoto AllScenes
func CmdListScenes(args []string) {
	res, err := c.AllScenes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, s := range res.Scene {
		if verbose {
			fmt.Printf("%s: %s (%s) - %s\n", s.Idx, s.Name, s.Type, s.Status)
		} else {
			fmt.Printf("%s: %s\n", s.Idx, s.Name)
		}
	}
	if verbose {
		fmt.Println(res.Status)
	}
}

// sceneRun triggers a scene or group, used by other commands
func sceneRun(args []string, cmd string) {
	if len(args) == 0 {
		fmt.Println("Missing scene/group id")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, err := c.SceneSwitch(id, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if verbose {
		fmt.Println(res.Status)
	}
}

// CmdSceneRun executes a scene or turns a group "on"
func CmdSceneRun(args []string) {
	sceneRun(args, "On")
}

// CmdGroupOff sets a group to "off"
func CmdGroupOff(args []string) {
	sceneRun(args, "Off")
}
