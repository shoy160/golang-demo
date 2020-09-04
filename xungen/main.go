package main

import (
	"os"

	"github.com/kardianos/service"
)

var logger = service.ConsoleLogger

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	start()
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "xungen_reader",       //服务显示名称
		DisplayName: "xungen reader",       //服务名称
		Description: "巡更棒数据读取服务 - Go语言版本.", //服务描述
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Error(err)
	}

	if err != nil {
		logger.Error(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			s.Install()
			logger.Info("服务[xungen_reader]安装成功!")
			s.Start()
			logger.Info("服务[xungen_reader]启动成功!")
			break
		case "start":
			s.Start()
			logger.Info("服务[xungen_reader]启动成功!")
			break
		case "stop":
			s.Stop()
			logger.Info("服务[xungen_reader]关闭成功!")
			break
		case "restart":
			s.Stop()
			logger.Info("服务[xungen_reader]关闭成功!")
			s.Start()
			logger.Info("服务[xungen_reader]启动成功!")
			break
		case "remove":
			s.Stop()
			logger.Info("服务[xungen_reader]关闭成功!")
			s.Uninstall()
			logger.Info("服务[xungen_reader]卸载成功!")
			break
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
