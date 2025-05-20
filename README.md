# go_pidea-get-users-by-realm
Go: get PrivacyIdea users via REST API by given realm

For use in Windows mostly(linux not tested)

<h2>Data file(data.json)</h2>

```
{
    "pideaUrl": "https://<YOUR PRIVACY IDEA URL>",
    "pideaApiUser": "<YOUR PIDEA API USER(superadmin)>"
}
```

<h2>User fields</h2>

User fileds in result(default values for actual version PrivacyIdea 3.11.2, but may be the same in older versions):

* "editable": bool, # excluded from result
* "email": string,
* "givenname": string,
* "memberOf": []string,
* "mobile": string,
* "phone": string,
* "resolver": string,
* "surname": string,
* "userid": string, # excluded from result
* "username": string

<h2>Flags</h2>

* -log-dir - full path to your custom log dir(default is .logs_pi-get-users-by-realm)
* -keep-logs - number of last logs to keep, 7 by default
* -realm - Pidea realm to search users in, must exist and can't be unset

<h2>Result</h2>

Result is a *.csv file in <b>Results</b> dir, which will be created in app's root dir and file name format is:

```
result_<REALM>_<dd.mm.YYYY>.csv
```

<h2>Usage example</h2>

```
pidea-get-users-by-realm.exe -realm DOMAIN.EX.COM
```
