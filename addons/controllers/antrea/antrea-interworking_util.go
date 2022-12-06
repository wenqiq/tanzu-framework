// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterapiutil "sigs.k8s.io/cluster-api/util"
	clusterapipatchutil "sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/vmware-tanzu/tanzu-framework/addons/pkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/addons/pkg/util"

	cniv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/addonconfigs/cni/v1alpha1"
)

// antreaInterworkingConfigSpec defines the desired state of AntreaInterworkingConfig
type antreaInterworkingConfigSpec struct {
	NSXCert        string         `yaml:"nsx_cert"`
	NSXKey         string         `json:"nsxKey"`
	ClusterName    string         `json:"clusterName"`
	NSXIP          string         `json:"NSXIP"`
	VPCPath        string         `json:"VPCPath"`
	MpAdapterConf  MpAdapterConf  `json:"mp_adapter_conf"`
	CcpAdapterConf CcpAdapterConf `json:"ccp_adapter_conf"`
}

type MpAdapterConf struct {
	NSXClientTimeout     int `json:"NSXClientTimeout,omitempty"`
	InventoryBatchSize   int `json:"InventoryBatchSize,omitempty"`
	InventoryBatchPeriod int `json:"InventoryBatchPeriod,omitempty"`
	EnableDebugServer    int `json:"EnableDebugServer,omitempty"`
	APIServerPort        int `json:"APIServerPort,omitempty"`
	DebugServerPort      int `json:"DebugServerPort,omitempty"`
	NSXRPCDebug          int `json:"NSXRPCDebug,omitempty"`
	ConditionTimeout     int `json:"ConditionTimeout,omitempty"`
}

type CcpAdapterConf struct {
	EnableDebugServer               bool    `json:"EnableDebugServer,omitempty"`
	APIServerPort                   int     `json:"APIServerPort,omitempty"`
	DebugServerPort                 int     `json:"DebugServerPort,omitempty"`
	NSXRPCDebug                     bool    `json:"NSXRPCDebug,omitempty"`
	RealizeTimeoutSeconds           int     `json:"RealizeTimeoutSeconds,omitempty"`
	RealizeErrorSyncIntervalSeconds int     `json:"RealizeErrorSyncIntervalSeconds,omitempty"`
	ReconcilerWorkerCount           int     `json:"ReconcilerWorkerCount,omitempty"`
	ReconcilerQPS                   float64 `json:"ReconcilerQPS,omitempty"`
	ReconcilerBurst                 int     `json:"ReconcilerBurst,omitempty"`
	ReconcilerResyncSeconds         int     `json:"ReconcilerResyncSeconds,omitempty"`
}

func (r *AntreaConfigReconciler) ReconcileAntreaInterworkingConfig(ctx context.Context, req ctrl.Request, cluster *clusterapiv1beta1.Cluster) (_ ctrl.Result, retErr error) {
	var log logr.Logger

	antreaInterworkingConfig := &cniv1alpha1.AntreaInterworkingConfig{}
	if err := r.Client.Get(ctx, req.NamespacedName, antreaInterworkingConfig); err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info(fmt.Sprintf("AntreaInterworkingConfig resource '%v' not found", req.NamespacedName))
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	patchHelper, err := clusterapipatchutil.NewHelper(antreaInterworkingConfig, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Patch AntreaConfig before returning the function
	defer func() {
		log.Info("Patching AntreaConfig")
		if err := patchHelper.Patch(ctx, antreaInterworkingConfig); err != nil {
			log.Error(err, "Error patching antreaInterworkingConfig")
			retErr = err
		}
		log.Info("Successfully patched antreaInterworkingConfig")
	}()

	// If AntreaConfig is marked for deletion, then no reconciliation is needed
	if !antreaInterworkingConfig.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}

	ownerReference := metav1.OwnerReference{
		APIVersion: clusterapiv1beta1.GroupVersion.String(),
		Kind:       cluster.Kind,
		Name:       cluster.Name,
		UID:        cluster.UID,
	}

	if !clusterapiutil.HasOwnerRef(antreaInterworkingConfig.OwnerReferences, ownerReference) {
		log.Info("Adding owner reference to AntreaConfig")
		antreaInterworkingConfig.OwnerReferences = clusterapiutil.EnsureOwnerRef(antreaInterworkingConfig.OwnerReferences, ownerReference)
	}

	if err := r.ReconcileAntreaInterworkingConfigDataValue(ctx, antreaInterworkingConfig, cluster, log); err != nil {
		log.Error(err, "Error creating antreaConfig data value secret")
		return ctrl.Result{}, err
	}

	// update status.secretRef
	dataValueSecretName := util.GenerateDataValueSecretName(cluster.Name, constants.AntreaAddonName)
	antreaInterworkingConfig.Status.SecretRef = dataValueSecretName

	// if err := r.ReconcileAntreaConfigNormal(ctx, antreaConfig, cluster, log); err != nil {
	// 	log.Error(err, "Error reconciling AntreaConfig to create data value secret")
	// 	return ctrl.Result{}, err
	// }

	log.Info("Successfully reconciled antreaInterworkingConfig")
	return ctrl.Result{}, nil
}

func (r *AntreaConfigReconciler) ReconcileAntreaInterworkingConfigDataValue(
	ctx context.Context,
	antreaInterworkingConfig *cniv1alpha1.AntreaInterworkingConfig,
	cluster *clusterapiv1beta1.Cluster,
	log logr.Logger) (retErr error) {

	// prepare data values secret for AntreaConfig
	antreaDataValuesSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GenerateDataValueSecretName(cluster.Name, constants.AntreaAddonName),
			Namespace: antreaInterworkingConfig.Namespace,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: clusterapiv1beta1.GroupVersion.String(),
				Kind:       cluster.Kind,
				Name:       cluster.Name,
				UID:        cluster.UID,
			}},
		},
	}

	antreaDataValuesSecretMutateFn := func() error {
		antreaDataValuesSecret.Type = corev1.SecretTypeOpaque
		antreaDataValuesSecret.StringData = make(map[string]string)

		// marshall the yaml contents
		antreaConfigYaml, err := mapAntreaInterworkingConfigSpec(antreaInterworkingConfig)
		if err != nil {
			return err
		}

		dataValueYamlBytes, err := yaml.Marshal(antreaConfigYaml)
		if err != nil {
			log.Error(err, "Error marshaling AntreaConfig to Yaml")
			return err
		}

		antreaDataValuesSecret.StringData[constants.TKGDataValueFileName] = string(dataValueYamlBytes)

		return nil
	}

	result, err := controllerutil.CreateOrPatch(ctx, r.Client, antreaDataValuesSecret, antreaDataValuesSecretMutateFn)
	if err != nil {
		log.Error(err, "Error creating or patching antrea data values secret")
		return err
	}

	log.Info(fmt.Sprintf("Resource %s data values secret %s", constants.AntreaAddonName, result))

	return nil
}

func mapAntreaInterworkingConfigSpec(config *cniv1alpha1.AntreaInterworkingConfig) (*antreaInterworkingConfigSpec, error) {
	configSpec := &antreaInterworkingConfigSpec{}

	configSpec.NSXCert = config.Spec.AntreaInterworking.AntreaConfigDataValue.NSXCert
	configSpec.NSXKey = config.Spec.AntreaInterworking.AntreaConfigDataValue.NSXKey
	configSpec.NSXIP = config.Spec.AntreaInterworking.AntreaConfigDataValue.NSXIP
	configSpec.ClusterName = config.Spec.AntreaInterworking.AntreaConfigDataValue.ClusterName
	configSpec.VPCPath = config.Spec.AntreaInterworking.AntreaConfigDataValue.VPCPath
	configSpec.CcpAdapterConf = CcpAdapterConf{
		EnableDebugServer:               config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.EnableDebugServer,
		APIServerPort:                   config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.APIServerPort,
		DebugServerPort:                 config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.DebugServerPort,
		NSXRPCDebug:                     config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.NSXRPCDebug,
		RealizeTimeoutSeconds:           config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.RealizeTimeoutSeconds,
		RealizeErrorSyncIntervalSeconds: config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.RealizeErrorSyncIntervalSeconds,
		ReconcilerWorkerCount:           config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.ReconcilerWorkerCount,
		ReconcilerQPS:                   config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.ReconcilerQPS,
		ReconcilerBurst:                 config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.ReconcilerBurst,
		ReconcilerResyncSeconds:         config.Spec.AntreaInterworking.AntreaConfigDataValue.CcpAdapterConf.ReconcilerResyncSeconds,
	}
	configSpec.MpAdapterConf = MpAdapterConf{
		NSXClientTimeout:     config.Spec.AntreaInterworking.AntreaConfigDataValue.MpAdapterConf.NSXClientTimeout,
		InventoryBatchSize:   config.Spec.AntreaInterworking.AntreaConfigDataValue.MpAdapterConf.InventoryBatchSize,
		InventoryBatchPeriod: config.Spec.AntreaInterworking.AntreaConfigDataValue.MpAdapterConf.InventoryBatchPeriod,
		EnableDebugServer:    config.Spec.AntreaInterworking.AntreaConfigDataValue.MpAdapterConf.EnableDebugServer,
		APIServerPort:        config.Spec.AntreaInterworking.AntreaConfigDataValue.MpAdapterConf.APIServerPort,
		DebugServerPort:      config.Spec.AntreaInterworking.AntreaConfigDataValue.MpAdapterConf.DebugServerPort,
		NSXRPCDebug:          config.Spec.AntreaInterworking.AntreaConfigDataValue.MpAdapterConf.NSXRPCDebug,
		ConditionTimeout:     config.Spec.AntreaInterworking.AntreaConfigDataValue.MpAdapterConf.ConditionTimeout,
	}
	return configSpec, nil
}
