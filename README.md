## The deCerver

This is the decerver repository.

## What is it?

It is beta software.

It is an application platform. It lets you load and run applications. We call them distributed applications because those are the types of applications we expect the decerver will mainly be used for. These applications consist of a UI (http/css/js), back-end javascript (this would be the server-side scripts in a normal web application), an optional smart contract back end, and some application meta data such as package and configuration files. 

It is node.js but in Go. The way you call decerver methods in your (distributed) application code is through special back-end javascript. Decerver exposes some basic functionality, such as networking, event-processing, big integer math, logging, and other things. It also exposes the API methods of its modules, because: 

It is a module platform, such as for example the Netbeans platform. It allows you to make a program or library into a module that can be utilized by the applications that are run. In order to make something into a module you need to do two things: 

1) You must give the library a module wrapper. The module interface requires a number of functions, such as Start, Init, Subscribe (for events), and Shutdown, and those should be bound to library functions or complemented in the module wrapper code if such functionality does not exist. If 'Start' does not make sense for a certain lib, then the wrappers Start method should do nothing, but it would still have to be included. 

2) You must give the library a javascript API. The way we do it in our modules is by simply calling the corresponding api method, and then convert the return value into a proper javascript object using some utility functions that are included with the decerver. It is very simple.  

In the decerver-modules repository, you can see how our modules are built. Legal markdown is by far the simplest one, and only has the basic module wrapper + a single API method. The others are more complex.

Currently, modules has to be compiled together with the decerver in order to be used, which makes it look as if  Thelonious, IFPS and the other modules we've included are part of decerver itself, but they are not. IPFS, for example, is a stand alone library. We just put a module wrapper on it and gave it a javascript API. It has the same module wrapper as all the other modules we use. Decerver and dapp programmers could do this with any library they want, however, we currently do not support dynamic module loading... Once we get dynamic module loading in place (https://github.com/eris-ltd/decerver/issues/86), it will be possible to just switch modules out depending on what the (distributed) application needs. I would not recommend getting too deep into module development at this point, since the module API is going to change. Right now it's about application and not module making.


## What is it not?

It is not a blockchain. It does not have a blockchain do any of its work. Decerver itself runs perfectly fine without a blockchain present. It is however possible to use external libraries as modules, such as ethereum or thelonious, and those utilize blockchains as part of their functionality. These modules can then be used by the DApps that are run through the decerver.

It is not a key manager. Decerver does not store keys, it does not create them, and it does not sign any messages. If a module uses cryptographic keys, then it is also managing those keys. Central key management will be added to decerver, but only in the form of utility functions. The actual management itself will still be done externally.

It is not a program that is in competition with Ethereum. At least not Ethereum in its current form. It is a complement (see all the previous points). To be perfectly honest, though, the Eris system has its place in a distributed application stack. It is legitimate software that solves a lot of problems us early dapp makers faced in a good way. It needs and deserves to be around. If it happens to be in competition with other systems out there then too bad (for them mostly).

## Overview of the deCerver

![deCerver architecture](docs/images/deCerver Structure.png)

In the image shown above the deCerver package is in Blue. As stated, it acts as the hub for the system.

TODO