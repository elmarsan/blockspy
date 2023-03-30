# Blockspy

Little script that connects to a Bitcoin node using the Electrum protocol and subscribes to new block headers. 
Whenever a new block is received, the program decodes the block header, deserializes the block, and logs the height and hash of the new block.

It's also compatible with some Bitcoin forks.

## Usage
`go run main.go <client address>`