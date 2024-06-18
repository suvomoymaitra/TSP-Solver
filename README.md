# What is TSP Solver?
TSP Solver is a web service built in Golang to solve **Travelling Salesman Problem** for a given number of points. It is made to deploy as a Lambda function in AWS.

# How to use it?
+ Open powershell and run the following commands :
  ```
  $env:GOOS = "linux"
  $env:CGO_ENABLED = "0"
  $env:GOARCH = "amd64"
  go build -o bootstrap main.go
  chmod +x bootstrap
  zip main.zip bootstrap
  ```
+ Go to AWS console and create a Lambda function
  - Select a runtime for Go
  - Choose x86_64 architecture
+ Now go to the code section and upload the main.zip file generated
+ Go to runtime settings and change handler to **main**
+ Then go to add trigger and add an API gateway
+ Select REST API and set security to open. (You can setup security as API KEY or AWS_IAM but here I have kept open for simplicity)
+ Now you can access the web service by hitting the url

# Sample Request Body
The web service expects two values
- **distance_matrix** : A distance matrix which is converted to string format
- **number_of_points** : The number of points in the distance matrix
```
{
    "distance_matrix" : "0,10,15,20,25,30,10,0,35,25,30,20,15,35,0,30,20,25,20,25,30,0,35,15,25,30,20,35,0,10,30,20,25,15,10,0",
    "number_of_points" : 6
}
```
# Sample Response Body
- ## Success response (200 OK)
    ```
    {
        "min_distance": "95.000000",
        "optimal_path": "0, 1, 3, 5, 4, 2, 0"
    }
    ```
- ## Unsuccessful (502 Internal server error)
    ```
    {
        "message": "Internal server error"
    }
    ```
