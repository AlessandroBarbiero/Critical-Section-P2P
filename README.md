# Critical-Section-P2P
Implement distributed mutual exclusion between nodes in a distributed system. 

## Run network

```bash
./run.sh
```

The script reads port from config file `clients.info` and runs clients on those ports.
Each stdin is spawned as new terminal emulator. Emulator name has to be saved in variable `$TERM`

## Config file

Ports for all clients have to be specified in `clients.config` file. 
Each line of file represents single node and only thing on that line should be port of that node.
Next node to send the token to is the node on the next line in config file.
If the node is the last one in config file, next node to send the token to is specified on first line of the file.
