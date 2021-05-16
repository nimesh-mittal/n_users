# n_users

[![Go Report Card](https://goreportcard.com/badge/github.com/nimesh-mittal/n_users)](https://goreportcard.com/report/github.com/nimesh-mittal/n_users)

## Background

Aim of this service is to provide ability to manage users. A user entity contains:

1. Basic attributes: location and contact information
2. A User can have multiple profiles
3. A User can have multiple preferences
4. A User can have profile Image

While this list covers some common usecases it is no way exhaustive.

## Requirments

User service should provide following abilities

- Ability to create and delete user
- Ability to update attributes of user
- Ability to search users by name or other attributes
- Ability to add/update profile image

## Service SLA

- Availability
User service should target 99.99% uptime
- Latency
User service should aim for less than 100 ms P95 latency
- Throughput
User service should provide 2000 QPS per node
- Freshness
User service should provide newly created users in search results immediately
- Consistency
User service should ensure only one and one user exist in the system, no dublicate users

## Architecture

![image](https://github.com/nimesh-mittal/n_users/blob/main/.github/images/user_service.png)

## Implementation

### API Interface

```go

type Interface{
  // Create new user profile
  Create(profile Profile, tenant string) (string, error)
  // Delete exisiting user profile
  Delete(profileID uuid, tenant string) (bool, error)
  // Search user profiles by keyword
  Search(query string, limit int, offset int, sortBy string, tenant string) ([]Profile, error)
  // Update user profile
  Update(profileID uuid, fieldsToUpdate map[string]interface{}, filter map[string]interface{}, tenant string) (bool, error)
  // Update user profile image
  UpdateProfileImage(profileID uuid, image []byte) (bool, error)
}

```

## Data Model

| Table Name | Description | Columns |
| ------- | ---- | ---- |
| profile | Represents user profile | (tenant_name, profile_id, profile_type, *attributes...*, *who...*)

where attributes are:

- UserID          string
- FullName        string
- Gender          string
- EmailID         string
- Mobile          string
- BirthDate       time.Time
- CityID          string
- CountryID       string
- Address         string
- Latitude        float64
- Longitude       float64
- ProfileImageURL string

Who columns are"

- Active    bool
- CreatedBy string
- CreatedAt int64
- UpdatedBy string
- UpdatedAt int64
- DeletedBy string
- DeletedAt int64

## Database choice

A close look at the API request reveals large amount to read over write requests. Data consistency is important.

Assuming number of profiles will not exceed more than few billions and number of active users at any given time will be few millions, amount of data that needs to be store in database will be couple of billion rows which can be handle by any relational database easily.

Postgres can be used to store all the user profiles. Images can be store in S3. Database table stores s3 image url as profile attribute.

## Scalability and Fault tolerance

Inorder to survive host failure, multiple instances of the service needs to be deployed behind load balancer. Load balance should detect host failure and transfer any incoming request to only healthy node. One choice is to use ngnix server to perform load balancing.

Given one instance of service can serve not more than 2000 request per second, one must deploy more than one instance to achive required throughput from the system

Load balancer should also rate limiting incoming requests to avoid single user/client penalising all other user/client due to heavy load.

Given the service is going to get more read request than write, and to avoid database contention we can deploy cache and database read replica to increase throughput.

## Functional and Load testing

Service should implement good code coverage and write functional and load test cases to maintain high engineering quality standards.

## Logging, Alerting and Monitoring

Service should also expose health end-point to accuretly monitor the health of the service.

Service should also integrate with alerting framework like new relic to accuretly alert developers on unexpected failured and downtime

Service should also integrate with logging library like zap and distributed tracing library like Jager for easy debugging

## Security

Service should validated the request using Oauth token to ensure request is coming from authentic source

## Documentation

This README file provides complete documentation. Link to any other documentation will be provided in the Reference section of this document.

## Local Development Setup

- Setup environment variables
export AWS_REGION="ap-south-1"
export AWS_SECRET_ID="<key here>"
export AWS_SECRET="<secret here>"
export POSTGRE_URL_VALUE="<url in the form postgres://userid:password@host:port/databasename>"

- Start service
```go run main.go```

- Run testcases with coverage
```go test ./... -cover```
