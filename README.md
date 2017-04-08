# Spartan

Spartan is a data process much like Logstash. It uses a programatic pipeline structure to move events generated
from inputs through a set of filters and eventually to outputs. Spartan was built mainly for academic/experimentation
purposes. I wanted to see what I took to build something like Logstash.

Spartan is written in Go and should build on all platforms targeted by he Go runtime. Development is mainly on
Linux but if any problems arise, please put in an issue and I'll take a look at it. The project is targeted to
Go 1.7.

## Configuration

The filter configuration syntax is very similar to Logstash. It shares the same basic structure with only
minor modifications necessary.

## Inputs

Currently supported inputs:

- File

## Filters

Currently supported filters:

- Grok
- Date
- Mutate

## Outputs

Currently supported outputs:

- Stdout
