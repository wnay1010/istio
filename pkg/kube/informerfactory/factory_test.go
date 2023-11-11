// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package informerfactory

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
	"istio.io/istio/pkg/kube/kubetypes"
)

// Dummy function to simulate creating a SharedIndexInformer
func newDummyInformer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, nil, time.Minute, cache.Indexers{})
}

func TestNewSharedInformerFactory(t *testing.T) {
	factory := NewSharedInformerFactory()
	if factory == nil {
		t.Error("Expected non-nil factory")
	}
}

func TestInformerFor(t *testing.T) {
	factory := NewSharedInformerFactory().(*informerFactory)
	gvr := schema.GroupVersionResource{Group: "testgroup", Version: "v1", Resource: "testresource"}

	informer := factory.InformerFor(gvr, kubetypes.InformerOptions{}, newDummyInformer)
	if informer.Informer == nil {
		t.Error("Expected non-nil informer")
	}
}

func TestStartAndShutdown(t *testing.T) {
	factory := NewSharedInformerFactory().(*informerFactory)
	stopCh := make(chan struct{})

	// Start the factory
	factory.Start(stopCh)

	// Shutdown the factory
	factory.Shutdown()

	if factory.shuttingDown != true {
		t.Error("Expected factory to be in the shutting down state")
	}
}

func TestWaitForCacheSync(t *testing.T) {
	factory := NewSharedInformerFactory().(*informerFactory)
	stopCh := make(chan struct{})
	doneCh := make(chan bool)

	go func() {
		doneCh <- factory.WaitForCacheSync(stopCh)
	}()

	close(stopCh)
	syncStatus := <-doneCh

	if syncStatus != true {
		t.Error("Expected WaitForCacheSync to return true")
	}
}
