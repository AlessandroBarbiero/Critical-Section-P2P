# Critical-Section-P2P
Implement distributed mutual exclusion between nodes in a distributed system. 

## Run network
In order to easily run all the nodes at ones we wrote a shell and a batch file.
The scripts read ports from config file `clients.info` and runs clients on those ports.

### Unix

```bash
./run.sh ${port of first node with access to critical area}
```
Each stdin is spawned as new terminal emulator. Emulator name has to be saved in variable `$TERM`  
To set required terminal emulator run `TERM=${your terminal emulator}` before running `./run.sh`

### Windows

```bash
run.bat [port of first node with access to critical area]
```
For example `run.bat 5000` \
Each stdin is spawned as new cmd prompt.


## Config file

Ports for all clients have to be specified in `clients.info` file. 
Each line of the file represents a single node and the only thing on that line should be the port of that node.
Next node to send the token to is the node on the next line in config file.
If the node is the last one in the config file, next node to send the token to is specified on first line of the file.

## Use the network

After running the script you will see `n` (n being the number of nodes in your clients.info file) terminals.
In each terminal you will see the message showing witch port the node is running on.
If you want to access the critical area just input random new line into terminal.
Node will access the critical area, when it's its turn.

## Script tested on

- `5.18.0-1parrot1-amd64` (basically debian)
- `Windows 10`