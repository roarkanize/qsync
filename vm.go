package qsync

import (
	"bytes"
	"os/exec"
	"strings"
)

type VMState uint8
type VMClass uint8

const (
	VMHalted VMState = iota
	VMTransient
	VMRunning
)

const (
	NilVM VMClass = iota
	AdminVM
	AppVM
	TemplateVM
	StandaloneVM
	DispVM
)

type VM struct {
	Name  []byte
	Class VMClass
	State VMState
}

func ListVMs() ([]VM, error) {
	cmd := exec.Command("qrexec-client-vm", "$adminvm", "admin.vm.List")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	vms := []VM(nil)

	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}

		field := strings.Fields(line)

		vms = append(vms, VM{
			Name: []byte(field[0]),
			Class: (func() VMClass {
				class := strings.Split(field[1], "=")[1]
				switch class {
				case "AdminVM":
					return AdminVM
				case "AppVM":
					return AppVM
				case "TemplateVM":
					return TemplateVM
				case "StandaloneVM":
					return StandaloneVM
				case "DispVM":
					return DispVM
				}

				return NilVM
			})(),
			State: (func() VMState {
				state := strings.Split(field[2], "=")[1]
				switch state {
				case "Halted":
					return VMHalted
				case "Transient":
					return VMTransient
				case "Running":
					return VMRunning
				}

				return VMHalted
			})(),
		})
	}

	return vms, nil
}

func VMExists(name []byte) (bool, error) {
	vms, err := ListVMs()
	if err != nil {
		return false, err
	}

	for _, vm := range vms {
		if bytes.Compare(name, vm.Name) == 0 {
			return true, nil
		}
	}

	return false, nil
}
