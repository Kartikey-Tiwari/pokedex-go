package repl

import "os/exec"

func enableRawMode() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func restoreNormalTTYSettings() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
}
