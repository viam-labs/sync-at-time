# `sync-at-time` modular resource

This module allows you to configure Cloud Sync to occur only at a specific time frame by implementing a sensor, `naomi:sync-at-time:timesyncsensor`, that enables sync when within a specified timeframe and disables sync when outside that timeframe.

## Requirements

Before configuring your sensor, you must [create a machine](https://docs.viam.com/fleet/machines/#add-a-new-machine).

To use the `sync-at-time` module you also need to:

1. Enable [data capture](https://docs.viam.com/data/capture/).
2. Enable [cloud sync](https://docs.viam.com/data/cloud-sync/).

## Build and run

To use this module, follow these instructions to [add a module from the Viam Registry](https://docs.viam.com/registry/configure/#add-a-modular-resource-from-the-viam-registry) and select the `naomi:sync-at-time:timesyncsensor` model from the [`sync-at-time` module](https://app.viam.com/module/naomi/sync-at-time).

## Configure your `sync-at-time` sensor

Navigate to the **Config** tab of your machine's page in [the Viam app](https://app.viam.com/).
Click on the **Components** subtab and click **Create component**.
Select the `sensor` type, then select the `sync-at-time:timesyncsensor` model.
Click **Add module**, then enter a name for your sensor and click **Create**.

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

```json
{
  "attributes": {
    "additional_sync_paths": [],
    "capture_dir": "",
    "capture_disabled": false,
    "selective_syncer_name": "<SENSOR-NAME>",
    "sync_disabled": false,
    "sync_interval_mins": 0.1,
    "tags": []
  },
  "name": "Data-Management-Service",
  "type": "data_manager"
}
```

For an example configuration see [example.json](./example.json).

### Next steps

You have now configured sync to happen during a specific time slot.
To test your setup, [configure a webcam](https://docs.viam.com/components/camera/webcam/) or another component and [enable data capture on one of the component methods](https://docs.viam.com/data/capture/#configure-data-capture-for-individual-components), for example `ReadImage`.
The data manager will now capture data.
Go to the [**Control** tab](https://docs.viam.com/fleet/machines/#control). You should see the sensor.
Click on `GetReadings`.
If you are in the time frame for sync, the sensor will return true.
You can confirm that no data is currently syncing by going to the [**Data** tab](https://app.viam.com/data/view).
If you are not in the time frame for sync, adjust the configuration of your sensor.
Then check again on the **Control** and **Data** tab to confirm data is syncing.

## Local Development

This module is written in Go.

```bash
go mod tidy
go build
```

## License

Copyright 2021-2023 Viam Inc. <br>
Apache 2.0