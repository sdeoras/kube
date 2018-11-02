package main

const (
	consulImage      = "consul:1.2.2"
	serverReplicas   = 3
	datacenter       = "gke"
	kvKey            = "bar/baz"
	kvPrefix         = "foo"
	consulLoadApp    = "consul-load"
	consulWatchApp   = "consul-watch"
	consulServerName = "consul"
	consulService    = consulServerName
)
