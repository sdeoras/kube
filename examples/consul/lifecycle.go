package main

import v1 "k8s.io/api/core/v1"

func getLifeCycle() *v1.Lifecycle {
	lifecycle := new(v1.Lifecycle)
	lifecycle.PreStop = new(v1.Handler)
	lifecycle.PreStop.Exec = new(v1.ExecAction)
	lifecycle.PreStop.Exec.Command = []string{"consul", "leave"}

	return lifecycle
}
