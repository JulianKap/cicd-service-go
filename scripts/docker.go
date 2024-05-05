package scripts

import (
	"bufio"
	"cicd-service-go/pipeline"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type rewriter struct {
	ctx    context.Context
	prefix string
	input  *bufio.Reader
}

func (rw *rewriter) watch() {
	for {
		select {
		case <-rw.ctx.Done():
			return
		case err := <-rw.rewriteInput():
			if err != nil {
				rw.writeToOutput(fmt.Sprintf("Error while reading command output: %v", err))
				return
			}
		}
	}
}

func (rw *rewriter) writeToOutput(line string) {
	fmt.Printf("%s[%s]%s %s", rw.prefix, line)
}

func (rw *rewriter) rewriteInput() <-chan error {
	e := make(chan error)

	go func() {
		line, err := rw.input.ReadString('\n')
		if err == nil || err == io.EOF {
			rw.writeToOutput(line)
			e <- nil

			return
		}

		e <- err
	}()

	return e
}

func newRewriter(ctx context.Context, prefix string) io.Writer {
	pr, pw, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	w := &rewriter{
		ctx:    ctx,
		prefix: prefix,
		input:  bufio.NewReader(pr),
	}

	go w.watch()

	return pw
}

// PullImage получаем образ
func PullImage(wg *sync.WaitGroup, name string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Done()
	}()

	output := newRewriter(ctx, name)

	cmd := exec.Command("docker", "pull", name)
	cmd.Stdout = output
	cmd.Stderr = output

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func RunCommandImage(wg *sync.WaitGroup, s pipeline.Step) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Done()
	}()

	output := newRewriter(ctx, s.Image)

	//cmd := exec.Command("docker", "pull", name)

	// Создание команды для запуска Docker образа с командами шага
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", s.Image, "/bin/bash", "-c", fmt.Sprintf("%s", s.Commands[0]))

	cmd.Stdout = output
	cmd.Stderr = output

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
