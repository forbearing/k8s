/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package endpoint

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	discovery "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	utilpointer "k8s.io/utils/pointer"
)

func TestDetermineNeededServiceUpdates(t *testing.T) {
	testCases := []struct {
		name  string
		a     sets.String
		b     sets.String
		union sets.String
		xor   sets.String
	}{
		{
			name:  "no services changed",
			a:     sets.NewString("a", "b", "c"),
			b:     sets.NewString("a", "b", "c"),
			xor:   sets.NewString(),
			union: sets.NewString("a", "b", "c"),
		},
		{
			name:  "all old services removed, new services added",
			a:     sets.NewString("a", "b", "c"),
			b:     sets.NewString("d", "e", "f"),
			xor:   sets.NewString("a", "b", "c", "d", "e", "f"),
			union: sets.NewString("a", "b", "c", "d", "e", "f"),
		},
		{
			name:  "all old services removed, no new services added",
			a:     sets.NewString("a", "b", "c"),
			b:     sets.NewString(),
			xor:   sets.NewString("a", "b", "c"),
			union: sets.NewString("a", "b", "c"),
		},
		{
			name:  "no old services, but new services added",
			a:     sets.NewString(),
			b:     sets.NewString("a", "b", "c"),
			xor:   sets.NewString("a", "b", "c"),
			union: sets.NewString("a", "b", "c"),
		},
		{
			name:  "one service removed, one service added, two unchanged",
			a:     sets.NewString("a", "b", "c"),
			b:     sets.NewString("b", "c", "d"),
			xor:   sets.NewString("a", "d"),
			union: sets.NewString("a", "b", "c", "d"),
		},
		{
			name:  "no services",
			a:     sets.NewString(),
			b:     sets.NewString(),
			xor:   sets.NewString(),
			union: sets.NewString(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			retval := determineNeededServiceUpdates(testCase.a, testCase.b, false)
			if !retval.Equal(testCase.xor) {
				t.Errorf("%s (with podChanged=false): expected: %v  got: %v", testCase.name, testCase.xor.List(), retval.List())
			}

			retval = determineNeededServiceUpdates(testCase.a, testCase.b, true)
			if !retval.Equal(testCase.union) {
				t.Errorf("%s (with podChanged=true): expected: %v  got: %v", testCase.name, testCase.union.List(), retval.List())
			}
		})
	}
}

func TestShouldPodBeInEndpoints(t *testing.T) {
	testCases := []struct {
		name               string
		pod                *v1.Pod
		expected           bool
		includeTerminating bool
	}{
		// Pod should not be in endpoints:
		{
			name: "Failed pod with Never RestartPolicy",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
				},
				Status: v1.PodStatus{
					Phase: v1.PodFailed,
					PodIP: "1.2.3.4",
				},
			},
			expected: false,
		},
		{
			name: "Succeeded pod with Never RestartPolicy",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
				},
				Status: v1.PodStatus{
					Phase: v1.PodSucceeded,
					PodIP: "1.2.3.4",
				},
			},
			expected: false,
		},
		{
			name: "Succeeded pod with OnFailure RestartPolicy",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyOnFailure,
				},
				Status: v1.PodStatus{
					Phase: v1.PodSucceeded,
					PodIP: "1.2.3.4",
				},
			},
			expected: false,
		},
		{
			name: "Empty Pod IPs, Running pod with OnFailure RestartPolicy",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
				},
				Status: v1.PodStatus{
					Phase:  v1.PodRunning,
					PodIP:  "",
					PodIPs: []v1.PodIP{},
				},
			},
			expected: false,
		},
		{
			name: "Terminating Pod with includeTerminating=false",
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: &metav1.Time{
						Time: time.Now(),
					},
				},
				Spec: v1.PodSpec{},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					PodIP: "1.2.3.4",
				},
			},
			expected: false,
		},
		{
			name: "Failed pod with Always RestartPolicy",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyAlways,
				},
				Status: v1.PodStatus{
					Phase: v1.PodFailed,
					PodIP: "1.2.3.4",
				},
			},
			expected: false,
		},
		// Pod should be in endpoints:
		{
			name: "Pending pod with Never RestartPolicy",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
				},
				Status: v1.PodStatus{
					Phase: v1.PodPending,
					PodIP: "1.2.3.4",
				},
			},
			expected: true,
		},
		{
			name: "Unknown pod with OnFailure RestartPolicy",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyOnFailure,
				},
				Status: v1.PodStatus{
					Phase: v1.PodUnknown,
					PodIP: "1.2.3.4",
				},
			},
			expected: true,
		},
		{
			name: "Running pod with Never RestartPolicy",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					PodIP: "1.2.3.4",
				},
			},
			expected: true,
		},
		{
			name: "Multiple Pod IPs, Running pod with OnFailure RestartPolicy",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
				},
				Status: v1.PodStatus{
					Phase:  v1.PodRunning,
					PodIPs: []v1.PodIP{{IP: "1.2.3.4"}, {IP: "1234::5678:0000:0000:9abc:def0"}},
				},
			},
			expected: true,
		},
		{
			name: "Terminating Pod with includeTerminating=true",
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: &metav1.Time{
						Time: time.Now(),
					},
				},
				Spec: v1.PodSpec{},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					PodIP: "1.2.3.4",
				},
			},
			expected:           true,
			includeTerminating: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := ShouldPodBeInEndpoints(test.pod, test.includeTerminating)
			if result != test.expected {
				t.Errorf("expected: %t, got: %t", test.expected, result)
			}
		})
	}
}

func TestShouldSetHostname(t *testing.T) {
	testCases := map[string]struct {
		pod      *v1.Pod
		service  *v1.Service
		expected bool
	}{
		"all matching": {
			pod:      genSimplePod("ns", "foo", "svc-name"),
			service:  genSimpleSvc("ns", "svc-name"),
			expected: true,
		},
		"all matching, hostname not set": {
			pod:      genSimplePod("ns", "", "svc-name"),
			service:  genSimpleSvc("ns", "svc-name"),
			expected: false,
		},
		"all set, different name/subdomain": {
			pod:      genSimplePod("ns", "hostname", "subdomain"),
			service:  genSimpleSvc("ns", "name"),
			expected: false,
		},
		"all set, different namespace": {
			pod:      genSimplePod("ns1", "hostname", "svc-name"),
			service:  genSimpleSvc("ns2", "svc-name"),
			expected: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			result := ShouldSetHostname(testCase.pod, testCase.service)
			if result != testCase.expected {
				t.Errorf("expected: %t, got: %t", testCase.expected, result)
			}
		})
	}
}

func genSimplePod(namespace, hostname, subdomain string) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Hostname:  hostname,
			Subdomain: subdomain,
		},
	}
}

func genSimpleSvc(namespace, name string) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func TestServiceSelectorCache_GetPodServiceMemberships(t *testing.T) {
	fakeInformerFactory := informers.NewSharedInformerFactory(&fake.Clientset{}, 0*time.Second)
	for i := 0; i < 3; i++ {
		service := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("service-%d", i),
				Namespace: "test",
			},
			Spec: v1.ServiceSpec{
				Selector: map[string]string{
					"app": fmt.Sprintf("test-%d", i),
				},
			},
		}
		fakeInformerFactory.Core().V1().Services().Informer().GetStore().Add(service)
	}
	var pods []*v1.Pod
	for i := 0; i < 5; i++ {
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Name:      fmt.Sprintf("test-pod-%d", i),
				Labels: map[string]string{
					"app":   fmt.Sprintf("test-%d", i),
					"label": fmt.Sprintf("label-%d", i),
				},
			},
		}
		pods = append(pods, pod)
	}

	cache := NewServiceSelectorCache()
	tests := []struct {
		name   string
		pod    *v1.Pod
		expect sets.String
	}{
		{
			name:   "get servicesMemberships for pod-0",
			pod:    pods[0],
			expect: sets.NewString("test/service-0"),
		},
		{
			name:   "get servicesMemberships for pod-1",
			pod:    pods[1],
			expect: sets.NewString("test/service-1"),
		},
		{
			name:   "get servicesMemberships for pod-2",
			pod:    pods[2],
			expect: sets.NewString("test/service-2"),
		},
		{
			name:   "get servicesMemberships for pod-3",
			pod:    pods[3],
			expect: sets.NewString(),
		},
		{
			name:   "get servicesMemberships for pod-4",
			pod:    pods[4],
			expect: sets.NewString(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services, err := cache.GetPodServiceMemberships(fakeInformerFactory.Core().V1().Services().Lister(), test.pod)
			if err != nil {
				t.Errorf("Error from cache.GetPodServiceMemberships: %v", err)
			} else if !services.Equal(test.expect) {
				t.Errorf("Expect service %v, but got %v", test.expect, services)
			}
		})
	}
}

func TestServiceSelectorCache_Update(t *testing.T) {
	var selectors []labels.Selector
	for i := 0; i < 5; i++ {
		selector := labels.Set(map[string]string{"app": fmt.Sprintf("test-%d", i)}).AsSelectorPreValidated()
		selectors = append(selectors, selector)
	}
	tests := []struct {
		name   string
		key    string
		cache  *ServiceSelectorCache
		update map[string]string
		expect labels.Selector
	}{
		{
			name:   "add test/service-0",
			key:    "test/service-0",
			cache:  generateServiceSelectorCache(map[string]labels.Selector{}),
			update: map[string]string{"app": "test-0"},
			expect: selectors[0],
		},
		{
			name:   "add test/service-1",
			key:    "test/service-1",
			cache:  generateServiceSelectorCache(map[string]labels.Selector{"test/service-0": selectors[0]}),
			update: map[string]string{"app": "test-1"},
			expect: selectors[1],
		},
		{
			name:   "update test/service-2",
			key:    "test/service-2",
			cache:  generateServiceSelectorCache(map[string]labels.Selector{"test/service-2": selectors[2]}),
			update: map[string]string{"app": "test-0"},
			expect: selectors[0],
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := test.cache.Update(test.key, test.update)
			if !reflect.DeepEqual(selector, test.expect) {
				t.Errorf("Expect selector %v , but got %v", test.expect, selector)
			}
		})
	}
}

func generateServiceSelectorCache(cache map[string]labels.Selector) *ServiceSelectorCache {
	return &ServiceSelectorCache{
		cache: cache,
	}
}

func BenchmarkGetPodServiceMemberships(b *testing.B) {
	// init fake service informer.
	fakeInformerFactory := informers.NewSharedInformerFactory(&fake.Clientset{}, 0*time.Second)
	for i := 0; i < 1000; i++ {
		service := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("service-%d", i),
				Namespace: "test",
			},
			Spec: v1.ServiceSpec{
				Selector: map[string]string{
					"app": fmt.Sprintf("test-%d", i),
				},
			},
		}
		fakeInformerFactory.Core().V1().Services().Informer().GetStore().Add(service)
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "test-pod-0",
			Labels: map[string]string{
				"app": "test-0",
			},
		},
	}

	cache := NewServiceSelectorCache()
	expect := sets.NewString("test/service-0")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		services, err := cache.GetPodServiceMemberships(fakeInformerFactory.Core().V1().Services().Lister(), pod)
		if err != nil {
			b.Fatalf("Error from GetPodServiceMemberships(): %v", err)
		}
		if len(services) != len(expect) {
			b.Errorf("Expect services size %d, but got: %v", len(expect), len(services))
		}
	}
}

func Test_podChanged(t *testing.T) {
	testCases := []struct {
		testName      string
		modifier      func(*v1.Pod, *v1.Pod)
		podChanged    bool
		labelsChanged bool
	}{
		{
			testName:      "no changes",
			modifier:      func(old, new *v1.Pod) {},
			podChanged:    false,
			labelsChanged: false,
		}, {
			testName: "change NodeName",
			modifier: func(old, new *v1.Pod) {
				new.Spec.NodeName = "changed"
			},
			// NodeName can only change before the pod has an IP, and we don't care about the
			// pod yet at that point so we ignore this change
			podChanged:    false,
			labelsChanged: false,
		}, {
			testName: "change ResourceVersion",
			modifier: func(old, new *v1.Pod) {
				new.ObjectMeta.ResourceVersion = "changed"
			},
			// ResourceVersion is intentionally ignored if nothing else changed
			podChanged:    false,
			labelsChanged: false,
		}, {
			testName: "add primary IPv4",
			modifier: func(old, new *v1.Pod) {
				new.Status.PodIP = "1.2.3.4"
				new.Status.PodIPs = []v1.PodIP{{IP: "1.2.3.4"}}
			},
			podChanged:    true,
			labelsChanged: false,
		}, {
			testName: "modify primary IPv4",
			modifier: func(old, new *v1.Pod) {
				old.Status.PodIP = "1.2.3.4"
				old.Status.PodIPs = []v1.PodIP{{IP: "1.2.3.4"}}
				new.Status.PodIP = "2.3.4.5"
				new.Status.PodIPs = []v1.PodIP{{IP: "2.3.4.5"}}
			},
			podChanged:    true,
			labelsChanged: false,
		}, {
			testName: "add primary IPv6",
			modifier: func(old, new *v1.Pod) {
				new.Status.PodIP = "fd00:10:96::1"
				new.Status.PodIPs = []v1.PodIP{{IP: "fd00:10:96::1"}}
			},
			podChanged:    true,
			labelsChanged: false,
		}, {
			testName: "modify primary IPv6",
			modifier: func(old, new *v1.Pod) {
				old.Status.PodIP = "fd00:10:96::1"
				old.Status.PodIPs = []v1.PodIP{{IP: "fd00:10:96::1"}}
				new.Status.PodIP = "fd00:10:96::2"
				new.Status.PodIPs = []v1.PodIP{{IP: "fd00:10:96::2"}}
			},
			podChanged:    true,
			labelsChanged: false,
		}, {
			testName: "add secondary IP",
			modifier: func(old, new *v1.Pod) {
				old.Status.PodIP = "1.2.3.4"
				old.Status.PodIPs = []v1.PodIP{{IP: "1.2.3.4"}}
				new.Status.PodIP = "1.2.3.4"
				new.Status.PodIPs = []v1.PodIP{{IP: "1.2.3.4"}, {IP: "fd00:10:96::1"}}
			},
			podChanged:    true,
			labelsChanged: false,
		}, {
			testName: "modify secondary IP",
			modifier: func(old, new *v1.Pod) {
				old.Status.PodIP = "1.2.3.4"
				old.Status.PodIPs = []v1.PodIP{{IP: "1.2.3.4"}, {IP: "fd00:10:96::1"}}
				new.Status.PodIP = "1.2.3.4"
				new.Status.PodIPs = []v1.PodIP{{IP: "1.2.3.4"}, {IP: "fd00:10:96::2"}}
			},
			podChanged:    true,
			labelsChanged: false,
		}, {
			testName: "remove secondary IP",
			modifier: func(old, new *v1.Pod) {
				old.Status.PodIP = "1.2.3.4"
				old.Status.PodIPs = []v1.PodIP{{IP: "1.2.3.4"}, {IP: "fd00:10:96::1"}}
				new.Status.PodIP = "1.2.3.4"
				new.Status.PodIPs = []v1.PodIP{{IP: "1.2.3.4"}}
			},
			podChanged:    true,
			labelsChanged: false,
		}, {
			testName: "change readiness",
			modifier: func(old, new *v1.Pod) {
				new.Status.Conditions[0].Status = v1.ConditionTrue
			},
			podChanged:    true,
			labelsChanged: false,
		}, {
			testName: "mark for deletion",
			modifier: func(old, new *v1.Pod) {
				now := metav1.NewTime(time.Now().UTC())
				new.ObjectMeta.DeletionTimestamp = &now
			},
			podChanged:    true,
			labelsChanged: false,
		}, {
			testName: "add label",
			modifier: func(old, new *v1.Pod) {
				new.Labels["label"] = "new"
			},
			podChanged:    false,
			labelsChanged: true,
		}, {
			testName: "modify label",
			modifier: func(old, new *v1.Pod) {
				old.Labels["label"] = "old"
				new.Labels["label"] = "new"
			},
			podChanged:    false,
			labelsChanged: true,
		}, {
			testName: "remove label",
			modifier: func(old, new *v1.Pod) {
				old.Labels["label"] = "old"
			},
			podChanged:    false,
			labelsChanged: true,
		},
	}

	orig := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "pod",
			Labels:    map[string]string{"foo": "bar"},
		},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{
				{Type: v1.PodReady, Status: v1.ConditionFalse},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			old := orig.DeepCopy()
			new := old.DeepCopy()
			tc.modifier(old, new)

			podChanged, labelsChanged := podEndpointsChanged(old, new)
			if podChanged != tc.podChanged {
				t.Errorf("Expected podChanged to be %t, got %t", tc.podChanged, podChanged)
			}
			if labelsChanged != tc.labelsChanged {
				t.Errorf("Expected labelsChanged to be %t, got %t", tc.labelsChanged, labelsChanged)
			}
		})
	}
}

func TestEndpointsEqualBeyondHash(t *testing.T) {
	tests := []struct {
		name     string
		ep1      *discovery.Endpoint
		ep2      *discovery.Endpoint
		expected bool
	}{
		{
			name: "No change",
			ep1: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			ep2: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			expected: true,
		},
		{
			name: "NodeName changed",
			ep1: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			ep2: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				NodeName:  utilpointer.StringPtr("node-2"),
			},
			expected: false,
		},
		{
			name: "Zone changed",
			ep1: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				Zone:      utilpointer.StringPtr("zone-1"),
			},
			ep2: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				Zone:      utilpointer.StringPtr("zone-2"),
			},
			expected: false,
		},
		{
			name: "Ready condition changed",
			ep1: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			ep2: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(false),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			expected: false,
		},
		{
			name: "Serving condition changed from nil to true",
			ep1: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready:       utilpointer.BoolPtr(true),
					Serving:     nil,
					Terminating: nil,
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			ep2: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready:       utilpointer.BoolPtr(true),
					Serving:     utilpointer.BoolPtr(true),
					Terminating: utilpointer.BoolPtr(false),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			expected: true,
		},
		{
			name: "Serving condition changed from false to true",
			ep1: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready:       utilpointer.BoolPtr(true),
					Serving:     utilpointer.BoolPtr(false),
					Terminating: utilpointer.BoolPtr(false),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			ep2: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready:       utilpointer.BoolPtr(true),
					Serving:     utilpointer.BoolPtr(true),
					Terminating: utilpointer.BoolPtr(false),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			expected: false,
		},
		{
			name: "Pod name changed",
			ep1: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			ep2: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod1"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			expected: false,
		},
		{
			name: "Pod resourceVersion changed",
			ep1: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0", ResourceVersion: "1"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			ep2: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0", ResourceVersion: "2"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			expected: true,
		},
		{
			name: "Pod resourceVersion removed",
			ep1: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0", ResourceVersion: "1"},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			ep2: &discovery.Endpoint{
				Conditions: discovery.EndpointConditions{
					Ready: utilpointer.BoolPtr(true),
				},
				Addresses: []string{"10.0.0.1"},
				TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "default", Name: "pod0", ResourceVersion: ""},
				Zone:      utilpointer.StringPtr("zone-1"),
				NodeName:  utilpointer.StringPtr("node-1"),
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndpointsEqualBeyondHash(tt.ep1, tt.ep2); got != tt.expected {
				t.Errorf("EndpointsEqualBeyondHash() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEndpointSubsetsEqualIgnoreResourceVersion(t *testing.T) {
	copyAndMutateEndpointSubset := func(orig *v1.EndpointSubset, mutator func(*v1.EndpointSubset)) *v1.EndpointSubset {
		newSubSet := orig.DeepCopy()
		mutator(newSubSet)
		return newSubSet
	}
	es1 := &v1.EndpointSubset{
		Addresses: []v1.EndpointAddress{
			{
				IP:        "1.1.1.1",
				TargetRef: &v1.ObjectReference{Kind: "Pod", Name: "pod1-1", Namespace: "ns", ResourceVersion: "1"},
			},
		},
		NotReadyAddresses: []v1.EndpointAddress{
			{
				IP:        "1.1.1.2",
				TargetRef: &v1.ObjectReference{Kind: "Pod", Name: "pod1-2", Namespace: "ns2", ResourceVersion: "2"},
			},
		},
		Ports: []v1.EndpointPort{{Port: 8081, Protocol: "TCP"}},
	}
	es2 := &v1.EndpointSubset{
		Addresses: []v1.EndpointAddress{
			{
				IP:        "2.2.2.1",
				TargetRef: &v1.ObjectReference{Kind: "Pod", Name: "pod2-1", Namespace: "ns", ResourceVersion: "3"},
			},
		},
		NotReadyAddresses: []v1.EndpointAddress{
			{
				IP:        "2.2.2.2",
				TargetRef: &v1.ObjectReference{Kind: "Pod", Name: "pod2-2", Namespace: "ns2", ResourceVersion: "4"},
			},
		},
		Ports: []v1.EndpointPort{{Port: 8082, Protocol: "TCP"}},
	}
	tests := []struct {
		name     string
		subsets1 []v1.EndpointSubset
		subsets2 []v1.EndpointSubset
		expected bool
	}{
		{
			name:     "Subsets removed",
			subsets1: []v1.EndpointSubset{*es1, *es2},
			subsets2: []v1.EndpointSubset{*es1},
			expected: false,
		},
		{
			name:     "Ready Pod IP changed",
			subsets1: []v1.EndpointSubset{*es1, *es2},
			subsets2: []v1.EndpointSubset{*copyAndMutateEndpointSubset(es1, func(es *v1.EndpointSubset) {
				es.Addresses[0].IP = "1.1.1.10"
			}), *es2},
			expected: false,
		},
		{
			name:     "NotReady Pod IP changed",
			subsets1: []v1.EndpointSubset{*es1, *es2},
			subsets2: []v1.EndpointSubset{*es1, *copyAndMutateEndpointSubset(es2, func(es *v1.EndpointSubset) {
				es.NotReadyAddresses[0].IP = "2.2.2.10"
			})},
			expected: false,
		},
		{
			name:     "Pod ResourceVersion changed",
			subsets1: []v1.EndpointSubset{*es1, *es2},
			subsets2: []v1.EndpointSubset{*es1, *copyAndMutateEndpointSubset(es2, func(es *v1.EndpointSubset) {
				es.Addresses[0].TargetRef.ResourceVersion = "100"
			})},
			expected: true,
		},
		{
			name:     "Pod ResourceVersion removed",
			subsets1: []v1.EndpointSubset{*es1, *es2},
			subsets2: []v1.EndpointSubset{*es1, *copyAndMutateEndpointSubset(es2, func(es *v1.EndpointSubset) {
				es.Addresses[0].TargetRef.ResourceVersion = ""
			})},
			expected: true,
		},
		{
			name:     "Ports changed",
			subsets1: []v1.EndpointSubset{*es1, *es2},
			subsets2: []v1.EndpointSubset{*es1, *copyAndMutateEndpointSubset(es1, func(es *v1.EndpointSubset) {
				es.Ports[0].Port = 8082
			})},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndpointSubsetsEqualIgnoreResourceVersion(tt.subsets1, tt.subsets2); got != tt.expected {
				t.Errorf("semanticIgnoreResourceVersion.DeepEqual() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
