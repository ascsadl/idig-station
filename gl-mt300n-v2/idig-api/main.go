package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type Out struct {
	Env []string `json:"env"`
	In  string   `json:"in"`
}

type StationResponse struct {
	Enabled bool     `json:"enabled"`
	Name    string   `json:"name"`
	Log     []string `json:"log"`
	Code    int      `json:"code"`
}

func getStation() *StationResponse {
	_, err := os.Stat("/var/run/idig-station.pid")
	enabled := err == nil

	name := "iDig Station"
	data, err := os.ReadFile("/etc/idig-station")
	if err == nil {
		name = strings.TrimSpace(string(data))
	}

	logEntries := make([]string, 0)
	stdout, err := exec.Command("/sbin/logread", "-e", "idig-station").CombinedOutput()
	if err == nil {
		for _, line := range strings.Split(string(stdout), "\n") {
			_, msg, _ := strings.Cut(strings.TrimSpace(line), ": ")
			if msg != "" {
				logEntries = append(logEntries, msg)
			}
		}
	} else {
		logEntries = append(logEntries, err.Error())
	}

	return &StationResponse{
		Enabled: enabled,
		Name:    name,
		Log:     logEntries,
		Code:    0,
	}
}

func setStation() *StationResponse {
	resp := getStation()

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		resp.Code = 1
		return resp
	}
	v, err := url.ParseQuery(string(data))
	if err != nil {
		resp.Code = 2
		return resp
	}

	name := v.Get("name")
	if name == "" {
		name = "iDig Station"
	}
	if os.WriteFile("/etc/idig-station", []byte(name), 0o644) != nil {
		resp.Code = 3
		return resp
	}
	resp.Name = name

	if v.Get("enabled") == "true" {
		if exec.Command("/etc/init.d/idig-station", "enable").Run() != nil {
			resp.Code = 4
			return resp
		}
		var cmd string
		if resp.Enabled {
			cmd = "restart"
		} else {
			cmd = "start"
		}
		if exec.Command("/etc/init.d/idig-station", cmd).Run() != nil {
			resp.Code = 5
			return resp
		}
		resp.Enabled = true
	} else {
		if exec.Command("/etc/init.d/idig-station", "disable").Run() != nil {
			resp.Code = 6
			return resp
		}
		if exec.Command("/etc/init.d/idig-station", "stop").Run() != nil {
			resp.Code = 7
			return resp
		}
		resp.Enabled = false
	}

	resp.Code = 0
	return resp
}

func main() {
	var resp *StationResponse

	switch os.Getenv("PATH_INFO") {
	case "/station/get":
		resp = getStation()
	case "/station/set":
		resp = setStation()
	}

	fmt.Println("Content-Type: application/json")
	fmt.Println()

	if resp == nil {
		fmt.Print(`{"code": -1}`)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Print(`{"code": -2}`)
	} else {
		fmt.Print(string(data))
	}
}
