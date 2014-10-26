### Server 

This is a brief explanation of how the server works. It's all based on a Martini server. The martini 
server uses something called the SRPCServer to handle socket-based RPC.

##### SRPCServer

This is the head of the socket RPC system. It is an http handler, first of all, and listens for 
incoming ws connections on /srpc. When a new websocket connection is created, the server will wrap 
it in a 'SessionHandler' and add it to a list of active sessions.

##### SessionHandler

The session handler keeps a reference to the websocket connection, and a list of SRPCService 
objects. It handles incoming and outgoing messages over the socket.

##### SRPCService

The server may allow different types of RPC services. A service is basically just a collection 
of related methods. The EthereumSRPC service has methods that deals with Ethereum. An IPFS service 
would deal with IPFS related calls, and so forth. Every service uses the same basic object for requests 
(the RPCRequest object), and also for responses (RPCResponse). 


```go
type RPCRequest struct {
	Method string 
	Params *json.RawMessage
	Timestamp int 
}

type RPCResponse struct {
	Id string 
	Result interface{}
	Error interface{} 
	Timestamp int
}
```


A normal call from a client is on the form: Method = "ServiceName.MethodName". If the client 
wants to get the balance of a certain ethereum account, for example, it would be:

``` 
Method: "EthereumSRPC.BalanceAt"
Params: {SVal : "accountAddressAsHexString"}
Timestamp: (millisecond time)
```  

These calls are mapped to js functions though, that automatically creates the request object, 
JSON.stringify them and send them to the server.

When the server receives this request, it will be a SessionHandler that gets it. The session-handler 
will choose the SRPCService based on the first part of the method, "EthereumSRPC", and the method 
it calls from EthereumSRPC is of course "BalanceAt".

##### SRPCServiceFactory

When a the SRPCServer successfully negotiates a new websocket connection, and creates a new 
SessionHandler for it, the new session-handler needs its own set of services. It cannot 
share its services with other sessions. The way this is done is by registering a SRPCServiceFactory 
with the SRPCServer (before it starts to listen for connections). Each session type needs its own 
factory. When a new session-handler is created the server will iterate over its map of factories, and 
create a session of each type, and then add them to the new handlers map of services.

Sometimes services needs to share parts of their state with other services (of the same type). 
EthereumSRPC is an example of this; it has its own unique pipe to communicate with the eth client, 
but all pipes lead to the same ethereum. This is managed by the factory. It keeps a reference to 
ethereum that it passes along to each new service instance it creates.  

##### Connection Management

SRPCServer has a 'maxConnections' field, and an 'activeConnections field. It will not accept any 
more connections when it reaches capacity. If it isn't at capacity, it will create the websocket 
connection, create a new SessionHandler, up the number of active connections by one, give the handler 
an id, and add it to its list of active sessions. The id is provided from a pool. When a connection 
is closed it will subtract one from the number of active connections, reclaim the id, and remove the 
handler from the list.  