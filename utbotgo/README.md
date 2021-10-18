# UTBotGO

Program for generating unit tests for Go files.

## How to install

*   ```bash
    make install
    ```
    
    It installs command `utbotgo` on your system, so you will be able to run the
    command `utbotgo` from any directory.

    Note: You may need to run this command together with `sudo`.

*   ```bash
    make build
    ```
    
    It creates executable file `utbotgo` in this directory. You will be able to
    run command `utbotgo` only from this directory.

## How to run

```bash
make run <file_or_dir>
```
or
```bash
make build
./utbotgo <file_or_dir>
```
or
```bash
make install
utbotgo <file_or_dir>
```

Note: Last variant allow you to run `utbotgo` command from any directory.

## How to use

Instead of `<file_or_dir>` specify Go file or directory with Go files, for
which you want to generate testing files (files `*_test.go`).

Output will contain amount of generated files.

For example, if you have Go file `hello.go` and run
```bash
make run hello.go
```
testing file `hello_test.go` will be generated, and output will be `1`.

## How to uninstall

```bash
make uninstall
```

It removes command `utbotgo` from your system.

Note: You may need to run this command together with `sudo`.
