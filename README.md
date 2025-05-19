# go_pidea-get-users-by-realm
Go: get PrivacyIdea users via REST API by given realm

User fileds in result(default values for actual version PrivacyIdea 3.11.2, but may be the same in older versions):
* "editable": bool, # excludede
* "email": string,
* "givenname": string,
* "memberOf": []string,
* "mobile": string,
* "phone": string,
* "resolver": string,
* "surname": string,
* "userid": string, # excluded
* "username": string