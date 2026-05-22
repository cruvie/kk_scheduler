package kk_scheduler

func (x *PBRegisterService) Check() error {
	if x.GetTarget() == "" {
		return ErrTargetEmpty
	}
	if x.GetServiceName() == "" {
		return ErrServiceNameEmpty
	}
	if x.GetAuthToken() == "" {
		return ErrServiceAuthTokenEmpty
	}
	return nil
}

func (x *PBRegisterJob) Check() error {
	if x.GetServiceName() == "" {
		return ErrServiceNameEmpty
	}
	if x.GetFuncName() == "" {
		return ErrFuncNameEmpty
	}
	return nil
}
