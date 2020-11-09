# Example snippet for grafana streamed json datasource from http backend
This example represents an http backend server that streams json data to a locally hosted grafana server.
Very crude example which is still heavily under development.
Examples used:
[Build streaming datasource for grafana](https://grafana.com/docs/grafana/latest/developers/plugins/build-a-streaming-data-source-plugin/)
[Simple streaming datasrouce](https://github.com/seanlaff/simple-streaming-datasource)

# How it works?
1. The example consist of an http backend and a grafana docker container.
2. When the example DataStream Dashboard is opened the datasource will request the streaming from the backend
3. Once the response starts to arrive to the data source the data will be displayed in the charts

# How to run?
## Setup
- Run ```docker-compose up --build --force-recreate -d main-server```
- Access grafana on `http://localhost:3000`
- Set data-stream datasource with the host address `http://localhost:8080`
- Import the demo dashboard from under the Dashboards tab

## Running
- Once the setup steps are done the line chart panel can be accessed from the backend through `http://localhost:8080/show`
- The data name can be changed in the panel edit page under the dataText field. If multiple series need to be shown, names need to be separated by a comma (For example: testData1,testData2).

## Build the datasource
- navigate in `./grafana/json-data-stream`
- run `yarn install` (may requires to install some linux node.js packages)
- run `yarn dev`

# Known issues
This example is meant to be a POC and it is in a very early stage yet.

[Only handles panels in a single dashboard] (https://github.com/artofimagination/grafana-json-streaming-datasource/issues/1)

[Data display gets slow after a while] (https://github.com/artofimagination/grafana-json-streaming-datasource/issues/2)

The panel visualization is cleared every time you update the dashboard. If you have access to historical data, you can add, or backfill, it to the data frame before the first call to subscriber.next().

[There are building issues with the datasource (the fail does not prevent the usage of the datasource)](https://github.com/artofimagination/grafana-json-streaming-datasource/issues/3)
