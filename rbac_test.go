package cmd

import (
	"testing"
)

func TestCurrentContext(t *testing.T) {
	if output, err := currentContext(); output != "docker-desktop" {
		t.Errorf("out: %s, err %s", output, err)
	}
}

func TestHelm(t *testing.T) {
	kc := kubeconfig{}
	kc.setup()
	if err := freshDocClu(); err != nil {
		t.Fatal(err)
	}
	if err := makeUser(&kc, "client_rbac.yaml"); err != nil {
		t.Error(err)
	}
	if err := tryHelm(); err != nil {
		t.Error(err)
	}
	if err := delUser(&kc); err != nil {
		t.Error(err)
	}
}

func TestSmoke(t *testing.T) {
	kc := kubeconfig{}
	kc.setup()
	kc.cleanup()
	if err := freshDocClu(); err != nil {
		t.Fatal(err)
	}
	if err := makeUser(&kc, "client_rbac.yaml"); err != nil {
		t.Error(err)
	}
}

func TestMain(t *testing.T) {
	kc := kubeconfig{}
	kc.setup()
	if err := freshDocClu(); err != nil {
		t.Fatal(err)
	}
	if err := makeUser(&kc, "client_rbac.yaml"); err != nil {
		t.Error(err)
	}
	if err := tryConnect(); err != nil {
		t.Error(err)
	}
	if err := delUser(&kc); err != nil {
		t.Error(err)
	}
}

func TestConnect(t *testing.T) {
	if err := tryConnect(); err != nil {
		t.Error(err)
	}
}
