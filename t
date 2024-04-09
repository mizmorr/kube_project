apiVersion: v1
kind: Service
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"app":"listener","service":"listener"},"name":"listener","namespace":"default"},"spec":{"ports":[{"name":"http","port":9080}],"selector":{"app":"listener"},"type":"NodePort"}}
  creationTimestamp: "2024-04-04T20:30:22Z"
  labels:
    app: listener
    service: listener
  name: listener
  namespace: default
  resourceVersion: "515443"
  uid: 662205fe-97a2-4abc-80e7-17b13deb681c
spec:
  clusterIP: 10.111.20.85
  clusterIPs:
  - 10.111.20.85
  externalTrafficPolicy: Cluster
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - name: http
    nodePort: 32162
    port: 9080
    protocol: TCP
    targetPort: 9080
  selector:
    app: listener
  sessionAffinity: None
  type: NodePort
status:
  loadBalancer: {}
