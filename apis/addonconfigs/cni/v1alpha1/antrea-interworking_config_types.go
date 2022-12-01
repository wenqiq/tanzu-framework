package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AntreaInterworkingON struct {
	AntreaNsx AntreaNsx `json:"antrea_nsx"`
}

type AntreaNsx struct {
	Enable bool `json:"enable"`
}

// AntreaInterworkingConfigSpec defines the desired state of AntreaInterworkingConfig
type AntreaInterworkingConfigSpec struct {
	AntreaInterworking AntreaInterworking `json:"antrea_interworking,omitempty"`
}

type AntreaInterworking struct {
	AntreaConfigDataValue AntreaInterworkingConfigDataValue `json:"config,omitempty"`
}

type AntreaInterworkingConfigDataValue struct {
	// Specifies nsxCert file path
	// +kubebuilder:validation:Optional
	NSXCert string `json:"nsxCert"`
}

// AntreaInterworkingConfigStatus defines the observed state of AntreaInterworkingConfig
type AntreaInterworkingConfigStatus struct {
	// Reference to the data value secret created by controller
	// +kubebuilder:validation:Optional
	SecretRef string `json:"secretRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=antrea-interworkingconfigs,shortName=antrea-interworkingconf,scope=Namespaced
// +kubebuilder:printcolumn:name="NSXCert",type="string",JSONPath=".spec.antrea_interworking.config.nsxCert",description="The NSXCert file path"

// AntreaInterworkingConfig is the Schema for the antrea-interworking configs API
type AntreaInterworkingConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AntreaInterworkingConfigSpec   `json:"spec"`
	Status AntreaInterworkingConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AntreaInterworkingConfigList contains a list of AntreaInterworkingConfig
type AntreaInterworkingConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AntreaInterworkingConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AntreaInterworkingConfig{}, &AntreaInterworkingConfigList{})
}
