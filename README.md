# How to Run the code


# Architecture
The architecture of the system is like:

![architecture](https://github.com/froyobin/downloadledger/raw/master/figures/architecture.png)


## 1. install the smart contract
* following the instrusctions from [https://hyperledger-fabric.readthedocs.io/en/release-1.4/tutorials.html](https://hyperledger-fabric.readthedocs.io/en/release-1.4/tutorials.html) to set up the hyperledger environment.

* put contract/chaincode_example02.go under fabric-samples/chaincode/chaincode_example02/go
* put contract/docker-compose-simple.yml under fabric-samples/chaincode-docker-devmode
* modify the volume mapping in the yml file.

** *************************** **
*  change directory to fabric-samples/chaincode-docker-devmode and run the command **docker-compose -f docker-compose-simple.yaml up**

*  
   **under  docker exec -it chaincode sh**
    - cd chaincode_example02/go
    - go build -o chaincode_example02
    - CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./chaincode_example02
    
    
   **under  docker exec -it cli bash**
    -   peer chaincode install -p chaincodedev/chaincode/chaincode_example02/go -n mycc -v 0
   -   peer chaincode instantiate -n mycc -v 0 -c '{"Args":["init","a","100","b","200"]}' -C myc

## 2. In another terminal run the data server, data owner and client seperately.
  
* server:
  - cd /app/server and run ./server 

* dataowner:
  - cd /app/dataowner and run ./server 

* client:
  - cd /app/client and run ./main

you should see the recovery of the master key you encrypt the fille.
```
mast key mykey44444444444
```


**once you change the source file, dont forget to re-compile the source file.**