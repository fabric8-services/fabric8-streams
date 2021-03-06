# Deployment guide

This guide shows how to deploy RH AMQ Streams based upon the [official
documentation](https://access.redhat.com/documentation/en-us/red_hat_amq/7.2/html/using_amq_streams_on_openshift_container_platform/getting-started-str#downloads-str) with customizations suited for needs.

## Optionally: Start a local OpenShift cluster using Minishift

For a local deployment 

WAIT for 

If you don't have Minishift, please consult [its documentation](https://docs.okd.io/latest/minishift/getting-started/preparing-to-install.html) on how to get and install it.

**NOTE:**
This guide is tested using [minishift v1.30.0](https://github.com/minishift/minishift/releases/tag/v1.30.0) on Fedora 27 so it might be that some instructions only work on Linux but not on MacOS or MS Windows Mashines.

**NOTE:**
If `minishift version` shows something smaller than `1.30.0`, please
consider running `minishift update` before opening an issue.

For a clean experience and setup, please follow these steps to bring up
minishift and prepare get the OpenShift `oc` command.

```bash
minishift stop
# Optionally:
#   minishift delete --force
minishift \
    --profile fabric8-streams \
    start \
        --memory=8192 \
        --cpus=4 \
        --disk-size=10g \
        --vm-driver=kvm

# Use "oc" binary from minishift
eval $(minishift oc-env)
```

## Depoloying RH AMQ Streams in your OpenShift cluster

Throughout this guide you'll find that we're referring to an OpenShift project
that goes by the name `$YOURPROJECT`. In order to follow this guide, please
export this environment variable before you continue:

```bash
# Replace the project name "mynewproject" with i.e. "dsaas-preview" or
# "dsaas-production" depending on where you want to deploy.
export YOURPROJECT=mynewproject
```

### Create or use OpenShift project

Optionally login as a user (e.g. "developer") with the right to create a new
project or if already done or just choose the right project to run on.

```bash
# Optional:
oc login $(minishift ip):8443 -u developer -p developer
# Optional:
oc new-project $YOURPROJECT
```

Switch to an OpenShift user with `cluster-admin` role (e.g. user
`system:admin`):

```bash
# Login as a user with "cluster-admin" role needs to be used,
# for example, "system:admin".
oc login -u system:admin
# Choose "yourproject" to run on
oc project $YOURPROJECT
```

Deploy the Cluster Operator to OpenShift ([documentation](https://access.redhat.com/documentation/en-us/red_hat_amq/7.2/html/using_amq_streams_on_openshift_container_platform/getting-started-str#deploying-cluster-operator-openshift-str)):

```bash
sed -i "s/namespace: .*/namespace: $YOURPROJECT/" cluster-operator/*RoleBinding*.yaml
oc apply -f cluster-operator -n $YOURPROJECT
```

**WARNING:** I suggest that you run `watch oc get pods -n preview` after each
`oc apply` to wait until the listed pods have the `STATUS` called `Running`.

Create templates to build upon when deploying the Kafka resources

```bash
oc apply -f templates/cluster-operator -n $YOURPROJECT
```

Deploy persistent Kafka cluster to OpenShift ([documentation](https://access.redhat.com/documentation/en-us/red_hat_amq/7.2/html/using_amq_streams_on_openshift_container_platform/getting-started-str#deploying-kafka-cluster-openshift-str)):

```bash
oc apply -f kafka/kafka-persistent.yaml
```

### Test deployment so far

Open two terminals, one for a *producer* and one for a *consumer*.

In the *producer* terminal run this command, and type in a `Hello World!` and hit `<ENTER>`:

```bash
# Optionally: eval $(minishift oc-env)
oc run kafka-producer \
    -ti \
    --image=registry.access.redhat.com/amqstreams-1/amqstreams10-kafka-openshift:1.0.0 \
    --rm=true \
    --restart=Never \
    -- bin/kafka-console-producer.sh \
        --broker-list fabric8-streams-cluster-kafka-bootstrap:9092 \
        --topic my-topic
```

In the *consumer* terminal run this command to receive all messages created by the producer:

```bash
# Optionally: eval $(minishift oc-env)
oc run kafka-consumer \
    -ti \
    --image=registry.access.redhat.com/amqstreams-1/amqstreams10-kafka-openshift:1.0.0 \
    --rm=true \
    --restart=Never \
    -- bin/kafka-console-consumer.sh \
        --bootstrap-server fabric8-streams-cluster-kafka-bootstrap:9092 \
        --topic my-topic \
        --from-beginning
```

If the test was successful, you can close the terminals. 

# Examples for producers and consumers

## Deploy OpenShift Container Registry

For further testing with custom built images, we [deploy a container image registry](https://docs.openshift.com/container-platform/3.11/install_config/registry/deploy_registry_existing_clusters.html) (aka the OpenShift Container Registry) using a user with cluster admin priviliges (e.g. `system:admin`):

```bash
oc login -u system:admin
oc adm registry --service-account=registry
```

The created pod uses an ephemeral volume that is destroyed if the pod exits. This perfectly fits our testing puroposes.

We're going to use the docker daemon as it is deployed inside OpenShift and we're loggin into the registry as user `developer`.

```bash
eval $(minishift docker-env)
oc login -u developer -p developer
docker login -u developer -p $(oc whoami -t) $(minishift openshift registry)
```

## Build and deploy a producer written in Go

**NOTE:** This section is not finished as is the `producer-go`

Now, we create an image stream for our example go producer, build and then push the image to our OpenShift container registry.

```bash
oc apply -f examples/producer-go/image-stream.yaml
docker build \
    examples/producer-go/ \
    -t $(minishift openshift registry)/$YOURPROJECT/producer-go:latest
docker push $(minishift openshift registry)/$YOURPROJECT/producer-go:latest
```

Next, let's create a deployment config that tells OpenShift to deploy a new version if it is pushed to the image stream that we've created before.

```bash
oc apply -f examples/producer-go/deployment-config.yaml
```






... To be continued ...