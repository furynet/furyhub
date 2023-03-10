---
order: 1
---

# Introduction

## GRID Network

The GRID network is an internet of blockchains intended to provide a technology foundation that facilitates construction of distributed business applications.

![Figure of GRID Network](../pics/chap2-1.jpg)

The GRID network is part of the larger Cosmos network -- all zones in the network would be able to interact with any other zone in the Cosmos network over the standard IBC protocol.  By introducing a layer of service semantics into the network, we are going to provide an innovative solution that enables a whole new set of business scenarios, which would result in an increase in scale and diversity of the Cosmos network.

## GRID Hub

At the "center" of the GRID network is a blockchain known as the *GRID Hub*, which is a Proof-of-Stake (PoS) blockchain built on Cosmos SDK and Tendermint.  It will be the first regional hub that connects to the Cosmos Hub as one of its zones.  The GRID hub is equipped with a service protocol that coordinates on-chain transaction processing with off-chain data processing and business logic execution.  We have also enhanced the IBC protocol to facilitate cross-chain invocation of those off-chain services, if so desired.

The service protocol and enhanced IBC protocol could eventually be contributed back into Cosmos SDK, allowing SDK users to develop blockchains that are compatible with the GRID network.

Client-facing, programming language specific SDKs will also be available to make it easy to provide and consume off-chain services within the GRID network.

## GRID Tokens

The GRID hub has its own native token known as *GRID*.  It is designed to serve three purposes in the network.

* **Staking.**  Similar to the ATOM token in the Cosmos Hub, the GRID token will be used as a staking token to secure the PoS blockchain.

* **Transaction Fee.**  The GRID token will also be used to pay fees for all transactions in the GRID network.

* **Service Fee.**  It is required that service providers in the GRID network charge service fees denominated in the GRID token.

It is intended that the GRID network will eventually support all whitelisted fee tokens from the Cosmos network, which can be used to pay the transaction fees and service fees.

## GRID Services

*GRID Services* (a.k.a. "iService") are introduced to bridge the gap between the blockchain world and the conventional business application world, mediating the complete lifecycle of off-chain services -- from their definition, binding (provider registration), invocation, to their governance (profiling and arbitration).

### Lifecycle

* **Definition:** Definition of what an off-chain iService can do in terms of an Interface Definition Language (IDL) file.

* **Binding:** Declaration of the location (address), pricing and QoS of a provider endpoint that implements a given iService definition.

* **Invocation:** Handling of consumer requests to and provider responses from a given iService provider endpoint.

### Providers

*Providers* are network users who offer the implementation of one or more iService definitions and often act as *adaptors* of off-chain services and resources located in other public and consortium chains, as well as in enterprise legacy systems.  Providers monitor and process incoming requests and send responses back to the network.  A provider could at the same time act as a consumer by sending requests to other providers.  As planned, providers would be required to charge a fee for any services they might offer, and the service fee, by default, would be priced in the GRID token.

### Consumers

*Consumers* are those users who consume iService by sending requests to designated provider endpoints and receiving responses from providers in question.

### Profilers

*Profilers* are special consumers who act on behalf of the GRID Foundation Limited, a Hong Kong incorporated company limited by guarantee that takes the lead in building the GRID network.  Profilers are the sole users authorized to invoke iService in the *profiling mode*, which is intended to help create and maintain objective provider profiles that consumers refer to for provider screening.

### Arbitrators

*Arbitrators* are self-declared users who, working collectively, arbitrate consumer complaints against provider performance.  The details about the arbitration mechanism are being actively worked on, please keep an eye on our [whitepaper](../resources/whitepaper.md).
