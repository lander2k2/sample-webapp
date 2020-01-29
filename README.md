# Sample Web App

A one-file go sample application that connects to a postgres database for
testing, demos and examples in Kubernetes.

## Prerequisites

- a running Postgres database called "webapp_sample"
- psql installed locally

## Prepare the Database

If your database is in a Kubernetes cluster you can use the following to forward
a port from the database:

    kubectl port-forward -n web-app pod/<database pod name> 6432:5432

Use psql to create the database schema:

    psql -U webapp_sample_appuser -p 6432 -h localhost -d webapp_sample -a -f
    db/v1.0_schema.sql

Add some mock data:

    psql -U webapp_sample_appuser -p 6432 -h localhost -d webapp_sample -a -f
    db/planets-mock.sql

## Run Sample App in Kubernetes Cluster

Note: The manifests used here use a namespace called "web-app".

Create secret for database credentials:

    export DB_USERNAME=$(echo "<db username>" | base64)
    export DB_PASSWORD=$(echo "<db user password>" | base64)
    cat secret.yaml | envsubst | kubectl create -f -

Run the pod:

    kubectl apply -f k8s.yaml

Forward a port to your local machine:

    kubectl port-forward -n web-app pod/sample-webapp 8080:8000

Now you can view the sample app at http://localhost:8080/.

