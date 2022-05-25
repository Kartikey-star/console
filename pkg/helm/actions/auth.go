package actions

import (
	"context"
	"fmt"
	"os"

	"github.com/openshift/api/helm/v1beta1"
	"helm.sh/helm/v3/pkg/action"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

func setUpAuthentication(chartPathOptions *action.ChartPathOptions, connectionConfig *v1beta1.ConnectionConfig, coreClient corev1client.CoreV1Interface, namespace string, isClusterScoped bool) ([]*os.File, error) {
	tlsFiles := []*os.File{}
	var secretNamespace, configMapName, configMapNameSpace, secretName string
	//set up tls cert and key
	if connectionConfig.TLSClientConfig != nil {
		secretName = connectionConfig.TLSClientConfig.Name
		secretNamespace = connectionConfig.TLSClientConfig.Namespace
		if secretNamespace == "" && isClusterScoped {
			secretNamespace = configNamespace
		} else if secretNamespace == "" && !isClusterScoped {
			secretNamespace = namespace
		}
		secret, err := coreClient.Secrets(secretNamespace).Get(context.TODO(), secretName, v1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("Failed to GET secret %s from %vreason %v", secretName, secretNamespace, err)
		}
		tlsCertBytes, found := secret.Data[tlsSecretCertKey]
		if !found {
			return nil, fmt.Errorf("Failed to find %s key in secret %s", tlsSecretCertKey, secretName)
		}
		tlsCertFile, err := writeTempFile((tlsCertBytes), tlsSecretPattern)
		if err != nil {
			return nil, err
		}
		chartPathOptions.CertFile = tlsCertFile.Name()
		tlsFiles = append(tlsFiles, tlsCertFile)
		tlsKeyBytes, found := secret.Data[tlsSecretKey]
		if !found {
			return nil, fmt.Errorf("Failed to find %s key in secret %s", tlsSecretKey, secretName)
		}
		tlsKeyFile, err := writeTempFile(tlsKeyBytes, tlsKeyPattern)
		if err != nil {
			return nil, err
		}
		chartPathOptions.KeyFile = tlsKeyFile.Name()
		tlsFiles = append(tlsFiles, tlsKeyFile)
	}
	//set up ca certificate
	if connectionConfig.CA != nil {
		configMapName = connectionConfig.CA.Name
		configMapNameSpace = connectionConfig.CA.Namespace
		if configMapNameSpace == "" && isClusterScoped {
			configMapNameSpace = configNamespace
		} else if configMapNameSpace == "" && !isClusterScoped {
			configMapNameSpace = namespace
		}
		configMap, err := coreClient.ConfigMaps(configMapNameSpace).Get(context.TODO(), configMapName, v1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("Failed to GET configmap %s, reason %v", configMapName, err)
		}
		caCertBytes, found := configMap.Data[caBundleKey]
		if !found {
			return nil, fmt.Errorf("Failed to find %s key in configmap %s", caBundleKey, configMapName)
		}
		caCertFile, caCertGetErr := writeTempFile([]byte(caCertBytes), "cacert-*")
		if caCertGetErr != nil {
			return nil, caCertGetErr
		}
		chartPathOptions.CaFile = caCertFile.Name()
		tlsFiles = append(tlsFiles, caCertFile)
	}
	return tlsFiles, nil
}
