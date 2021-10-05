package cmd

import (
	"encoding/base64"
	"errors"
	"strings"
)

func makeUser(kc *kubeconfig, path string) error {
	_, err := Output("kubectl", "apply", "-f", path)
	if err != nil {
		return err
	}

	secret, err := Output("kubectl", "get", "sa", kc.testContext, "-o", "jsonpath={.secrets[0].name}")
	if err != nil {
		return err
	}
	encSecret, err := Output("kubectl", "get", "secret", secret, "-o", "jsonpath={.data.token}")
	if err != nil {
		return err
	}
	token, err := base64.StdEncoding.DecodeString(encSecret)
	if err != nil {
		return err
	}
	_, err = Output("kubectl", "config", "set-credentials", kc.testContext, "--token", string(token))
	if err != nil {
		return err
	}

	_, err = Output("kubectl", "config", "set-context", kc.testContext, "--user", kc.testContext, "--cluster", kc.userContext)
	if err != nil {
		return err
	}

	_, err = Output("kubectl", "config", "use-context", kc.testContext)
	if err != nil {
		return err
	}

	return err
}

func delUser(kc *kubeconfig) error {
	err := kc.cleanup()
	if err != nil {
		return err
	}

	_, err = Output("kubectl", "delete", "ClusterRole", "telepresence-role")
	if err != nil {
		return err
	}
	_, err = Output("kubectl", "delete", "ServiceAccount", kc.testContext)
	if err != nil {
		return err
	}
	_, err = Output("kubectl", "delete", "ClusterRoleBinding", "telepresence-clusterrolebinding")
	if err != nil {
		return err
	}
	return nil
}

func tryConnect() error {
	err := connect()
	if err != nil {
		return err
	}

	return uninstall()
}

func smokeTest() error {
	out, err := Output("run_smoke_test.sh")
	if err != nil {
		return err
	}
	if !strings.Contains(out, "has been smoke tested and took") {
		return errors.New(out)
	}
	return nil
}

func tryHelm() error {
	_, err := Output("helm", "install", "traffic-manager", "-n", "ambassador", "datawire/telepresence")
	if err != nil {
		return err
	}
	_, err = Output("helm", "uninstall", "traffic-manager", "-n", "ambassador")
	if err != nil {
		return err
	}
	return nil
}

func Full() error {
	kc := kubeconfig{}
	kc.setup()
	if err := freshDocClu(); err != nil {
		return err
	}
	if err := makeUser(&kc, "client_rbac.yaml"); err != nil {
		return err
	}
	if err := tryConnect(); err != nil {
		return err
	}
	if err := delUser(&kc); err != nil {
		return err
	}
	return nil
}
