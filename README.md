`env`

# Overview

`env` is a tiny library designed to simplify interactions with the environment in which your application runs. It provides an easy-to-use API for fetching environment variables, setting them, and performing various other environment-related tasks.

It exists mainly to support functionality of `nv`.

# Installation

```shell
~> go get github.com/jcouture/env
```

# Features

- Clear the environment with exceptions
- Check if an environment variable exists
- Retrieve all environment variables
- Set environment variables dynamically
  And more!

# Usage

_**Important note**: this library works on a copy of your environment, just like Go’s stdlib._

## Import the library

```go
import "github.com/jcouture/env"
```

## Clear the environment

Clears the environment but will skip any (optional) variable names passed as parameter.

```go
env.Clear("HOME")

```

## Check if a variable exists

```go
if !env.Exists("FOO") {
  fmt.Printf("FOO not found\n")
}
```

## Retrieve all variables

```go
vars := env.Getvars()

for k, v := range vars {
  fmt.Printf("%v=%v\n", k, v)
}
```

## Retrive all variable names

```go
names := env.Getnames(env.Getvars())

for _, n := range names {
  fmt.Printf("%v\n", n)
}
```

## Join two maps together to override variables

```go
base := env.Getvars()

override := map[string]string{
  "PATH": "/usr/bin:/bin:/usr/sbin:/sbin",
  "HOME": "/home/user",
}

base = env.Join(base, override)
```

## Set variables

```go
vars := map[string]string{
  "SHELL": "/usr/bin/zsh",
  "TERM":  "xterm-256color",
}

env.Setvars(vars)

```

# License

`env` is released under the MIT license. See [LICENSE](./LICENSE) for details.
