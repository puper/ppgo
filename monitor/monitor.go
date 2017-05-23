package monitor

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/puper/ppgo/listener"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	reloadSignal    chan os.Signal
	stopSignal      chan os.Signal
	cmdSignal       chan *exec.Cmd
	currentCmd      *exec.Cmd
	lastCmdSignalAt = time.Now()
)

func Run() {
	reloadSignal = make(chan os.Signal, 1)
	stopSignal = make(chan os.Signal, 1)
	cmdSignal = make(chan *exec.Cmd, 1)
	var err error
	currentCmd, err = startProcess()
	if err != nil {
		log.Errorf("monitor start process error: %v", err)
		return
	}
	signal.Notify(stopSignal,
		os.Kill,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	signal.Notify(reloadSignal, syscall.SIGHUP)
	handleSignal()

}

func handleSignal() {
	var err error
	for {
		select {
		case <-reloadSignal:
			oldCmd := currentCmd
			for {
				log.Infof("monitor reload serve...")
				currentCmd, err = startProcess()
				if err == nil {
					log.Infof("monitor reloaded serve...")
					stopProcess(oldCmd)
					log.Infof("monitor stoped old serve...")
					break
				}
				log.Infof("monitor reload serve error...")
				time.Sleep(5 * time.Second)
			}
		case <-stopSignal:
			log.Infof("monitor stop serve")
			stopProcess(currentCmd)
			return
		case cmd := <-cmdSignal:
			if cmd == currentCmd {
				log.Infof("monitor current serve killed...")
				remain := time.Now().Sub(lastCmdSignalAt)
				if remain < 5*time.Second {
					time.Sleep(5*time.Second - remain)
				}
				lastCmdSignalAt = time.Now()
				for {
					currentCmd, err = startProcess()
					if err == nil {
						break
					}
				}
			}
		}
	}
}

func stopProcess(cmd *exec.Cmd) {
	for {
		cmd.Process.Signal(syscall.SIGTERM)
		fromCmd := <-cmdSignal
		if fromCmd == cmd {
			log.Println("monitor child serve exited...")
			return
		}
		time.Sleep(time.Second)
	}
}

func startProcess() (*exec.Cmd, error) {
	var err error
	args := os.Args[1:]
	args[0] = "serve"
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.ExtraFiles, err = listener.GetFiles([]string{viper.GetString("server.Addr")})
	if err != nil {
		log.Infof("monitor files error: %v", err)
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	go func(cmd *exec.Cmd) {
		err := cmd.Wait()
		if err != nil {
			log.Infof("monitor wait serve error:", err)
		}
		cmdSignal <- cmd
	}(cmd)
	return cmd, nil
}
