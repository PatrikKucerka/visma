# Image Service

Image pre-processing is an important task in many machine learning applications, but can be an intensive task when trying to keep request times low and throughput high. In the VML team we use a micro-service architecture running on Kubernetes, along with autoscaling implementations to achieve image processing systems that scale well.

For this task we want to implement a small service that scales and grayscales incoming images to some normalized settings, so they could be used as input for a machine learning system.

We also want to containerize the service and deploy it to Minikube.

## The sample code
We have provided the code for a simple HTTP API for the image service, along with YAML manifests for deploying the API to a Kubernetes cluster. The code currently just echos the image sent in the request, no pre-processing done yet.

## Prequisites
Install: 
- [Golang](https://go.dev/dl/)
- [Docker](https://docs.docker.com/engine/install/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start)

## Local testing & setup
1. To start the server run `go run cmd/service/v1/main.go`
2. To call the server run `go run cmd/client/v1/main.go ./test.jpg ./output.jpg` 

### Kubernetes setup
1. Run `minikube start` to start up the cluster. 
   1. Test the connection to the cluster with `minikube kubectl -- get pods -A` 
2. Run `eval $(minikube docker-env)` to setup docker environment
   1. this needs to be done in every session of the terminal

### Deploying to Kubernetes
1. First we need to build docker image `docker build . -t image-api:latest --file dockerfiles/Dockerfile.api`
2. Then we need to apply the API manifest to the Minikube `minikube kubectl -- apply -f manifests/api.yaml`
3. We can create a tunnel to the Minikube by port forwarding from localhost to the Minikube `minikube kubectl -- port-forward image-api 8080:8080`
4. Now we can call the Minikube as if we were calling the localhost `go run cmd/client/v1/main.go ./test.jpg ./output.jpg` 

## Your Task

### Part One - Programming
Create you own micro-service to do the image processing, seperate from the API service. You can choose to write it in whatever programming language you are comfortable with. We expect it to:

* Be a simple HTTP server
* Take as part of the request:
    * The image originally sent to the API
    * Parameters for scaling the image or not
    * Parameters for grayscaling the image or not
    * The service needs to be able to handle JPEGs and PNGs

Apart from that you will need to change [the API endpoint](./pkg/api/api.go) so that it calls and passes the image to your service. The API should satisfy the following:

* Be configurable, such that the desired image size and grayscale options should be supplied at startup. How is up to you. You can use 1024x768 as the target size.

### Part Two - Kubernetes
Create Dockerfile for containerization of your service as well as manifest to deploy your new service to the Minikube. Take an inspiration from already existing manifest and dockerfile for the API.

Deploy your new service onto Minikube. The API should be able to call your service within Minikube. You should be able to call the API through port forward. Useful link: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/ 

### Part Three - Optional extras / Questions to think about
#### Extras
- Have an option to supply the API with a URL for downloading images, instead of image bytes. Whether you want to download the image content in the API and send bytes to your service - or send the URL and download in the service - is up to you.
- Implement caching of the images so that we don't do processing of the images twice. It's up to you where and how.

#### Questions
- How would you horizontally scale your service? 
- If we wanted to store the images for further processing - how and where would you store them?
- What would you change to support rasterization of PDFs?
