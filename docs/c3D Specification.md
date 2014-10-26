# The c3D Specification

#### Last Modified: 18 September, 2014

## Introduction to c3D

c3D compatible contracts have a standardized structure which is used by the c3D package to push and pull information to the contracts. By standardizing the ways in which individual contracts and systems of contracts are built, smart contract designers are able to build unparalleled extensibility into their systems.

c3D, or `contract-controlled-content-dissemination` is a method of disseminating content between nodes on a network in a manner whereby the *rendering* and the *dissemination* of this content is controlled by a system of smart contracts.

Using this specification, smart contract and decentralized application designers can develop any arbitrary system of content which can be disseminated to the proper people in a manner which is not *ad hoc* but rather is controlled in the ways and manners in which the smart contract designers desire. This arbitrary system need not have hard coded rules for its user interface, instead the user interface and their components will be able to understand relative positioning of content and other semantic criteria simply by their developers adhering to this specification.

Smart contracts wishing to have c3D-compatible clients, servers, and nodes disseminate their content must be built to the specifications of this document. If these designers feel this specification is useful, they should speak with their users to lobby for the continued integration of this specification into the entirety of the DApp ecosystem.

Similarly, clients, servers, and nodes wishing to *render* and/or *disseminate* content-heavy dapps and smart contract systems built to the specifications of this document must build parsers, renderers and disseminators capable of parsing smart contracts built to this specification.

Example disseminators, parsers and renders are currently available in:

* [Ruby](https://github.com/project-douglas/c3d); and
* [Go](https://github.com/eris-ltd/deCerver/c3d).

If you know of other renderers, parsers, or disseminators built in other languages please submit a pull request to this document.

The designers of this specification have attempted to simplify and harmonize the process of matching blobs of content with smart contracts so that it will be as open and widely utilized within the DApp community as possible.

If you do not find this specification useful, please do not hesitate to submit an issue and let your voice be heard.

## Blobs

c3D compatible content is stored in a file. These files can contain any arbitrary type or form of content. These files are called blobs throughout the remainder of this specification document.

c3D compatible smart contracts are capable of registering blobs of content which can be shared, distributed, disseminated, or served from any number of centralized or decentralized systems. By definition, there are no limits on filesize. Blobs **SHALL NOT BE** distributed with executable permissions within any c3D compatible system for security reasons.

Blobs must have a filename (also called, in bittorent parlance, a display name). This filename **SHALL BE** the SHA1 checksum hash when the blob is *first submitted* to the blockchain -- truncated to XXX bytes. If the blob is later updated the display name **MAY** be changed within the smart contract, **OR** it **MAY** stay with the same name as that used upon initiation.

## c3D Contract Types

*Notes*: in the below documentation storage slots can hold different types of data. These are the types of data used in the documentation:

* (A) indicates this field points to a Storage address
* (B) indicates this field should be filled with a blob_id value
* (C) indicates this field should be filled with a contract address
* (B||C) indicate either (B) or (C) is accepted
* (K) indicates this field should be a (public) Address
* (I) indicates this field is an indicator and can take one of the values in []
* (V) indicates a value

c3D compatible contracts come in three flavors, each of which is connoted in the 0x10 storage slot:

1. 0x10 : (I) : 0x88554646AA -- c3D Action Contract (no content attached)
2. 0x10 : (I) : 0x88554646AB -- c3D Content Details Contract
3. 0x10 : (I) : 0x88554646BA -- c3D Content Grouping Contract

The `0x10` slot of all c3D compatible contract **SHALL** contain one of the three values discussed above.

For the purposes of this documentation, `AA` contracts are simply indicators for the c3D system. c3D parsers will not take any action when passed an `AA` contract.

### BA Contracts -- c3D Content Grouping Contracts

`BA` Contracts are designed to provide parsers, renderers, and disseminators with the ability to see large groupings of contracts. As such they do not provide contract designers with reliable storage space which can be gauranteed to be collision free. As such they should use limited storage space outside of the c3D compatible slots. Storage heavy information should be pushed into an `AB` contract rather than be kept in a `BA` contract.

Nearly the entire storage space of an `BA` contract will consist of the top level slots with a running linked list of individual entries into the group. The following slots are required for each of the entries within the linked list as well as the top level slots.

#### Top Level Contract Storage Slots

The following slots are **REQUIRED** to be added to every `BA` contract. Information which is not needed, nor used, **MAY** be left blank.

* 0x10 : (I)    : [0x88554646BA]
* 0x11 : (B||C) : pointer to the contract's data model
* 0x12 : (B||C) : pointer to the contract's UI files
* 0x13 : (B||C) : pointer to the contract's content
* 0x14 : (C)    : Parent of this contract
* 0x15 : (K)    : Owner of this contract
* 0x16 : (C)    : Creator of this contract
* 0x17 : (V)    : TimeStamp this contract was created
* 0x18 : (I)    : Behaviour slot
* 0x19 : (A)    : Linked list start

The above should all be self-explanatory. The Behaviour slot is explained in greater detail below.

#### Individual Entity Entries

The following slots are **REQUIRED** to be added to every entry within every `BA` contract's linked list of individual entries. Information which is not needed, nor used, **MAY** be left blank.

Below the address of the first entry in the linked list is referenced as `linkID`.

* (linkID)+0 : (A)    : ContractTarget
* (linkID)+1 : (A)    : pointer to the previous entry in the linked list
* (linkID)+2 : (A)    : pointer to the next entry in the linked list
* (linkID)+3 : (I)    : Type slot
* (linkID)+4 : (V)    : Behaviour slot
* (linkID)+5 : (B||C) : pointer to the entry's content
* (linkID)+6 : (B||C) : pointer to the entry's data model
* (linkID)+7 : (B||C) : pointer to the entry's UI files
* (linkID)+8 : (V)    : Timestamp this entry in the linked list was added

The above should all be self-explanatory. The Type slot is explained in greater detail below. *Note* if the content is a pointer to an `AB` contract the data model pointers and UI pointers would typically be blank as it is much easier for parsers and renders to handle those when they are in the top level slots rather than within the linked list entries.

### AB Contracts -- c3D Content Details Contract

`AB` Contracts are meant to provide a large amount of storage space to smart contract designers who desire to add additional content to a blockchain smart contract. As such, these contracts are designed to be used with an individual piece of content rather than with groups of contents.

That said, any contracts which are not in direct violation of this specification will be parsed, rendered, and disseminated by c3D compatible parsers, renderers, and disseminaters without problem.

#### Top Level Contract Storage Slots

The following slots are **REQUIRED** to be added to every `AB` contract. Information which is not needed, nor used, **MAY** be left blank.

* 0x10 : (I)    : [0x88554646AB]
* 0x11 : (B||C) : pointer to the contract's data model
* 0x12 : (B||C) : pointer to the contract's UI files
* 0x13 : (B)    : pointer to the contract's content
* 0x14 : (C)    : Parent of this contract
* 0x15 : (K)    : Owner of this contract
* 0x16 : (C)    : Creator of this contract
* 0x17 : (V)    : TimeStamp this contract was created
* 0x18 : (I)    : Behaviour slot

#### Individual Entries for BA Contracts

*Note*: BA contracts will never have linked lists as they are predominantly used to track meta information regarding individual content blobs.

### Behaviour Slot

TODO - Explain

[0 => Ignore || 1 => Treat Normally || 2 => UI structure ||
                                        3 => Flash Notice || 4 => Datamodel list || 5 => Blacklist]

### Type Slot

TODO - Explain

[ 0 => Contract || 1 => Blob || 2 => Datamodel Only ]
