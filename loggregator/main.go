package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func main() {

	f := os.NewFile(uintptr(syscall.Stdout), "/proc/"+os.Args[1]+"/fd/1")
	defer f.Close()

	for {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		// if err := scanner.Err(); err != nil {
		// 	log.Fatal(err)
		// }
	}

}

// FileDescriptors returns the currently open file descriptors of a process.
// func FileDescriptors() ([]uintptr, error) {
// 	names, err := p.fileDescriptors()
// 	if err != nil {
// 		return nil, err
// 	}

// 	fds := make([]uintptr, len(names))
// 	for i, n := range names {
// 		fd, err := strconv.ParseInt(n, 10, 32)
// 		if err != nil {
// 			return nil, fmt.Errorf("could not parse fd %s: %s", n, err)
// 		}
// 		fds[i] = uintptr(fd)
// 	}

// 	return fds, nil
// }

// https://github.com/prometheus/procfs/blob/master/proc.go#L201

func fileDescriptors(pid string) ([]string, error) {
	d, err := os.Open("/proc/" + pid + "/fd")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %s", d.Name(), err)
	}

	return names, nil
}
func FileDescriptors(pid string) ([]uintptr, error) {

	channels := []string{"1"}

	fds := make([]uintptr, len(channels))
	for i, n := range channels {
		fd, err := strconv.ParseInt(n, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("could not parse fd %s: %s", n, err)
		}
		fds[i] = uintptr(fd)
	}

	return fds, nil
}
