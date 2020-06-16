<h1 align="center">
  <br>
  <a href="https://github.com/gopinath-langote/1build">
    <img src="assets/logo.png" alt="2lazy" height="200"></a>
  <br>
</h1>

2lazy is a tool allowing you to create and use project specific command aliases.

[1build](https://github.com/gopinath-langote/1build/) is a similar tool with a different execution model. `1build` allows you to run multiple commands in a single run, but offers no way of interacting with them. `2lazy` allows you to specify additional arguments for your aliases and fully interact with them (stdout, stderr, stdin and exit code are redirected) but you can only run one command per invocation.

## Basic usage

First of all add a `2lazy.yml` file to your project. Here you can find minimal example.

```yaml
commands:
  hello: echo Hello
```

```console
$ 2lazy hello
Hello
INFO Finished                                      elapsedTime="637.598µs"
$ 2lazy hello world
Hello world
INFO Finished                                      elapsedTime="566.128µs"
```

## Advanced usage

```yaml
quiet: true

project_dir: /tmp
start_in_project_dir: true

commands:
  interactive: docker run --rm -it debian /bin/bash
  cut_spaces: tr -d ' '
  ls: ls -l
  pwd: pwd
```

### Running interactive software

```console
$ 2lazy interactive
root@fe74e2d6b84b:/# echo "That's great!"
That's great!
```

### Input/output redirection

Remember to set `quiet: true` if you want to use standard output!

```console
$ echo "2 + 2" | 2lazy cut_spaces | tee /dev/stderr | bc
2+2
4
```

(`tee /dev/stderr` makes the output go to both standard error so you can see the result and to standard output which is redirected to bc for an actual calculation.)

### Exit code redirection

```console
$ 2lazy ls unknown.png
ls: cannot access 'unknown.png': No such file or directory
$ echo $?
2
```

### Overriding working directory

You can start `2lazy` in any of the child directories of the project - it'll automatically find configuration in one of the parent files. By default commands will be spawned i the same directory as `2lazy` was spawned, but you can override this behavior with `start_in_project_dir` and `project_dir` options.

```console
$ 2lazy pwd
/tmp
```

## Available configuration keys

| Key name               | Default value                                   | Description                                                                                                      |
| ---------------------- | ----------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `quiet`                | `false`                                         | If `true` will suppress messages with level lower or equal to `INFO` (like the `--quiet` option in CLI)          |
| `start_in_project_dir` | `false`                                         | If `true` will start commands in `project_dir`, else in the same directory the `2lazy` was spawned               |
| `project_dir`          | Directory where first `2lazy.yml` file is found | Allows you to specify where to start commands                                                                    |
| `commands`             | Empty dictionary                                | Dictionary with actual list of aliases. Currently all parameters are passed 1:1, but it may change in the future |

## Help section

```console
NAME:
   2lazy - when you're too lazy to type full command

USAGE:
   2lazy [global options] [alias name] [alias arguments] [...]

GLOBAL OPTIONS:
   --debug     whether to show debug messages (default: false)
   --quiet     whether to hide info messages (default: false)
   --help, -h  show help (default: false)
```

## Possible future enhancements

- bash-like argument substitution (think `$1` and similar)
- parsing all the `2lazy.yml` files in the parent tree
- support for other configuration formats (like [Jsonnet](https://jsonnet.org/), JSON, etc.)
- support for saving environment variables in configuration file
- suppoprt for special arguments like `~`

## Contributions

Feel free to submit pull requests. Files are formated with gofmt with default settings.

## License

MIT 2020 by [pidpawel](https://pidpawel.eu/)
