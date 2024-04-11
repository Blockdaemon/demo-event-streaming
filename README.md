
# Blockdaemon Chain Watch - streaming on-chain events in real-time 
  - In this example we'll use the Event Streaming service to monitor any on-chain transfers to the PEPE smart contract
  - Requirements:
    - Go lang
    - Blockdaemon REST API key
    - Reverse proxy service
    - [HTTPie](https://httpie.io/cli) & [jq](https://jqlang.github.io/jq/download/) CLI tools
  - Docs: https://docs.blockdaemon.com/reference/events-introduction

## Step 1. Run webhook receiver (includes a Challenge Response Check)
  - the sample webhook server will respond to CRC token checks with the secret "mysecret123"
  - all webhook posts will be printed to the console
  - the service will listen on port 8082
```shell
go run webhookserver.go
```

## Step 2. Publish the webhook server socket
  - [ngrok reverse proxy](https://ngrok.com/docs/getting-started/) is used as an example to publish the server
  - register for ngrok account or publish through your own methods
```shell
# set the Webhook Server FQDN
export WEBHOOKSERVERFQDN=XXXX.ngrok-free.app

# start the reverse proxy
ngrok http --domain=$WEBHOOKSERVERFQDN http://localhost:8082
```

## Step 3. 

```shell
# set Blockdaemon API key variable
export XAPIKey=XXXXX

# note Event Streaming protocol capabilities
http GET https://svc.blockdaemon.com/streaming/v2/ \
  X-Api-Key:$XAPIKey

# create webhook target - set the ngrok domain and note secret matches webhook server
TARGET_ID=$(http POST https://svc.blockdaemon.com/streaming/v2/targets \
    X-Api-Key:$XAPIKey \
    name='smartcontract webhook target' \
    description='Sample Go Webhookserver target' \
    max_buffer_count:=100 \
    settings:='{"destination": "https://'"$WEBHOOKSERVERFQDN"'", "method": "POST", "secret": "mysecret123"}' \
    type='webhook' \
    | jq -r .id)

# get the webhook target status - should be "connected"
http GET https://svc.blockdaemon.com/streaming/v2/targets \
    X-Api-Key:$XAPIKey 
```
## Step 4.
```shell
# create a smart contract address variable key
VARIABLE_ID=$(http POST https://svc.blockdaemon.com/streaming/v2/variables \
    X-Api-Key:$XAPIKey \
    name='smart contract addresses' \
    description='smart contract addresses' \
    type='string' \
    | jq -r .id)

# set the variable key value to the PEPE token contract address 0x6982508145454Ce325dDbE47a25d4ec3d2311933
http POST https://svc.blockdaemon.com/streaming/v2/variables/$VARIABLE_ID/values \
    X-Api-Key:$XAPIKey \
    value='0x6982508145454Ce325dDbE47a25d4ec3d2311933'
```
## Step 5.
```shell
# create rule combing the webhook target, the blockchain (Ethereum Mainnet), and the target address
RULE_ID=$(http POST https://svc.blockdaemon.com/streaming/v2/rules \
    X-Api-Key:$XAPIKey \
    name='Ethereum mainnet rule' \
    description='Ethereum mainnet rule' \
    network='mainnet' \
    protocol='ethereum' \
    condition:='[{"variable_type": "address", "variable_id": "'"$VARIABLE_ID"'"}]' \
    target="$TARGET_ID" \
    condition_type='match_var' \
    is_active:=true \
    template='ALL_DATA')

# validate the rule
http GET https://svc.blockdaemon.com/streaming/v2/rules \
    X-Api-Key:$XAPIKey

```

## Step 6.
Observe the webhook events in the webhookserver output. Example below:
```json
{
  "data": {
    "accessList": [],
    "additionalFields": {
      "yParity": "0x0"
    },
    "blockHash": "0xfc57c0ce8f94d133863f21dda0943b553e792e5fd1b9e687c680aa1b286bfd1f",
    "blockNumber": "0x12b845b",
    "chainId": "0x1",
    "contractAddress": null,
    "creates": null,
    "cumulativeGasUsed": "0x1bfc64d",
    "effectiveGasPrice": "0x2bad348ab",
    "from": "0xcAFb5420CE411476ef43CCeEa50e71A95b6Ad5B0",
    "gas": "0x163ed",
    "gasPrice": "0x2bad348ab",
    "gasUsed": "0xd986",
    "hash": "0x06424f8c03c746ea1ccf3167d454a4ca54351bb4b179c9b3803455872cc31248",
    "input": "0xa9059cbb00000000000000000000000031da1b45d2570b5722618b8b43b0c805114048ee00000000000000000000000000000000000000000031c7ea725f3bca431e6000",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000008000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000010000004000000000000000000000000000000000000000000000000000000000200000000000040000000000000000000000000000000002000000000000000000000000000000002000000000400000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000008000",
    "maxFeePerGas": "0x40c933ae2",
    "maxPriorityFeePerGas": "0xecd10",
    "nonce": "0x3d6",
    "publicKey": null,
    "r": "0x2b12ee7984f96db6cf16a190e257caa0392fa8961243488ff07caa4931adaffe",
    "raw": null,
    "root": null,
    "s": "0x6ed8f306cdfa24c4aea4f30ec5cfb853dfe90adc6d32db37d60195e36e3b1422",
    "status": "0x1",
    "timestamp": 1712799383,
    "to": "0x6982508145454Ce325dDbE47a25d4ec3d2311933",
    "transactionIndex": "0xd0",
    "type": "0x2",
    "uuid": "8e695b82-7ed3-4ab0-894d-bce7a91077a6",
    "v": "0x0",
    "value": "0x0"
  },
  "event_type": "confirmed_tx",
  "network": "mainnet",
  "protocol": "ethereum",
  "rule_id": "34b775e7-4425-4487-9181-0436022588fc",
  "target_id": "9c6b6a09-dd23-40b7-bb57-7aa0673f2a20"
}
```
where:
```json
data.input: 0xa9059cbb00000000000000000000000031da1b45d2570b5722618b8b43b0c805114048ee00000000000000000000000000000000000000000031c7ea725f3bca431e6000
```
ABI decodes to:
```json
      "definition": "transfer(address,uint256)",
      "decodedInputs": [
        "0x31DA1B45d2570B5722618B8b43b0C805114048EE",
        "60181440870692720000000000"
      ]
```
meaning `60181440870692720000000000` PEPE was transferred from `0xcAFb5420CE411476ef43CCeEa50e71A95b6Ad5B0` to `0x31DA1B45d2570B5722618B8b43b0C805114048EE`



## Step 7.
```shell
# delete the rule to halt further notifitions
http DELETE https://svc.blockdaemon.com/streaming/v2/rules/$RULE_ID \
    X-Api-Key:$XAPIKey
```