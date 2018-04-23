package main

import "testing"

func TestBuildName(t *testing.T) {
	name := buildName("deployment", "hello")
	if name != "hello-deploy.yaml" {
		t.Fatalf("name invalid %s", name)
	}
}
