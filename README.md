# Daily Diary
Daily Diary is a RESTful backend server implemented in Golang using Clean Architecture by Uncle Bob. It helps the user to organize, preserve and track their daily activities, experiences, thoughts and feelings.

## Getting Started
The following instructions will help you to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites
* Make sure you have MySQL installed on your local machine. You can download it from  [here](https://dev.mysql.com/downloads/installer/).

### Installing
* To clone this repository, you need to have [GIT](https://git-scm.com) installed on your local machine.
* Paste the following on the command line:
```
$ git clone https://github.com/sammy9867/daily-diary.git
```

### Deployment
* After you have cloned the repository, navigate to ***.env*** file and enter your MySQL database configurations.
```
# Mysql
API_SECRET=9867sammy9867 # For JSON Web tokens, can be anything
DB_HOST=127.0.0.1
DB_DRIVER=mysql 
DB_USER=DB_USER_NAME
DB_PASSWORD=DB_USER_PASSWORD
DB_NAME=YOUR_DB
DB_PORT=3306

# Mysql Test
API_SECRET_TEST=9867sammy9867 # For JSON Web tokens, can be anything
DB_HOST_TEST=127.0.0.1
DB_DRIVER_TEST=mysql 
DB_USER_TEST=DB_USER_NAME
DB_PASSWORD_TEST=DB_USER_PASSWORD
DB_NAME_TEST=YOUR_DB_TEST
DB_PORT_TEST=3306
```
* Run the ***diary_db.sql*** file in your database workbench.
* Navigate to the folder where ***main.go*** resides and enter the following command to run the program:
```
go run main.go
```
## Running the tests
Each repository folder has a test file. Make sure you have created a separate database for testing purposes. In order to run a particular test, run the following command:
```
go test --run TestName
```


## Contributing
All pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Author
* **Samuel Menezes**
