## Introduction

There are four major paradigms which we need to keep in mind when we are thinking about refactoring how the `eth` mining algorithm works. Each of these paradigms will exist somewhere on a spectrum between a completely centralized solution and a completely decentralized solution.

In order to properly explain each of the four places on the spectrum which will be operative for Eris Industries, I shall provide an example scenario. These scenarios should not be considered exclusive but rather as illustrative.

## Scenario 1: Paranoid Bank

Let us imagine we have a client who wants to deploy a smart contract system and use Eris' smart contracts but they do not want any other nodes to be able to interact with that system. Let us imagine that the client will deploy the system within an outer ring of security provided by their IT department. Let us imagine that the system will operate on servers within a client-owned data center. The only nodes able to make API calls to the server will be a predefined and relatively small number of nodes. This is Jack's Centralized Solution.

## Scenario 2: Semi-Paranoid Insurance Company

Let us imagine that we have a client who wants to deploy a smart contract system and use Eris' smart contracts but they want their workers to be able to transact with the system. This client is OK having Eris host the peer servers for the blockchain on Bostejo, but the client wants to restrict the nodes which can enter transactions into blocks to *only* those nodes hosted on the Bostejo system.

## Scenario 3: Coordinating Construction Company

Think one chain per large construction project. Let us imagine that we have a client who wants to deploy a smart contract system and use Eris' smart contracts but they need other companies to also be able to transact with the system and have other company's primary (peer server) nodes be able to enter transactions into blocks. The other company's nodes will be registered but may change over time as different company's (think subcontractors) are only operational for a particular portion of the overall project which the chain is used for.

## Scenario 4: Open-'er-up DApps

This would likely utilize the primary Etherum mining and transaction transmitting formula -- once that is established. But we should take into account that over time our systems will be used in a completely open way and the necessities and precautions which Ethereum proper is taking into account will be necessary to properly secure the integrity of this DApps' blockchain (even if it does not run on Ethereum's chain but on its own DApp chain).

## Tyler Durden's Security translation

## Scenario 1: Completely Internal and hidden

The World State is inaccessible to the end users. All Communication with the chain is done through API's from applications. Servers will hand back transactions to do perform the action and the user will sign it. Being completely blind to what it will do. All transactions not from a node will be ignore (no one can send in transactions from outside the system and have them processed even if valid and signed by a valid key).

* Miners - Trusted - possibly unreliable (downtime)
* Transactions - Trusted, must come from a valid node, must be signed by a valid key
* Nodes - Verified, restricted (World state hidden)

Mining Sync needed but all miners can be assumed to be uncompromised and safe. Most likely canidate Paxos/Raft for synchronization and peer verification. Due to the system probably want it so not EVERY node must verify. This will give improved performance. "Blocks" not completely necessary. Each transaction can be processed as it comes

## Scenario 2: Locked down mining and transaction signers

The World state is accessible to any node who wants to download it. However The only one capable of verifying changes to the world state are a select number of miners which are trusted. Only certain coinbases are allowed to submit transactions (all others will be ignored)

* Miners - Trusted, Possibly unreliable (downtime)
* Transactions - Arbitrary, Must be signed be valid coinbase to be accepted
* Nodes - Anyone

Mining Sync may be sufficient in this case Paxos/Raft should work for submission and peer verification of blocks. Possibly help the timing be more reliable and consistent.

Other notes: This case has been discussed thoroughly so heres some additional notes

1. Philosopher Kings, Miners can be constrained from submitting transactions of their own allowing the mining to be carried out by a party that you would not want to interact
2. transactors should have a couple possibilities.
	a. Can submit transactions with "Create" opcode
	b. Can submit transactions which do not have "Create" opcode
3. The Genesis DOUG should require a key different from any of the miners (incase they are compromised) Furthermore each use of the Genesis DOUG key should require a new key to be created and submitted as the new Genesis DOUG key - this is to quantum proof genesis DOUG. For the extremely paranoid we can set up a similar situation for transactors. making them simple gate keeper contracts which will forward the transaction when the correct key is used.

## Scenario 3: Groups of mutually untrusting parties

This case is actually non-trivial to solve. If all parties agree to hosting and mining by a neutral 3rd party it becomes simpler and reverts tot he previous case. In that instance nodes would be verifiers that nothing is going weird. THIS REQUIRES TRUST IN THE THIRD PARTY - but it is verifiable trust.

The reason this is non trivial is because all parties have interest in the world state and in certain circumstances it would be very valuable for something to be a different way. Remember Miners simply are Nodes which decide on the current world state. Normally following rules. But all rules can be amnipulated and if all nodes suddenly agree something else is truth then it is. The following scenario is illustrative:

Take the bitcoin blockchain, Miners submit blocks in order to get a reward. When the pool is large enough it is hard to find enough insentive to sway the direction of block mining. However say there are 3 people on a chain who all know each other Alice Bob and Carol who all mine. Now Alice sends a large payment to Bob in exchange for some goods. Alice does not want to really pay for it but she does not have enough mining power to sway it. However if Alice stands to save 100 BTC if that transaction gets reversed, Alice can go to Carol and offer her 50 BTC to help her mine the fork. Its sneaky dirty and rotten but suddenly Carol is more incentivized to not mine on the longest chain.

In our scenarios we have arbitrary databases and the value of any particular transaction can be negliable to invaluable. For example the signing of a contract. As such we need to be aware that the standard assumptions made in longest chain or other consensus algorithms may be invalid

## Scenario 4: Everybody want to be a node

Bunch of miners the larger the pool the better. Must go back to untrusted miners -> PoS or PoW