# Atë Scripting System (aka, the Eris Stack Scripting Language)

## The Beginnings

### c3d

The Ate system has its roots in a specification that was developed in the Spring of 2014 by the Project Douglas team. That specification was purpose built to allow rationalized, efficient, and secure communciation between a smart contracts enabled blockchain, a distributed file store (or localized cache), and a user interface. This system was named `c3d` which stands for "Contract Controlled Content Distribution".

Contracts which are compliant with the c3d system generally follow a tree model. Each contract in the system is a node, and certain fields are set with a value so that reading the node contracts can be performed in a rational, programmatic way to see, for example, if the node was a leaf or a branch of the tree. For non leaves (branches), standard storage locations are used as an entry point so that the (linked) list of child branches and leaves can be iterated over. The c3d specification also enabled other data to be associated with each contract/node. For example, fields that held hashes associated with particular files which could hold content associated with the contract node (if the contract was tied to any out of chain data), references to UI elements, among other things.

A simple example of this model would be this: A tree with a contract A as the root. Its two children are contracts B and C, and C has a child D. Let's say the type of a contract is put at address 0x10, and the entry point to the list of children is at 0x11. If I wanted to get a summary of A, I would do this:

1. Check if A is a leaf, by reading 0x11.
2. If A is not a leaf (as in this example), I'd get the address to the first child, which would be stored at 0x11. If that address is 0 that means no children have been added. Otherwise I'd move on to the next element (normally it would be stored at the current list element address + 2). I would keep going until the next element is 0, and store all the child contract addresses as I go.

c3d uses a few other fields as well, let's say that each list element is stored on this (reduced) format:

```
childId -> contract address
childId + 1 -> address to previous child
childId + 2 -> address to next child
childId + 3 -> id of associated content file
```

This would be understood by any program which had been given an understanding of the specification, so in that case I could get the filename as well as the contract address to get the full summary.

Now, lets say I wanted to get a summary of the entire tree. Instead of stopping with Contract A, I would keep working my way through the tree of contracts. I would begin by picking the first of A's children. That would be contract B. Now I can do the same thing here, check the type address, and determine that it is a leaf. For the second child, C, I would see that it is a node rather than a leaf, and by continuing the iteration process I would get the address to the one child of contract C (which is contract D). Contract D is a leaf. By storing all the data for each child of each node, I end up with a full summary of the tree.

### The Tree Parser

This system has been implemented, in the form of a "decentralized reddit". Using the notation of the above examples, Contract A would have been the main site, its children would be forum contracts (in Reddit parlance, these would roughly analogize to subreddits), and its children's children would be the threads in that forum (or subreddit), then the children's children's children would be posts in those threads.

When someone opens up the web UI, they would get a list of all forums. Then they could click the forum and get the threads, working their way through the tree structure. The way that works was through the tree parser. The tree parser could get a series of ids, ("forumId","threadId", etc.). After this the tree parser would use those ids to navigate its way through the tree. The length of the path above and below the contract in question would allow the tree parser to determine if contract in question was a forum, or a thread, or a post. In turn the tree parser was enabled in such a way as to pass the location within the tree heirarchy to the user interface so that the rendered interface would modify itself appropriately based on where in the tree structure the rendered content was located.

The tree parser would start by checking the root contract for a child with id "forumId". If it finds the forum, then it moves on to search its children for the thread, etc.

What made this a decentralized reddit was the fact that every person has a copy of the blockchain (which was used as a decentralized data store), and because the content within the forum posts were shared as torrents everyone had the off-chain content data as well. If someone didn't have any of the data needed to use the system, it would be provided to them by the other peers rather than a server. And also, it worked like Reddit because the contract nodes and leaves allowed for up and down voting of the content to surface content which was interesting or important to the group using the system.

### Action Models

Action models are higher level representations of contracts that allow you to do batched reads/writes with one single method.
When the tree parser has found the correct contract for you, you could get a reference to the action model file from that
contract, and load the action model (a .json file). The model would have standard operations, like for example editing a forum
post. When called, the function (or 'action') would perform all the low level stuff, check certain addresses to make sure
everything is fine, do calls to the contract, maybe follow up on that call to see that things were done right before returning
the result of the call. It all worked through a basic scripting language that let you use certain variables.

Normally, the action model would be accessible through a contract as part of the c3d spec, most often as a reference to a file
with a certain name and a given SHA digest (for security).

Here's an example of a basic action model. It could be used to comment on a blog entry, for example.

BlogEntry
```
{
  "actions": {

    "addComment": {
      "precall": [ "$this:0x25" ],
      "call": [ "$this", $params[0], $params[1] ],
      "postcall": ["$this:0x25" ],
      "success": "$postcall[0] > $precall[0]",
      "result": "$postcall[0]"
    },
  },

  "data": {
  	something something
  }
}
```

Every action model has this format - a series of different actions (only 1 in this case), and some data. Actions had a sequence:

precall->call->postcall->success->result

This is the same for all actions, but some of the parts could be empty.

Precall was done before any calls was made. In this case, what it does is to get the value stored at address 0x25 in '$this',
which we can pretend is the number of already existing comments. Identifiers such as $this (the address of the contract being run)
and $params (the parameters that was sent along with the call) are added to the action model parser before running the
action itself.

Call makes a call to the given contract. The first value in the array is the contract to call (in this case it's '$this'). The
next few values are the tx payload.

Postcall is done after the transaction is made. We get the number of comments again.

Success is a conditional statement that determines if the action was successful. In this case we check that the number of
comments did in fact increase, which would mean that the transaction passed. It's a weak condition but this is just an example.

Result is the data sent back. Let's say each comment gets its own contract, then we could have used result to perhaps pass the
address to the new contract back to the caller.

A more detailed spec can be found here: https://github.com/eris-ltd/deCerver/blob/master/docs/models%20Specification.md

## Atë

Atë is an extension of this system. It still works in much the same way, but some things have been changed. Action models
may now use javascript (not just the object notation). This gives a lot more freedom. Also the c3d model itself is
being updated, as well as the data types and structures that the contracts use.

A javascript action model could look like this:

```javascript
// Genesis DOUG
function Model(address){

	var _address = address;
	var _name = "GenDoug";

	var _data = {
		// C3D
		"BAindicator" 	: "0x10",
		"BAdmpointer" 	: "0x11",
		"BAUIpointer" 	: "0x12",
		"BAblob" 		: "0x13",
		"BAparent"		: "0x14",
		"BAowner"  		: "0x15",
		"BAcreator"		: "0x16",
		"BAtime"		: "0x17",
		"BAbehaviour" 	: "0x18",
		"BALLstart"		: "0x19"
	};

	// GenDOUG
	function nextslot(params) {
		return Add(params,"2");
	};

	function prevslot(params) {
		return Add(params,"1");
	},

	function typeslot(params) {
		return Add(params,"3");
	};

	function behaviourslot(params) {
		return Add(params,"4");
	};

	function dataslot(params) {
		return Add(params,"5");
	};

	function modelslot(params) {
		return Add(params,"6");
	};

	function UIslot(params) {
		return Add(params,"7");
	};

	function timeslot(params) {
		return Add(params,"8");
	};

	function permname(params) {
		return Add(params, _data.offset)
	};

};


GenDoug.prototype.getname = function(params) {
	if(typeof params !== "string"){
		console.log("GenDoug.getname: Params should be a string.");
		return null;
	}
	// Make sure the name is an actual name.
	if (!IsZero(Mod(params,_data.GFzeros) ) ){
		return null;
	}
	return GetStorageAt(_address,params);
};


//Setters (these require calls to be made, and will halt execution until a block is mined)
GenDoug.prototype.register = function(params) {
	// TODO Decide how to handle errors and return them. Since we can't get return values in a
	// simple way yet (?), and at this point 'register' only returns 0 or 1, I add
	// some stuff here for now.
	if(!(params instanceof Array) || params.length != 2){
		console.log("GenDoug.register: Param length != 2 (should be: ['regname', '0xaddress'] (string,string)");
		return false;
	}
	if(typeof params[0] !== "string" || params[1] !== "string"){
		console.log("GenDoug.register: Param length != 2 (should be: ['regname', '0xaddress'] (string,string)");
		return false;
	}
	Transact(_address,["register"].concat(params));
	return true;
};
```

It is a javascript "class" that can have public and private variables and functions. Functions like Add and IsZero
is added to the javascript vm upon initialization. This makes it possible to do more advanced things, loops, all kinds
of arithmetic and branching etc.

### Calling actions from actions

It is easy to run one action from inside another. The actions are just functions in the same "class".
To call an action from an action inside another model, the model has to somehow get the other model.
This can be done in one of two ways. If the model is "above" the calling model in the tree, it can
get there by using its parent address (every contract has a reference to their parent). If the target
model is "below" the calling model in the tree, then it would get the appropriate child, etc. It can be
complicated and slow, depending on how the call is made, but it works.

### Otto

The library that is used to run javascript is Otto. It has been developed for about 2 years (which is a pretty long
time for a Go library), it is used in the Go implementation of Ethereum, and has a lot of great features. It can not
only interpret script, it can also pre-compile script into "Script" objects that can be run later, and even parse script
and return the AST (making action models -> contract code generation easier if we want to do that at some point).

The system currently keeps one single otto instance that is used for all script, but that will likely be changed.

### Running script & Security

There is no way to just run a snippet of javascript, as in passing a string to otto and run it. What you can run is:

1) Script that is held in a file and has been pre-added by someone with the persmissions to do so.
The SHA digest of the files are stored in a contract, so when the decerver boots up and loads the
files, it can check that the script files are valid. The people that are allowed to add script should
probably be the same people that are allowed to add contracts, eg those with root access.

2) Functions that are defined in those pre-approved files. When doing so, you don't just pass a js
formatted string containing a function name and parameters, all you can do is choose which file and
which function (often via the tree parser), and pass parameters. The parameters are then processed
(only basic types, and arrays of basic types are allowed) and passed into the vm inside a pre-rendered
function call. Often these calls will in turn be accessible only through RPC services or some other
components inside deCerver.

3) Functions that are (manually) exposed to otto by deCerver. Some example functions would be getting
stuff from contract storage, and making transactions. Later there will be bindings to the eris file
system as well. There are also another class of "injected" functions, such as math helpers. Note that
you cannot even run these functions directly. You can't call the exposed "GetStorageAt", unless it is
exposed to you in turn by a function within one of the pre-approved javascript files.

### Action models & tree parsing.

Tree parsing can be done in one of two ways - either by having a go function doing it, or by having
the tree parser run inside the javascript vm. Currently I run the parser from within the vm, because
it has all the functions and action models that the tree parser needs. Changing it would be very easy
because the parser itself is only made up of a few simple functions. If this is supposed to be run from
Go, then the alternatives are to either remove the dependency to action models, or have the parser darting
in and out of the vm, passing values back and forth (ie copying) to be able to run both action model
script and its own routines.

The alternative would be to have the tree parser run in go, and use hard coded addresses for
c3d stuff (meaning every contract has to have the same address for child list entry point etc).
That way it could parse the entire tree, get the action model and contract address at the
end and pass those to the vm to do the actual action. Currently it gets the action model for
each node in the tree, and uses its standard c3d functions (getchild etc.) to get the children.
This makes it possible to set the c3d addresses on a per contract basis.