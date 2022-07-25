package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:plural=projecthelmchartrepositories

// ProjectHelmChartRepository holds namespace-wide configuration for proxied Helm chart repository
//
// Compatibility level 2: Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=2
type ProjectHelmChartRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
	// +required
	Spec ProjectHelmChartRepositorySpec `json:"spec"`

	// Observed status of the repository within the namespace..
	// +optional
	Status HelmChartRepositoryStatus `json:"status"`
}

// Project Helm chart repository exposed within a namespace
type ProjectHelmChartRepositorySpec struct {

	// If set to true, disable the repo usage in the cluster/namespace
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// Optional associated human readable repository name, it can be used by UI for displaying purposes
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=100
	// +optional
	DisplayName string `json:"name,omitempty"`

	// Optional human readable repository description, it can be used by UI for displaying purposes
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	// +optional
	Description string `json:"description,omitempty"`

	// Required configuration for connecting to the chart repo
	ProjectConnectionConfig ConnectionConfigNamespaceScoped `json:"connectionConfig"`
}

type ConnectionConfigNamespaceScoped struct {

	// Chart repository URL
	// +kubebuilder:validation:Pattern=`^https?:\/\/`
	// +kubebuilder:validation:MaxLength=2048
	URL string `json:"url"`

	// ca is an optional reference to a config map by name containing the PEM-encoded CA bundle.
	// It is used as a trust anchor to validate the TLS certificate presented by the remote server.
	// The key "ca-bundle.crt" is used to locate the data.
	// If empty, the default system roots are used.
	// The namespace for this configmap can be provided by the user else we will search for the configmap in the namespace where the repo is instantiated.
	// +optional
	CA ConfigMapNameReference `json:"ca,omitempty"`

	// tlsClientConfig is an optional reference to a secret by name that contains the
	// PEM-encoded TLS client certificate and private key to present when connecting to the server.
	// The key "tls.crt" is used to locate the client certificate.
	// The key "tls.key" is used to locate the private key.
	// The namespace for this secret can be provided by the user else we will search for the secret in the namespace where the repo is instantiated.
	// +optional
	TLSClientConfig SecretNamespacedReference `json:"tlsClientConfig,omitempty"`

	// basicAuthConfig is an optional reference to a secret by name that contains
	// the basic authentication credentials to present when connecting to the server.
	// The key "username" is used locate the username.
	// The key "password" is used to locate the password.
	// The namespace for this secret can be provided by the user else we will search for the secret in the namespace where the repo is instantiated.
	// +optional
	BasicAuthConfig SecretNamespacedReference `json:"basicAuthConfig,omitempty"`
}

// Compatibility level 2: Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer).
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +openshift:compatibility-gen:level=2
type ProjectHelmChartRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ProjectHelmChartRepository `json:"items"`
}

// SecretNamespacedReference references a secret in a specific namespace.
// The namespace must be specified at the point of use.
type SecretNamespacedReference struct {
	// name is the metadata.name of the referenced secret
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// namespace is the metadata.namespace of the referenced secret
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}

// ConfigMapNameReference references a config map in a specific namespace.
// The namespace must be specified at the point of use.
type ConfigMapNameReference struct {
	// name is the metadata.name of the referenced config map
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// namespace is the metadata.namespace of the referenced config map
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}
