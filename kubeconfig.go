package cmd

type kubeconfig struct {
	userContext string
	testContext string
}

func (kc *kubeconfig) setup() *kubeconfig {
	out, _ := currentContext()
	kc.userContext = out
	kc.testContext = "telepresence-test-developer"
	return kc
}

func (kc *kubeconfig) cleanup() error {
	_, err := Output("kubectl", "config", "use-context", kc.userContext)
	if err != nil {
		return err
	}
	_, err = Output("kubectl", "config", "delete-context", kc.testContext)
	return err
}
