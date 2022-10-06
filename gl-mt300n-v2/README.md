# To build a new gl-mt300n-v2 firmware:

```sh
git clone https://github.com/gl-inet/imagebuilder.git
cd imagebuilder
cp /path/to/idig-station/gl-mt300n-v2/idig-station.json .
cp -r /path/to/idig-station/gl-mt300n-v2/files .
./gl_image -c idig-station.json -p idig-station
```

The new firmware should be created inside ```bin```.
