---
canonical: https://grafana.com/docs/alloy/latest/get-started/modules/
aliases:
  - ../concepts/modules/ # /docs/alloy/latest/concepts/modules/
description: Learn about modules
title: Modules
weight: 400
---

# Modules

A _Module_ is a unit of {{< param "PRODUCT_NAME" >}} configuration that combines all other concepts.
It contains a mix of configuration blocks, instantiated components, and custom component definitions.
The module you pass as an argument to [the `run` command][run] is called the _main configuration_.

You can [import modules](#import-modules) to reuse [custom components][] defined by that module.

## Import modules

You can _import_ a module to use its custom components in other modules, called _importing modules_.
Import modules from multiple locations using one of the `import` configuration blocks:

* [`import.file`][import.file]: Imports a module from a file on disk.
* [`import.git`][import.git]: Imports a module from a file in a Git repository.
* [`import.http`][import.http]: Imports a module from an HTTP request response.
* [`import.string`][import.string]: Imports a module from a string.

{{< admonition type="warning" >}}
You can't import a module that contains top-level blocks other than `declare` or `import`.
{{< /admonition >}}

Modules are imported into a _namespace_, exposing the top-level custom components of the imported module to the importing module.
The label of the import block specifies the namespace of an import.
For example, if a configuration contains a block called `import.file "my_module"`, then custom components defined by that module are exposed as `my_module.CUSTOM_COMPONENT_NAME`.
Namespaces for imports must be unique within a given importing module.

If an import namespace matches the name of a built-in component namespace, such as `prometheus`, the built-in namespace is hidden from the importing module.
Only components defined in the imported module are available.

{{< admonition type="warning" >}}
If you use a label for an `import` or `declare` block that matches an existing component, the component is shadowed and becomes unavailable in your configuration.
For example, if you use the label `import.file "mimir"`, you can't use existing components starting with `mimir`, such as `mimir.rules.kubernetes`, because the label refers to the imported module.
{{< /admonition >}}

## Example

This example module defines a component to filter out debug-level and info-level log lines:

```alloy
declare "log_filter" {
  // argument.write_to is a required argument that specifies where filtered
  // log lines are sent.
  //
  // The value of the argument is retrieved in this file with
  // argument.write_to.value.
  argument "write_to" {
    optional = false
  }

  // loki.process.filter is our component which executes the filtering,
  // passing filtered logs to argument.write_to.value.
  loki.process "filter" {
    // Drop all debug- and info-level logs.
    stage.match {
      selector = `{job!=""} |~ "level=(debug|info)"`
      action   = "drop"
    }

    // Send processed logs to our argument.
    forward_to = argument.write_to.value
  }

  // export.filter_input exports a value to the module consumer.
  export "filter_input" {
    // Expose the receiver of loki.process so the module importer can send
    // logs to our loki.process component.
    value = loki.process.filter.receiver
  }
}
```

You can save this module to a file called `helpers.alloy` and import it:

```alloy
// Import our helpers.alloy module, exposing its custom components as
// helpers.COMPONENT_NAME.
import.file "helpers" {
  filename = "helpers.alloy"
}

loki.source.file "self" {
  targets = LOG_TARGETS

  // Forward collected logs to the input of our filter.
  forward_to = [helpers.log_filter.default.filter_input]
}

helpers.log_filter "default" {
  // Configure the filter to forward filtered logs to loki.write below.
  write_to = [loki.write.default.receiver]
}

loki.write "default" {
  endpoint {
    url = LOKI_URL
  }
}
```

## Security

Since modules can load arbitrary configurations from potentially remote sources, carefully consider the security of your solution.
The best practice is to ensure attackers can't modify the {{< param "PRODUCT_NAME" >}} configuration.
This includes the main {{< param "PRODUCT_NAME" >}} configuration files and modules fetched from remote locations, such as Git repositories or HTTP servers.

[custom components]: ../custom_components/
[run]: ../../reference/cli/run/
[import.file]: ../../reference/config-blocks/import.file/
[import.git]: ../../reference/config-blocks/import.git/
[import.http]: ../../reference/config-blocks/import.http/
[import.string]: ../../reference/config-blocks/import.string/
