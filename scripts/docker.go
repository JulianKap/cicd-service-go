package scripts

import (
	"bufio"
	"cicd-service-go/pipeline"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
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
				log.Error("watch #0: ", err)
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

// PullImage получение docker образа
func PullImage(wg *sync.WaitGroup, name string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		//wg.Done()
	}()

	output := newRewriter(ctx, name)

	cmd := exec.Command("docker", "pull", name)
	cmd.Stdout = output
	cmd.Stderr = output

	err := cmd.Run()
	if err != nil {
		log.Error("PullImage #0: ", err)
		//panic(err)
		return err
	}

	return nil
}

func RunCommandImage(wg *sync.WaitGroup, s pipeline.Step) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		//wg.Done()
	}()

	output := newRewriter(ctx, s.Image)

	// Создание команды для запуска Docker образа с командами шага
	cmd := exec.CommandContext(
		ctx,
		"docker", "run", "--rm", "-i", s.Image, "/bin/bash", "-c", fmt.Sprintf("%s", s.Commands[0]))
	cmd.Stdout = output
	cmd.Stderr = output

	err := cmd.Run()
	if err != nil {
		log.Error("PullImage #0: ", err)
		//panic(err)
		return err
	}

	return nil
}
