# To build a new gl-mt300n-v2 firmware:

1. ```git clone https://github.com/gl-inet/imagebuilder.git```
2. ```cd imagebuilder```
3. ```cp /path/to/idig-station/gl-mt300n-v2/idig-station.json .```
4. ```cp -r /path/to/idig-station/gl-mt300n-v2/files .```
5. ```./gl_image -c idig-station.json -p idig-station```

The new firmware should be created inside ```bin```.
