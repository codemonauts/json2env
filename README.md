# json2env

In our CI/CD pipelines we regularly had the same pattern emerge which resulted in a lot of boilerplate code.
We always read a secret from the AWS SecretsManager, torn it apart with jq and loaded it into environment variables to use it e.g. with sed to put them into different configuration files.

This tool can either read json from stdin, or if `-secret-id` is given, directly read key/value pairs from AWS
SecretsManager.

It will output a list of `export KEY=value` lines to stdout so you could  just wrap
it into an `eval $()` call in your pipeline script.

## Usage
To directly query AWS SecretsManager use the *-secret-id* flag:
```bash
json2env -secret-id myAWSsecret
```

if no *-secret-id* argument is provided, the tool will expect json input via STDIN:

```bash
echo '{"username":"johndoe","password":"swordfish"}' | ./json2env             
```

The tool will print something like this to the terminal:
```bash
export username=johndoe
export password=swordfish
```

To get this directly into environment variables you can wrap the call into `eval` and execute it into your current
terminal:

```bash
$ eval $(echo '{"username":"johndoe","password":"swordfish"}' | ./json2env -prefix 'SM_')
$ env | grep "SM_"
SM_username=johndoe
SM_password=swordfish
```
