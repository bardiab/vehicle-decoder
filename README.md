# NHTSA Vehicle Decoder

A command line utility (CLI) written in Go, that can take in a vehicle's VIN number and provide the output of the API's response cleanly and efficiently.

### Notes:
I used the Cobra library for Golang for initial project setup


### CLI Usage:
To run the CLI, input a VIN as so, and choose one of the following arguments.

```
./vehicle-decoder decode --vin=<VIN>
```

Arguments supported: 
> --raw

Output the raw JSON response of the API.

> --yaml 

Output the raw YAML response of the API.

> --sparse 

Output how many fields and results were returned to the response.

> --meta

Output how many fields and results were returned to the response, including null ones.

> --fields

Query the fields that match the comma separated regular expression provided to the argument. Only returns fields with values that are not empty or null. 

example: 
```
--fields="Manufacturer Name, Make, Model, Plant.*"
```