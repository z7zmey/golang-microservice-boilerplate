package main

import "github.com/sirupsen/logrus"

type serviceLogHook struct {
	serviceName string
}

func newServiceLogHook(serviceName string) *serviceLogHook {
	return &serviceLogHook{serviceName: serviceName}
}

func (s serviceLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (s serviceLogHook) Fire(e *logrus.Entry) error {
	e.Data["service"] = s.serviceName
	return nil
}
