package main

//var (
//    yamlFilename = "../../testdata/examples/deployment.yaml"
//    jsonFilename = "../../testdata/examples/deployment.json"

//    unstructName = "mydep-unstruct"
//    unstructMap  = map[string]interface{}{
//        "apiVersion": "apps/v1",
//        "kind":       "Deployment",
//        "metadata": map[string]interface{}{
//            "name": unstructName,
//            "labels": map[string]interface{}{
//                "app":  unstructName,
//                "type": "deployment",
//            },
//        },
//        "spec": map[string]interface{}{
//            // replicas type is int32, not string.
//            "replicas": 1,
//            "selector": map[string]interface{}{
//                "matchLabels": map[string]interface{}{
//                    "app":  unstructName,
//                    "type": "deployment",
//                },
//            },
//            "template": map[string]interface{}{
//                "metadata": map[string]interface{}{
//                    "labels": map[string]interface{}{
//                        "app":  unstructName,
//                        "type": "deployment",
//                    },
//                },
//                "spec": map[string]interface{}{
//                    "containers": []map[string]interface{}{
//                        {
//                            "name":  "nginx",
//                            "image": "nginx",
//                            "resources": map[string]interface{}{
//                                "limits": map[string]interface{}{
//                                    "cpu": "100m",
//                                },
//                            },
//                        },
//                    },
//                },
//            },
//        },
//    }
//)

var (
	yamlFilename = "../../testdata/examples/namespace.yaml"
	jsonFilename = "../../testdata/examples/namespace.json"

	unstructName = "ns-unstruct"
	unstructMap  = map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]interface{}{
			"name": unstructName,
		},
	}
)
