package actions

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/openshift/api/helm/v1beta1"
	"github.com/openshift/console/pkg/helm/chartproxy"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"

	"k8s.io/client-go/dynamic"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

// constants
const (
	configNamespace         = "openshift-config"
	tlsSecretCertKey        = "tls.crt"
	tlsSecretKey            = "tls.key"
	caBundleKey             = "ca-bundle.crt"
	tlsSecretPattern        = "tlscrt-*"
	tlsKeyPattern           = "tlskey-*"
	cacertPattern           = "cacert-*"
	openshiftRepoUrl        = "https://charts.openshift.io"
	chartRepoPrefix         = "chart.openshift.io/chart-url-prefix"
	openshiftChartUrlPrefix = "https://github.com/openshift-helm-charts/"

	baUsernameKey = "username"
	baPasswordKey = "password"
)

// writeTempFile creates a temporary file with the given `data`. `pattern`
// is used by `os.CreateTemp` to create a file in the filesystem.
func writeTempFile(data []byte, pattern string) (*os.File, error) {
	f, createErr := os.CreateTemp("", pattern)
	if createErr != nil {
		return nil, createErr
	}

	_, writeErr := f.Write(data)
	if writeErr != nil {
		return nil, writeErr
	}

	closeErr := f.Close()
	if closeErr != nil {
		return nil, closeErr
	}

	return f, nil
}

type ChartInfo struct {
	Name                string
	Version             string
	RepositoryName      string
	RepositoryNamespace string
}

// getChartInfoFromChartUrl returns information for the chart contained in
// the given `url`.
//
// This function works by listing all available Helm Chart repositories (either
// scoped by the given `namespace` or cluster scoped), then comparing URLs of
// all existing charts in the repository manifest to match the given `chartUrl`.
func getChartInfoFromChartUrl(
	chartUrl string,
	namespace string,
	client dynamic.Interface,
	coreClient corev1client.CoreV1Interface,
) (*ChartInfo, error) {
	repositories, err := chartproxy.NewRepoGetter(client, coreClient).List(namespace)
	if err != nil {
		return nil, fmt.Errorf("error listing repositories: %w", err)
	}

	for _, repository := range repositories {
		idx, err := repository.IndexFile()
		if err != nil {
			return nil, fmt.Errorf("error producing the index file of repository %q in namespace %q is %q", repository.Name, repository.Namespace, err.Error())
		}
		for chartIndex, chartVersions := range idx.Entries {
			for _, chartVersion := range chartVersions {
				for _, url := range chartVersion.URLs {
					if chartUrl == url {
						return &ChartInfo{
							RepositoryName:      repository.Name,
							RepositoryNamespace: repository.Namespace,
							Name:                chartIndex,
							Version:             chartVersion.Version,
						}, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("could not find a repository for the chart url %q in namespace %q", chartUrl, namespace)
}

func findChartVersion(chartValue string) string {
	for i := 1; i < len(chartValue); i++ {
		if unicode.IsDigit(rune(chartValue[i])) && chartValue[i-1] == byte('-') {
			start := i
			var end int
			for j := i; j < len(chartValue)-1; j++ {
				if unicode.IsDigit(rune(chartValue[j])) == false && unicode.IsDigit(rune(chartValue[j+1])) == false {
					end = j - 2
				}
			}
			return chartValue[start:end]
		}
	}
	return ""
}

func getChartInfoFromIndexEntry(
	indexEntry, namespace, url string) *ChartInfo {
	indexArr := strings.Split(indexEntry, "--")
	pathsOfUrl := strings.Split(url, "/")
	version := findChartVersion(pathsOfUrl[len(pathsOfUrl)-1])
	return &ChartInfo{
		RepositoryName:      indexArr[1],
		RepositoryNamespace: namespace,
		Name:                indexArr[0],
		Version:             version,
	}
}

func GetChart(url string, conf *action.Configuration, repositoryNamespace string, client dynamic.Interface, coreClient corev1client.CoreV1Interface, filesCleanup bool, indexEntry string) (*chart.Chart, error) {
	var err error
	var chartInfo *ChartInfo
	tlsFiles := []*os.File{}
	cmd := action.NewInstall(conf)
	if repositoryNamespace == "" {
		chartLocation, err := cmd.ChartPathOptions.LocateChart(url, settings)
		if err != nil {
			return nil, err
		}
		return loader.Load(chartLocation)
	}

	chartInfo = getChartInfoFromIndexEntry(indexEntry, repositoryNamespace, url)

	cmd.ChartPathOptions.Version = chartInfo.Version

	connectionConfig, isClusterScoped, err := getRepositoryConnectionConfig(chartInfo.RepositoryName, chartInfo.RepositoryNamespace, client)
	if err != nil {
		return nil, err
	}
	if isClusterScoped {
		clusterConnectionConfig := connectionConfig.(v1beta1.ConnectionConfig)
		cmd.ChartPathOptions.RepoURL = clusterConnectionConfig.URL
		tlsFiles, err = setUpAuthentication(&cmd.ChartPathOptions, clusterConnectionConfig, coreClient)
		if err != nil {
			return nil, fmt.Errorf("error setting up authentication: %w", err)
		}
	} else {
		cmd.ChartPathOptions.RepoURL = connectionConfig.(v1beta1.ConnectionConfigNamespaceScoped).URL
		tlsFiles, err = setUpAuthenticationProject(&cmd.ChartPathOptions, connectionConfig.(v1beta1.ConnectionConfigNamespaceScoped), coreClient, repositoryNamespace)
		if err != nil {
			return nil, fmt.Errorf("error setting up authentication: %w", err)
		}
	}
	chartLocation, locateChartErr := cmd.ChartPathOptions.LocateChart(chartInfo.Name, settings)
	if locateChartErr != nil {
		return nil, fmt.Errorf("error locating chart: %w", locateChartErr)
	}
	defer func() {
		if filesCleanup == false {
			return
		}
		for _, f := range tlsFiles {
			os.Remove(f.Name())
		}
	}()
	return loader.Load(chartLocation)
}
