# `sync-at-time` modular resource

This module allows you to configure Cloud Sync to occur only at a specific time frame by implementing a sensor, `naomi:sync-at-time:timesyncsensor`, that enables sync when within a specified timeframe and disables sync when outside that timeframe.

## Requirements

Before configuring your sensor, you must [create a machine](https://docs.viam.com/manage/fleet/machines/#add-a-new-machine).

To use the `sync-at-time` module you also need to:

1. Enable [data capture](https://docs.viam.com/data/capture/).
2. Enable [cloud sync](https://docs.viam.com/data/cloud-sync/).

## Build and run

To use this module, follow these instructions to [add a module from the Viam Registry](https://docs.viam.com/registry/configure/#add-a-modular-resource-from-the-viam-registry) and select the `naomi:sync-at-time:timesyncsensor` model from the [`sync-at-time` module](https://app.viam.com/module/<INSERT API NAMESPACE>/<INSERT MODEL>).

## Configure your `sync-at-time` sensor

Navigate to the **Config** tab of your machine's page in [the Viam app](https://app.viam.com/).
Click on the **Components** subtab and click **Create component**.
Select the `<INSERT API NAME>` type, then select the `<INSERT MODEL>` model.
Enter a name for your <INSERT API NAME> and click **Create**.

On the new component panel, copy and paste the following attribute template into your sensor’s **Attributes** box:

```json
{
  "start": "HH:MM:SS",
  "end": "HH:MM:SS",
  "zone": "<TIMEZONE>"
}
```

> [!NOTE]
> For more information, see [Configure a Machine](https://docs.viam.com/manage/configuration/).

### Attributes

The following attributes are available for the `naomi:sync-at-time:timesyncsensor` sensor:

| Name    | Type   | Inclusion    | Description |
| ------- | ------ | ------------ | ----------- |
| `start` | string | **Required** | The start time for the time frame during which you want to sync, example: `"14:10:00"`.  |
| `end`   | string | **Required** | The start time for the time frame during which you want to sync, example: `"15:35:00"`. |
| `zone`  | string | **Required** | The time zone for the `start` and `end` time, for example: `"CET"`. |

### Example configuration

```json
{
  "start": "14:10:00",
  "end": "15:35:00",
  "zone": "CET"
}
```

### Configure data manager

On your machine's **Config** tab, switch to **JSON** mode and add a `selective_syncer_name` with the name for the sensor you configured:

```json {class="line-numbers linkable-line-numbers" data-line="6"}
{
  "attributes": {
    "additional_sync_paths": [],
    "capture_dir": "",
    "capture_disabled": false,
    "selective_syncer_name": "selective-syncer",
    "sync_disabled": false,
    "sync_interval_mins": 0.1,
    "tags": []
  },
  "name": "Data-Management-Service",
  "type": "data_manager"
}
```

### Next steps

_Add any additional information you want readers to know and direct them towards what to do next with this module._
_For example:_

- To test your...
- To write code against your...

## Local Development

This module is written in Go.

```bash
go mod tidy
go build
```

## License

Copyright 2021-2023 Viam Inc. <br>
Apache 2.0