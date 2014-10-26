## Overview

C3D compatible contracts integrate a datamodel for individual contracts. For factories which are distributed with the Eris platform, the datamodel is integrated at the factory level. This is because a factory will produce a given type of smart contract which will, in turn, have a specific data model. Other contracts as well will have their data models and UI files determined at the time they are deployed. Each c3D compatible contract will have a reference hash to a JSON compatible file in the `0x11` storage space of the contract (or this slot will be blank if there is not an associated model). This reference hash may be acquired via torrents, IPFS, or any other file sharing system which is c3D compatible.

The data model which the deCerver utilizes is contained in a JSON file. When a user wishes to interact with a smart contract, the deCerver will acquire the datamodel, read the JSON into memory, parse the API call from the user, and then perform the sequence of queries and transactions as mandated by the model. The model file **must** contain valid JSON or the deCerver will not be able to parse it.

Thus, the roll of the datamodel is to assemble the action calls and data required to render the view which the user can interact with. The top level of the JSON object contains two fields: actions and data. Each of these fields are covered in more detail below.

## Action Calls

Actions are given a name (a JSON string) which will be matched by the deCerver in order to perform the action. These actions can be anything which is allowed by the smart contract.

### Variables

Action calls will require variables in order to be viable in complex systems. Variables within the action models are denoted with the `$` symbol. The action models are able to reference the following variables:

* `$gendoug` refers to the address of the genesis DOUG for the chain in question. The deCerver will have the address for the genesis DOUG in memory, or will be able to find it quickly.
* `$this` refers to the contract in question. The `$this` variable will be replaced by the address of the contract in question throughout the action sequence.
* `$param[param_name]` refers to paramaters which are passed to the deCerver call. When the API call is sent to the deCerver the `param_name` will be a string and the value of the `param_name` will be substituted for `$param[param_name]` in the action sequence outlined below.
* `$precall[n]` refers to the results of the precall queries which have been sent to the blockchain (see below for more information). When subsequent queries or transactions are sent to the blockchain the return value of the query outlined at `$precall[n]` will be substitued. These values may only be referenced *after* they have been made. In other words you cannot reference `$precall[4]` in the first call of the $precall array.
* `$postcall[n]` works exactly like the $precall array.
* `$blob` refers to the hash of the file blob. This is a file blob pointer which can be used to send to the blockchain.

### Query Structure

The queries sent to the blockchain require an address and a storage location. Throughout the datamodel these are both included in one JSON string and denoted by a ":" between them.

Storage locations may begin either with a '0x' and continue with the hex storage location or may be a normal string (e.g., for a namereg type of contract) and the deCerver will translate the string into the hex required.

So each of the following are all correct:

```json
"precall": [ "$this:0x11", "$gendoug:Ponos", "0xd6e96ee6661367735c15894193bdca276bae27ba:feedface00000000000000000001000" ]
```

### Action Sequence

To make an action call, the deCerver will process the action model in a sequential five-step process. Each of these steps will allow the deCerver to (1) acquire sufficient information from the blockchain and smart contracts in order to formulate the necessary transaction in a meaningful way, (2) to test whether the transaction was successfully completed, and (3) providing a meaningful return to UIs which utilize the actionmodel capability of Eris derived smart contracts.

Within each action there are five fields which the deCerver will parse.

1. `precall`: these are queries to the blockchain which may be required previous to sending the call (transaction). Precall queries are passed as a JSON array and the deCerver will query each of the assembled queries in order. The result of each call will be stored in the `$precall` array for use by any queries or fields which follow the individual query. For example if the first precall query was to ask the Genesis DOUG what the value of a particular key was then the call (or any later query within the precall or field) will be able to access the result of that query by referencing `$precall[0]`.
2. `call`: call is the transaction which is assembled by the deCerver and sent to the blockchain. Call transactions are passed as a JSON array to the deCerver. If multiple transactions are necessary to perform an action then the call object will be an array of arrays. In the multiple transaction scenario, each element of the top level array will be a single transaction array formatted in the same way as a single transaction call is formated. The first element of any single transaction array *must* be the address of the destination contract of the transaction. The remainder of the single transaction array will be the data which is passed to the destination address. Call has access to any of the `$precall` array values via the `$precall` variable. Since transactions do not (at this time) return any value there is nothing which the `call` field will store.
3. `postcall`: these are queries to the blockchain which may be required after sending the call (transaction). All of the same rules as for precall exist here, save for the results of `$postcall` queries are stored in the `$postcall` array rather than the `$precall` array.
4. `success`: actions either succeed or fail. At this time there is no other answer which may be returned to the UI besides success=true or success=false. The UI will have to determine what that means and how to react. Success is defined as a comparison operator (`==`, `!=`, `>`, `>=`, `<`, `<=`) matched against two of the array values from either of the `$precall` or `$postcall` arrays. The other variables for the action model are valid to utilize in the success call but it is unlikely to be used with any practicality.
5. `result`: the result is returned, along with the result of the `success` field to the UI at the completion of the action sequence. It is the `return` in most programming languages. The `result` field must be either a comparison operator (used the same as the `success` field) matched against two of the array values from the `$precall` or `$postcall` arrays OR a single value from one of the arrays. The result field is passed as an array and will return a JSON array with the individual results assembled. Numerous `result` values may be returned to the UI by assembling them into a JSON array in the action model.

In addition, each of the functions has access to the addition and subtraction of hex numbers. Usually this is necessary in order to traverse within groups of c3d content blobs -- in particular in `BA` contracts. To add a number to a base, use something like `BASE+5` where BASE was something like `$postcall[1]`.

**Note**, this will ONLY work to modify the `$precall` and `$postcall` arrays and will not work for general calls. In addition, this will only work with storage locations and not with contract addresses. So this **will not** work: `$this:0x19+5`, but this **will** work: `$this:$postcall[0]+5`.

### Returns of Action Calls

The UI will receive a JSON formatted string from the deCerver after the result of the action call which will be formatted like so:

```json
{
  "success":true, // or...false
  "result": [
              "0xaaaaaaaaaaaaa" // or... whatever the result is.
            ]
}
```

###  Notes

When indicating to the datamodel a storage location it is **required** to use the '0x' prefix.

If you do not, then the deCerver will build the storage address according to a hex construction of the string. So if you indicate '0x19' then the deCerver will check the `0x19` slot.

If you indicate '19' then the deCerver will check the `0x3139000000000000000000000000000000000000000000000000000000000000` slot (which is the hex construction of the string which includes '1' and '9').