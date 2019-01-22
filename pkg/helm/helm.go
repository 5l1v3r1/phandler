package helm

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/helm"
)

func (client *Client) InstallHelmRelease(chstr, ns string) error {
	installOpts := []helm.InstallOption{
		helm.InstallDisableHooks(false),
		helm.InstallTimeout(30),
		// helm.ReleaseName(""),
	}
	result, err := client.InstallRelease(chstr, ns, installOpts...) //
	fmt.Print("%+v", result)
	if err != nil {
		return errors.Wrapf(err, "[ERROR] Failed to install %s", chstr)
	}
	return nil
}

func (client *Client) DeleteHelmRelease(rn string) error {
	deleteOpts := []helm.DeleteOption{
		helm.DeleteDisableHooks(false),
		helm.DeletePurge(true),
		helm.DeleteTimeout(30),
	}
	_, err := client.DeleteRelease(rn, deleteOpts...)
	if err != nil {
		return errors.Wrapf(err, "[ERROR] Failed to delete %s", rn)
	}
	return nil
}
