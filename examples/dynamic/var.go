package main

var (
	deployUnstructName = "mydep-unstruct"
	deployUnstructData = map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name": deployUnstructName,
			"labels": map[string]interface{}{
				"app":  deployUnstructName,
				"type": "deployment",
			},
		},
		"spec": map[string]interface{}{
			// replicas type is int32, not string.
			"replicas": 1,
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"app":  deployUnstructName,
					"type": "deployment",
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app":  deployUnstructName,
						"type": "deployment",
					},
				},
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":  "nginx",
							"image": "nginx",
							"resources": map[string]interface{}{
								"limits": map[string]interface{}{
									"cpu": "100m",
								},
							},
						},
					},
				},
			},
		},
	}

	podUnstructName = "pod-unstruct"
	podUnstructData = map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]interface{}{
			"name": podUnstructName,
			"labels": map[string]interface{}{
				"app":  "podUnstructName",
				"type": "pod",
			},
		},
		"spec": map[string]interface{}{
			"containers": []map[string]interface{}{
				{
					"name":  "nginx",
					"image": "nginx",
					"ports": []map[string]interface{}{
						{
							"name":          "http",
							"protocol":      "TCP",
							"containerPort": 80,
						},
						{
							"name":          "https",
							"protocol":      "TCP",
							"containerPort": 443,
						},
					},
				},
			},
		},
	}

	nsUnstructName = "ns-unstruct"
	nsUnstructData = map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]interface{}{
			"name": nsUnstructName,
		},
	}

	pvUnstructName = "pv-unstruct"
	pvUnstructData = map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "PersistentVolume",
		"metadata": map[string]interface{}{
			"name": pvUnstructName,
		},
		"spec": map[string]interface{}{
			"accessModes": []interface{}{"ReadWriteOnce"},
			"capacity": map[string]interface{}{
				"storage": "1Mi",
			},
			"hostPath": map[string]interface{}{
				"path": "/tmp/hahahahah",
			},
		},
	}

	crUnstructName = "cr-unstruct"
	crUnstructData = map[string]interface{}{
		"apiVersion": "rbac.authorization.k8s.io/v1",
		"kind":       "ClusterRole",
		"metadata": map[string]interface{}{
			"name": crUnstructName,
		},
	}
)
