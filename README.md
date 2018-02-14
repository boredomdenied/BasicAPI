

## How to use

Get Postman
* https://www.getpostman.com/

Use the server already setup for convenience at this IP: `188.166.69.159:3100`


- Signup `/signup`

https://www.evernote.com/l/AuJW7-NMRAlC5L4WBLZXfMbqA_DGeccd9CE

You need to store the token for 1 hr session and the id in shared preferences or sqlite to get the users private endpoint

- Login `/login`

https://www.evernote.com/l/AuKRG0RFsfVNgbfcHdctuvNc4xlLGLRdH60

Inserting the correct username & password will return the JWT which again should be stored in app at 

- User Page `/users/{id}`

https://www.evernote.com/l/AuJUXaJMcv1MB7Kb8boq20L0M89v4FMf508

Must insert Bearer Token & use the correct endpoint 


### Local development

If you would like to test and develop on a local environment please fork the repo, install Golang, and install MongoDB

Get Golang
* https://golang.org/doc/install

Get MongoDB
* https://docs.mongodb.com/manual/administration/install-community/

Please ensure your paths are set for proper Golang development/testing
You should be able to `go run app.go` from inside this repo root folder

Please ensure MongoDB is properly installed
You should be able to `mongo` to enter a mongo prompt

Be sure to setup the `users` database
Enter `mongo` to get inside prompt then enter `use users` to create users database, then type `exit` to leave the prompt

### How to contribute

You can test from the server without a need to setup a local development environment. 
If you would like to develop the API it will be necessary to follow the local development steps above.
If you would like to discuss a matter and are in GwG please find me `Brandon Chapman` and DM directly. 
If you find an error or a feature missing that is necessary for the service we're creating, please create an ISSUE.
From issues created we can discuss what steps to best take next. After issues are validated, development can begin.
When code resolving issue has been developed, please create a PR that will be reviewed before integrated into master.
If PR satisfied requirements, it will be pulled and issue can be closed. 
