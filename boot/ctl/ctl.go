package ctl

import (
	"errors"
	"log"
)

func Register(name string, fn func() I) error {
	if _, ok := serviceMap[name]; ok {
		return errors.New("service exists")
	}
	createServices = append(createServices, Created{
		Name: name,
		Func: fn,
	})
	serviceMap[name] = struct{}{}
	return nil
}

func Destroy() (err error) {
	for i := len(destroyServices) - 1; i >= 0; i-- {
		if err = destroyServices[i].Func(); err != nil {
			log.Printf("[Boot(Destroy)] service %s destroy fail, error: %s", destroyServices[i].Name, err)
			return
		}
	}
	log.Println("[Boot] >>>>>>>>>> all service are unmount <<<<<<<<<<")
	return nil
}

func Startup() error {
	size := len(createServices)
	if size == 0 {
		return errors.New("no services require register")
	}
	services := make([]serviced, size)
	for i := 0; i < size; i++ {
		services[i] = serviced{
			id:        uint(i),
			name:      createServices[i].Name,
			implement: createServices[i].Func(),
		}
		log.Printf("[Boot] service %s has mounted, serviced id: %d", services[i].name, i)
	}
	log.Println("[Boot] >>>>>>>>>> all service are mounted <<<<<<<<<<")
	var (
		errChan = make(chan error, 1)
		err     error
	)
	for i := 0; i < size; i++ {
		if len(errChan) == 1 {
			return <-errChan
		}
		service := services[i]
		if err = service.implement.Create(); err != nil {
			log.Printf("[Boot] service %s created fail, error: %s", service.name, err.Error())
			return err
		}
		destroyServices = append(destroyServices, Destroyed{
			Name: service.name,
			Func: service.implement.Destroy,
		})
		if !service.implement.IsAsync() {
			if err = service.implement.Start(); err != nil {
				log.Printf("[Boot] service %s startup fail, error: %s", service.name, err.Error())
				return err
			}
			log.Printf("[Boot] service %s startup", service.name)
		} else {
			go func() {
				log.Printf("[Boot(Async)] service %s startup", service.name)
				if er := service.implement.Start(); err != nil {
					log.Printf("[Boot(Async)] service %s startup fail, error: %s", service.name, er.Error())
				}
			}()
		}
	}
	return nil
}
