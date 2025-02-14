// Copyright 2021 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package nodereset

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/chaos-mesh/chaos-mesh/controllers/chaosimpl/gcpchaos/utils"
	"github.com/chaos-mesh/chaos-mesh/controllers/common"
)

type Impl struct {
	client.Client

	Log logr.Logger
}

var _ common.ChaosImpl = (*Impl)(nil)

func (impl *Impl) Apply(ctx context.Context, index int, records []*v1alpha1.Record, chaos v1alpha1.InnerObject) (v1alpha1.Phase, error) {
	gcpchaos, ok := chaos.(*v1alpha1.GCPChaos)
	if !ok {
		err := errors.New("chaos is not gcpchaos")
		impl.Log.Error(err, "chaos is not GCPChaos", "chaos", chaos)
		return v1alpha1.NotInjected, err
	}
	computeService, err := utils.GetComputeService(ctx, impl.Client, gcpchaos)
	if err != nil {
		impl.Log.Error(err, "fail to get the compute service")
		return v1alpha1.NotInjected, err
	}
	var selected v1alpha1.GCPSelector
	json.Unmarshal([]byte(records[index].Id), &selected)
	_, err = computeService.Instances.Reset(selected.Project, selected.Zone, selected.Instance).Do()
	if err != nil {
		impl.Log.Error(err, "fail to reset the instance")
		return v1alpha1.NotInjected, err
	}

	return v1alpha1.Injected, nil
}

func (impl *Impl) Recover(ctx context.Context, index int, records []*v1alpha1.Record, chaos v1alpha1.InnerObject) (v1alpha1.Phase, error) {
	return v1alpha1.NotInjected, nil
}

func NewImpl(c client.Client, log logr.Logger) *Impl {
	return &Impl{
		Client: c,
		Log:    log.WithName("nodereset"),
	}
}
