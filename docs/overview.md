# Overview

This file gives an overview of the different types in `hotomata` and how they
interract together.

## Machine

A machine that can be contacted by one of the runners (defaults to SSH) so that,
in term, commands are ran on it.

## Inventory

A list of `groups` of `machines` in **JSON** format. This file describes the
available `machines` than would later be targeted by `masterfile` directives.

For more more information: [Inventory file](https://github.com/merd/hotomata/blob/master/docs/inventory_file.md)

## Masterplan

For more more information: [Masterplan file](https://github.com/merd/hotomata/blob/master/docs/masterplan_file.md)

## Plan

## Command

In the end, the base primitive is a shell **command**, a string. All plans are
is a collection of commands often templated with vars given to it. You can see
any plan as tree with multiple branches which all end up in leading to leafs, in
our case command to run over SSH, locally or elsewhere

## File

A file is located in the `files` folder in the same directory as the `masterplan`
and can be used with different commands, they can be templated and often end up
uploaded on remote machines.

## Varfile / Vault

A varfile is a `file` containing variables in a **yaml** format and it can be
included in `masterfiles` directives.

Often, the goal of having those in a separate file is the ability to encrypt
this file independently. `hotomata` will have helper to make this seamless.
It is very much inspired by _Ansible Vault_ utility.
