package cmd

import (
	"encoding/base64"
)

func makeUser(kc *kubeconfig, path string) error {
	_, err := output("kubectl", "apply", "-f", path)
	if err != nil {
		return err
	}

	secret, err := output("kubectl", "get", "sa", kc.testContext, "-o", "jsonpath={.secrets[0].name}")
	if err != nil {
		return err
	}
	encSecret, err := output("kubectl", "get", "secret", secret, "-o", "jsonpath={.data.token}")
	if err != nil {
		return err
	}
	token, err := base64.StdEncoding.DecodeString(encSecret)
	if err != nil {
		return err
	}
	_, err = output("kubectl", "config", "set-credentials", kc.testContext, "--token", string(token))
	if err != nil {
		return err
	}

	_, err = output("kubectl", "config", "set-context", kc.testContext, "--user", kc.testContext, "--cluster", kc.userContext)
	if err != nil {
		return err
	}

	_, err = output("kubectl", "config", "use-context", kc.testContext)
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

	_, err = output("kubectl", "delete", "ClusterRole", "telepresence-role")
	if err != nil {
		return err
	}
	_, err = output("kubectl", "delete", "ServiceAccount", kc.testContext)
	if err != nil {
		return err
	}
	_, err = output("kubectl", "delete", "ClusterRoleBinding", "telepresence-clusterrolebinding")
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

func tryHelm() error {
	_, err := output("helm", "install", "traffic-manager", "-n", "ambassador", "datawire/telepresence")
	if err != nil {
		return err
	}
	_, err = output("helm", "uninstall", "traffic-manager", "-n", "ambassador")
	if err != nil {
		return err
	}
	return nil
}
