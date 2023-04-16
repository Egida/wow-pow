# Quote TCP Service with DDoS protection based on Proof of Work  

The possible solution for one of the interview task. 

## Problem 
Design and implement “Word of Wisdom” tcp server.

- TCP server should be protected from DDOS attacks with the Prof of Work, the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Prof Of Work verification, the server should send one of the quotes from “Word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge

## How to run
Built docker images:
```
make docker-images
```
Run server and client:
```
make docker-run-server
make docker-run-client
```

## Protocol
### Message structure 
#### Header
Size: 10 bytes
| Offset | Name         | Description                        |
|--------|--------------|------------------------------------|
| 0      | version      | version of used protocol           |
| 1      | type         | message type                       |
| 2      | payload size | size of upcoming payload in bytes |

#### Type of message 
| Type                         | Number | Description                                                    |
|------------------------------|--------|----------------------------------------------------------------|
| ProofOfWorkChallengeRequest  | 0      | Request to solve PoW puzzle (send by the server to the client) |
| ProofOfWorkChallengeResponce | 1      | Response with PoW solution (send by the client to the server   |
| Quote                        | 2      | Quote from the book                                            |

#### ProofOfWorkChallengeRequest
Size: depends on token size
| Offset | Name       | Description                          |
|--------|------------|--------------------------------------|
| 0      | difficulty | difficulty of the puzzle (up to 255) |
| 1      | token      | incoming token for the puzzle       |

#### ProofOfWorkChallengeResponce
Size: 8 bytes
| Offset | Name  | Description                 |
|--------|-------|-----------------------------|
| 0      | nonce | the solution for the puzzle |

#### Quote
Size: depends on quote size
| Offset | Name  | Description           |
|--------|-------|-----------------------|
| 0      | quote | a quote from the book |

### Interaction 
#### General 
- Server serialize a message into sequence of bytes (payload)
- Base on that sequence Server prepared a header a message and send it.
- Server send  the message. In version 1 it is 10 bytes. 
- Client read the header of the message (the first 10 bytes) that contains payload size
- Client received the payload and extract the original message

#### Flow of requisition a quote
- Client establishes connection with Server
- Server send a ProofOfWorkChallengeRequest message with a token and a difficulty
- Client solves the given puzzle 
- Client send a ProofOfWorkChallengeResponce message with a solution back
- Server checks the solution
- If the solution is correct Server send a quote from the book
- If the solution is wrong Server breaks off the connection


## Features
 - [Go client for external use](pkg/client/cleint.go:23)
 - [Go wrapper of protocol messages](pkg/protocol/massages.go)
 - [Concurrent implementation of Proof-of-Work fuction](internal/pow/pow.go:27)

```
cpu:   Intel(R) Core(TM) i7-10850H CPU @ 2.70GHz 8 Cores
cores: 8
BenchmarkFindNonceConcurrency/concurrency_1-8                 51        2298529773 ns/op            1374 B/op         13 allocs/op
BenchmarkFindNonceConcurrency/concurrency_2-8                100         183693295 ns/op            1171 B/op         15 allocs/op
BenchmarkFindNonceConcurrency/concurrency_4-8                126         317823195 ns/op            2067 B/op         23 allocs/op
BenchmarkFindNonceConcurrency/concurrency_8-8                172          91031197 ns/op            3175 B/op         39 allocs/op
BenchmarkFindNonceConcurrency/concurrency_16-8                55         517589201 ns/op            6035 B/op         72 allocs/op
```

## Proof-of-Work algorithm
PoW algorithm based Hashcash approach. Both Server and Client know shard Difficulty. Server generate a token. Client need to find nonce with which SHA256(token+nonce) has more leading bits then  value of Difficulty. [Implementation](internal/pow/pow.go)

This approach has following profs:
- simple implementation
- easiest way to parallelise calculation
- small footprint on the server size
- token size and puzzle difficulty can be dynamically adjusted for different type of client and network throughput

## Further development
- Continuous interaction
    The client supports a constant connection with the server and initiate new interaction by sending message with a special type.
- Adoptive puzzle generation
    Server can measure ddos response time for each client and adopt difficulty and token length for push client spend predefined time for finding puzzle solution.
- TechDept
    - Add retries to client
    - Add graceful shutdown to sever
    - Improve logging (change log to one of the well know logger)

