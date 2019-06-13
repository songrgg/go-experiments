## Golang Experiments
Do some experiments to find a more proficient way to write Golang.

### Start
Golang 1.11+ is needed.

### Questions
1. How to initialize gRPC client for better throughout?  
Single connection for gRPC client will have a high throughout, it's recommended.  
Pros:
    * gRPC use HTTP/2 as transport protocol, HTTP/2 leverages multiplexing & connection reuse & package compression etc, 
    it will decrease the connection creation cost and improve the transport efficiency.
    * gRPC client will handle the disconnect/reconnect situations

   Cons:
    * LoadBalancer need to be implemented when the destination is service IP mapping to service cluster, which is common
    when Kubernetes Service object is used.
    
    [gRPC connection test](grpc/clientconn_test)

1. So many log libraries, choose which one?  
TODO
