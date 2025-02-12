# Import a Kubernetes Context

After you've deployed Pachyderm with Helm, the Pachyderm context is not created.
Therefore, **you need to manually create a new Pachyderm context with
the embedded current Kubernetes context and activate that context**.

To import a Kubernetes context, complete the following steps:

1. Deploy a Pachyderm cluster using the [Helm installation commands](../helm_install/).

1. Verify that the cluster was successfully deployed:

   ```shell
   kubectl get pods
   ```

   **System Response:**

   ```shell
   NAME                                    READY   STATUS    RESTARTS   AGE
   console-6c989c8d56-ftxk7                1/1     Running   0          3d18h
   etcd-0                                  1/1     Running   0          3d18h
   pachd-f9fd5b6fc-8d774                   1/1     Running   0          3d18h
   pg-bouncer-794d8f68f-sjbbh              1/1     Running   0          3d18h
   ```

   You must see all the `console`, `etcd`, and `pachd` and postgreSQL pods running.

1. Create a new Pachyderm context with the embedded Kubernetes context:

   ```shell
   pachctl config import-kube <new-pachyderm-context> -k `kubectl config current-context`
   ```

1. Verify that the context was successfully created and view the context parameters:

   **Example:**

   ```shell
   pachctl config get context test-context
   ```

   **System Response:**

   ```shell
   {
     "source": "IMPORTED",
     "cluster_name": "minikube",
     "auth_info": "minikube",
     "namespace": "default"
   }
   ```

1. Activate the new Pachyderm context:

   ```shell
   pachctl config set active-context <new-pachyderm-context>
   ```

1. Verify that the new context has been activated:

   ```shell
   pachctl config get active-context
   ```
