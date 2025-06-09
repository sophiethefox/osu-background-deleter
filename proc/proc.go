package proc

import (
	"os/exec"
	"strings"
	"syscall"
	// "syscall"
)

// convert to using mem.FindProcess?
func CheckIfProcessRunning(name string) bool {
	cmd := exec.Command("powershell", "-Command", "Get-WmiObject Win32_Process | select commandline | Select-String -Pattern \""+name+"\"")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW

	output, _ := cmd.CombinedOutput()
	str := string(output)

	occurences := 0
	lines := strings.Split(str, "\n")

	for line := range lines {
		if strings.Contains(lines[line], ".exe") {
			occurences += 1
		}
	}

	return (occurences > 1)
}
