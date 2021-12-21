# Glossary

### Masterplan

A _masterplan_ is comprised of filters that target a set of _machines_, _variables_
that can later be used in _templates_ and a set of _plan_ that will be executed
against the _machines_ at run time.

### Plan

A _plan_ is comprised of a set of defaults variables and a set of sub-plans.
Those subplans can always be one of two things: a command line action to run
or a call to another plan.

### Machine

A _machine_ is the representation of a computer or server the user has access to
and often is associated with ssh login configuration and a hostname.

### Inventory

A set of _machines_. Often represented by the `inventory.json` file.

### Run

The act of running a masterplan against an inventory.

### File

A _file_ of any format that can used by a _plan_, if ending by `.tmpl` this file
will be considered as a _template_ by `hotomata`.

### Template

A _template_ is either a string or a _file_ that will be passed by a template
engine before being passed to a command. This allows that user to place variable
either coming from the _inventory_, the _masterplan_ or the _plan_ itself.
