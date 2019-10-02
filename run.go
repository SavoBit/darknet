package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	path := os.Args[1]
	outPath := os.Args[2]

	c := exec.Command("./darknet", "detect", "backup/yolo.cfg", "backup/yolo.backup", "-thresh", "0.3")
	stdin, err := c.StdinPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := c.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		r := bufio.NewReader(stderr)
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				panic(err)
			}
			line = strings.TrimSpace(line)
			fmt.Println(line)
		}
	}()
	r := bufio.NewReader(stdout)
	if err := c.Start(); err != nil {
		panic(err)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for i, fi := range files {
		fmt.Printf("%d/%d\n", i, len(files))
		if !strings.HasSuffix(fi.Name(), ".jpg") && !strings.HasSuffix(fi.Name(), ".png") {
			continue
		}
		inputFname := path + "/" + fi.Name()
		label := strings.Split(fi.Name(), ".")[0]
		predFname := outPath + "/" + label + ".jpg"
		txtFname := outPath + "/" + label + ".txt"
		var lines []string
		for {
			line, err := r.ReadString(':')
			if err != nil {
				panic(err)
			}
			fmt.Println(strings.TrimSpace(line))
			if strings.Contains(line, "Enter") {
				lines = append(lines, line)
				break
			}
			lines = append(lines, line)
		}
		if i > 1 {
			exec.Command("mv", "predictions.jpg", predFname).Run()
			metaBytes := []byte(strings.Join(lines, ""))
			if err := ioutil.WriteFile(txtFname, metaBytes, 0644); err != nil {
				panic(err)
			}
		}
		stdin.Write([]byte(inputFname + "\n"))
	}
}
