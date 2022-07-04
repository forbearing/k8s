package job

//rawData := map[string]interface{}{
//    "apiVersion": "patch/v1",
//    "kind":       "CronJob",
//    "metadata": map[string]interface{}{
//        "name": "mycj",
//        "labels": map[string]interface{}{
//            "name": "mycj",
//        },
//    },
//    "spec": map[string]interface{}{
//        "template": map[string]interface{}{
//            "spec": map[string]interface{}{
//                "restartPolicy": "Never",
//                "containers": []map[string]interface{}{
//                    {
//                        "name":            "echo",
//                        "image":           "busybox",
//                        "imagePullPolicy": "IfNotPresent",
//                        "command":         []string{"sh", "-c", "echo hello job"},
//                    },
//                },
//            },
//        },
//    },
//}
