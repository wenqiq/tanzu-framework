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

	// Specifies nsxKey file path
	// +kubebuilder:validation:Optional
	NSXKey string `json:"nsxKey"`

	// Specifies clusterName
	// +kubebuilder:validation:Optional
	ClusterName string `json:"clusterName"`

	// Specifies NSXIP
	// +kubebuilder:validation:Optional
	NSXIP string `json:"NSXIP"`

	// Specifies VPC path
	// +kubebuilder:validation:Optional
	VPCPath string `json:"VPCPath"`

	MpAdapterConf MpAdapterConf `json:"mp_adapter_conf"`

	CcpAdapterConf CcpAdapterConf `json:"ccp_adapter_conf"`
}

type MpAdapterConf struct {
	// Specifies NSXClientTimeout
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=120
	NSXClientTimeout int `json:"NSXClientTimeout,omitempty"`

	// Specifies InventoryBatchSize
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=50
	InventoryBatchSize int `json:"InventoryBatchSize,omitempty"`

	// Specifies InventoryBatchPeriod
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=5
	InventoryBatchPeriod int `json:"InventoryBatchPeriod,omitempty"`

	// Specifies EnableDebugServer
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	EnableDebugServer int `json:"EnableDebugServer,omitempty"`

	// Specifies APIServerPort
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=16664
	APIServerPort int `json:"APIServerPort,omitempty"`

	// Specifies DebugServerPort
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=16666
	DebugServerPort int `json:"DebugServerPort,omitempty"`

	// Specifies NSXRPCDebug
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	NSXRPCDebug int `json:"NSXRPCDebug,omitempty"`

	// Specifies ConditionTimeout
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=150
	ConditionTimeout int `json:"ConditionTimeout,omitempty"`
}

type CcpAdapterConf struct {
	// Specifies EnableDebugServer
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	EnableDebugServer bool `json:"EnableDebugServer,omitempty"`

	// Specifies APIServerPort
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=16665
	APIServerPort int `json:"APIServerPort,omitempty"`

	// Specifies DebugServerPort
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=16667
	DebugServerPort int `json:"DebugServerPort,omitempty"`

	// Specifies NSXRPCDebug
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	NSXRPCDebug bool `json:"NSXRPCDebug,omitempty"`

	// Specifies RealizeTimeoutSeconds, time to wait for realization
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=60
	RealizeTimeoutSeconds int `json:"RealizeTimeoutSeconds,omitempty"`

	// Specifies RealizeErrorSyncIntervalSeconds, an interval for regularly report latest realization error in background
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=600
	RealizeErrorSyncIntervalSeconds int `json:"RealizeErrorSyncIntervalSeconds,omitempty"`

	// Specifies ReconcilerWorkerCount
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=8
	ReconcilerWorkerCount int `json:"ReconcilerWorkerCount,omitempty"`

	// Specifies ReconcilerQPS, Average QPS = ReconcilerWorkerCount * ReconcilerQPS
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=5.0
	ReconcilerQPS float64 `json:"ReconcilerQPS,omitempty"`

	// Specifies ReconcilerBurst, Peak QPS =  ReconcilerWorkerCount * ReconcilerBurst
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=10
	ReconcilerBurst int `json:"ReconcilerBurst,omitempty"`

	// Specifies ReconcilerResyncSeconds, default 24 Hours
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=86400
	ReconcilerResyncSeconds int `json:"ReconcilerResyncSeconds,omitempty"`
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
