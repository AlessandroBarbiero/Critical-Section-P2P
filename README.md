# Critical-Section-P2P
Implement distributed mutual exclusion between nodes in a distributed system. 

## Run network

```bash
./run.sh ${port of first node with access to critical area}
```

The script reads port from config file `clients.info` and runs clients on those ports.
Each stdin is spawned as new terminal emulator. Emulator name has to be saved in variable `$TERM`  
To set required terminal emulator run `TERM=${your terminal emulator}` before running `./run.sh`

## Config file

Ports for all clients have to be specified in `clients.config` file. 
Each line of file represents single node and only thing on that line should be port of that node.
Next node to send the token to is the node on the next line in config file.
If the node is the last one in config file, next node to send the token to is specified on first line of the file.

## Use the network

After running `run.sh` you see `n` (n being the number of nodes in your clients.info file) terminals.
In each terminal you see the message with port of the node, which terminal belongs to.
If you want to access the critical area just input random new line into terminal.
Node will access the critical area, when it's its turn.

## Script tested on

`5.18.0-1parrot1-amd64` (basically debian)