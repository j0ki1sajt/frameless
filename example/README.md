# Example layered application

This example project want's to prove how to create software where the parts are independent even while they work with each other.
Most task this way can be developed parallel by different developers as well.

This is not a real world application, only for to create presentation how this idiom works.

## Development

The development workflow was the following:
* Create the use cases
* create/use-in a Storage
* create external interfaces
    * HTTP
    * CLI

for more check git history

## CLI

have two sub command:
* add
  * require -t and -c option
* list

## Web

have to fancy page:
* /add?Title=foo&Content=baz
* /list
