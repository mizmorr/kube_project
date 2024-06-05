# About
Project. Service mesh architecture

# Intallation
Check if [Istio](https://istio.io/latest/docs/setup/getting-started/), Docker and [Minikube](https://minikube.sigs.k8s.io/docs/start/?arch=%2Flinux%2Fx86-64%2Fstable%2Fbinary+download) are installed.

Then u need to execute the following

    kubectl apply -f referencer/
    kubectl apply -f requester/
    sh system/cluster_apply.sh

Well done!

# Getting Started

Run this command in a new terminal window to start a Minikube tunnel that sends traffic to your Istio Ingress Gateway. This will provide an external load balancer, EXTERNAL-IP, for service/istio-ingressgateway.

    minikube tunnel
Now, we should export current URL of gateway to variable:

    export GATEWAY_URL=$( sh system/gw.sh )
Ensure an IP address and port were successfully assigned to the environment variable:

    echo $GATEWAY_URL
Install Kiali and the other addons and wait for them to be deployed.
Change directory to Istio's. Then run:

    kubectl apply -f samples/addons
    kubectl rollout status deployment/kiali -n istio-system

Access the Kiali dashboard. Execute in another terminal:
    
    istioctl dashboard kiali

Go make requests! You can try just curl it:

    curl http://$GATEWAY_URL/hi_every_one/

Or use special script:

    sh system/system_test.sh

Look at dashboard!

# Clean up 

If you want to clear the kubernetes cluster (default namespace only), then execute folowing:

    sh system/delete.sh
