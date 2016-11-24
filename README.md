# About MNDP

The **MNDP** is a proprietary **M**ikrotik **N**eighbor **D**iscovery **P**rotocol. The protocol runs on [Mikrotik](https://www.mikrotik.com) made [network 
equipment](https://routerboard.com/) similar to *CDP* which runs on Cisco equipment. 

The repository contains source code for 3rd party implementation attempt of **MNDP** listener written in golang:
* `discover` command (commandline tool)
* `mndplib` package (golang library)

## Installation

1. Install [golang](https://golang.org/doc/install) and [git](https://www.git-scm.com/)
2. Setup environment:
    - Set `GOPATH` to go workspace directory
    - Set `PATH` to contain `GOROOT/bin` and `GOPATH/bin`
3. Install the discover command by 
    ```
    go get -u github.com/pjediny/mndp/discover/
    ```

## Using `discover`

Run the `discover` command from terminal and you should be able to see the incoming MNDP discovery messages.
- Press `<ENTER>` to request refresh
- Press `<q> <ENTER>` to quit

## License

[MIT License](LICENSE.txt)
