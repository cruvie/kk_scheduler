package kk_scheduler

import "errors"

var ErrJobNotFount = errors.New("job not found")

var ErrServiceNotFount = errors.New("service not found")

var ErrServiceHasJob = errors.New("service has job")

var ErrSpecIsEmpty = errors.New("spec is empty")

var ErrServiceNameEmpty = errors.New("serviceName empty")

var ErrServiceAuthTokenEmpty = errors.New("serviceAuthToken empty")

var ErrFuncNameEmpty = errors.New("funcName empty")

var ErrTargetEmpty = errors.New("target empty")

var ErrStopTask = errors.New("stop task")
