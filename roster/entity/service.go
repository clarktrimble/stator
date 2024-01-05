package entity

import (
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
)

type Service struct {
	Id          string
	Name        string
	Tags        []string
	IpAddress   string
	Port        int
	MonitorSpec string
}

func (svc *Service) NameId() string {

	return fmt.Sprintf("%s-%s", svc.Name, svc.Id)
}

func (svc *Service) Valid() (err error) {

	errs := []string{}

	if svc.Id == "" {
		errs = append(errs, "Id must not be blank")
	}

	if svc.Name == "" {
		errs = append(errs, "Name must not be blank")
	}

	parsed := net.ParseIP(svc.IpAddress)
	if parsed == nil {
		errs = append(errs, "IpAddress failed to parse")
	}

	if svc.Port < 1 || svc.Port > 65535 {
		errs = append(errs, "Port must be between 1 and 65535")
	}

	if svc.MonitorSpec == "" {
		errs = append(errs, "MonitorSpec must not be blank")
	}

	if len(errs) != 0 {
		err = errors.Errorf("invalid Service: %s", strings.Join(errs, ","))
	}
	return
}
