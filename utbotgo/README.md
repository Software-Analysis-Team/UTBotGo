# UTBotGo

## Building

```
make
```

This command build docker image from `Dockerfile` and build helpful utilities
in `utils` directory.

## Installing

```
make install
```

This install `utbotgo` command to your system.

> Note: it may require `root` rights.

## Running

1.  Go to example

```
cd ../examples/test1/
```

2.  Initialize directory for working `utbotgo` command

```
utbotgo init
```

3.  Edit `config.yml` file in directory `utbotgo/`

```
echo "tested-functions:" > utbotgo/config.yml
echo "  - getSign" >> utbotgo/config.yml
```

4.  Run testing

```
utbotgo test
```

5.  (optional) See generated test arguments for function `getSign`

```
cat utbotgo/functions/getSign/args.json
```
