# Image Service

Image pre-processing is an important task in many machine learning applications, but can be an intensive task when trying to keep request times low and throughput high. In the VML team we use a micro-service architecture running on Kubernetes, along with autoscaling implementations to achieve image processing systems that scale well.

For this task, we'd like you to implement a small service that scales and grayscales incoming images to some normalized settings, so they could be used as input for a machine learning system.

## The sample code
We have provided the code for a simple HTTP API for the image service, along with YAML manifests for deploying the API to a Kubernetes cluster. The code currently just echoes the image sent in the request with no pre-processing done yet.

## Prequisites
Install: 
- [Golang](https://go.dev/dl/)

## Local testing & setup
1. To start the server, run `go run cmd/service/v1/main.go`
2. To call the server, run `go run cmd/client/v1/main.go ./test.jpg ./output.jpg` 

## Your Task

### Part One - Programming
Create your own separate micro-service to do the image processing. We would like you to write the code in Go. We expect your new service to:

* Be able to receive requests with the following data:
    * The image originally sent to the API
    * Parameters that tell the service whether to scale the image or not
    * Parameters that tell the service whether to grayscale the image or not
    * The service needs to be able to handle JPEGs and PNGs
* Currently, the services communicate using standard HTTP setup. We would like you to change the code so the services use [gRPC](https://grpc.io/docs/languages/go/quickstart/) to exchange data between each other.
   * You will need to create a `.proto` file and evantually generate a Go code for the service to be able to use gRPC. The `protoc` command will help you with the code generation. To make things easier, you can just create the `.proto` and the generated Go file in the root directory of this task.

Apart from that, you will need to change [the API endpoint](./pkg/api/api.go) so that it calls and passes the image to your newly-created service. The API service should be configurable, so that the desired image size and grayscaling options (which will be sent to your new service) should be supplied at the API's startup. How you do this is up to you. You can use 1024x768 as the default size.

### Part Two - Optional extras / Questions to think about
#### Extras
- Have an option to supply the API with a URL for downloading images instead of providing image bytes. Whether you want to download the image content in the API and send the bytes to your service, or send the URL and download the byte content in your service, is up to you.
- Create Dockerfile for containerization of your service, as well as manifest to deploy your new service to kubernetes cluster. You can take inspiration from the already-existing manifest and dockerfile for the API. 
   - We only ask you to do this optional step on the theoretical level, but if you have the time and want to test whether your docker image and kubernetes manifest actually work, you can test it by deploying the services to Minikube (see the [last section](#optional-trying-out-the-code-in-kubernetes-with-minikube) of the readme). If you want to call the API in the Minikube cluster, you can use port-forwarding (here's a useful [link](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/)). 
- Implement caching of the images so that we don't do processing of the same image twice. It's up to you where you do the caching and how.

#### Questions
- If we wanted to store the images for further processing - how and where would you store them?
- What would you change to support rasterization of PDFs?
- What do you think are the advantages (and disadvantages) of using gRPC for inter-service communication over regular HTTP communication?

## Optional: Trying out the code in kubernetes with Minikube

### What you'll need
- [Docker](https://docs.docker.com/engine/install/) (optional)
- [Minikube](https://minikube.sigs.k8s.io/docs/start) (optional)

### Kubernetes setup
1. Run `minikube start` to start up the cluster. 
   - Test the connection to the cluster with `minikube kubectl -- get pods -A` 
2. Run `eval $(minikube docker-env)` to setup docker environment
   - this needs to be done in every session of the terminal

### Deploying to Kubernetes
1. First we need to build the docker image, do this by running `docker build . -t image-api:latest --file dockerfiles/Dockerfile.api`
2. Then we need to apply the API manifest to Minikube `minikube kubectl -- apply -f manifests/api.yaml`
3. We can create a tunnel to the Minikube by port forwarding from localhost to the Minikube `minikube kubectl -- port-forward image-api 8080:8080`
4. Now we can call the Minikube as if we were calling the localhost `go run cmd/client/v1/main.go ./test.jpg ./output.jpg` 
