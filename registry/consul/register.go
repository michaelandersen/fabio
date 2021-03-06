package consul

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/eBay/fabio/_third_party/github.com/hashicorp/consul/api"
	"github.com/eBay/fabio/config"
)

func serviceRegistration(addr, name string, interval, timeout time.Duration) (*api.AgentServiceRegistration, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	ipstr, portstr, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		return nil, err
	}

	ip, err := config.LocalIP()
	if err != nil {
		return nil, err
	}
	if ip == nil {
		givenip := net.ParseIP(ipstr)
                if givenip == nil {
                        return nil, errors.New("no local ip")
                }
                ip = givenip
	}

	serviceID := fmt.Sprintf("%s-%s-%d", name, hostname, port)

	checkURL := fmt.Sprintf("http://%s:%d/health", ip, port)
	if ip.To16() != nil {
		checkURL = fmt.Sprintf("http://[%s]:%d/health", ip, port)
	}

	service := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    name,
		Address: ip.String(),
		Port:    port,
		Check: &api.AgentServiceCheck{
			HTTP:     checkURL,
			Interval: interval.String(),
			Timeout:  timeout.String(),
		},
	}

	return service, nil
}
