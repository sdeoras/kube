package main

import v1 "k8s.io/api/core/v1"

func getEnv() []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name: "GOSSIP_ENCRYPTION_KEY",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "consul",
					},
					Key: "gossip-encryption-key",
				},
			},
		},
		{
			Name: "NAMESPACE",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		},
		{
			Name: "POD_IP",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
	}
}
