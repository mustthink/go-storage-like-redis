# go-storage-like-redis
like redis but written by me on Golang. 
Just Key-Value storage.

## Requirements
1) Go version 1.21 or above

## Install
1) Copy repository via `git clone` or other method
2) Install necessary modules via `go mod tidy` (inside project folder)

## Run
### 1) Via Docker
   1) Set up configuration
   2) Add necessary flags to Dockerfile (if you want use non-default config)
   3) `docker build -t image-name .`
   4) `docker run -p port:port image-name`
### 2)  Manual
   1) Set up configuration
   2) `go run cmd/main.go -config=pathToConfig`

## Flags
1) `-config` - set up path to your configuration

## Configuration struct
1) `server` - server settings
   1) `host` - server host
   2) `port` - server port
   3) `read_timeout_in_ms` - timeout for read requests
   4) `write_timeout_in_ms` - timeout for write requests
   5) `auth` - optional, if you want BaseAuth for server
      1) `user` - username 
      2) `pass` - password 
2) `storage` - storage settings 
   1) `ttl_in_seconds` - default TTL for objects
   2) `max_collections_count` - max collections count 
   3) `refresh_time_in_seconds` - refreshing time for collections

## Requests
### 1) POST
if you want to create new collection or set object into collection you should use this request
#### Struct:
1) `type` - type of request: 
   1) use "collection" for creating new collection
   2) use "object" for setting object into collection
2) `collection` - name of collection which you want to create or where you want to set object. If field is empty, object will set into `default` collection
3) `key` - key for setting object into collection. If you want to create collection just leave empty. 
4) `object` - object settings. If you want to create collection just leave empty.

#### Struct of object settings
1) `data` - binary object data
2) `timeout` - object TTL
3) `deadline` - object will expire in deadline
4) `timeless` - object will never expire 
> don't use options together, priority of options is timeout -> deadline -> timeless
> If options is empty, object TTL will be default.

### 2) GET
if you want to get collection or object you should use this request

#### Struct:
1) `type` - type of request:
    1) use "collection" for get all objects from collection
    2) use "object" for get object from collection
2) `collection` - name of collection. If field is empty, object will set into `default` collection
3) `key` - key for getting object from collection. If you want to get collection just leave empty.

### 2) DELETE
if you want to delete collection or object you should use this request

#### Struct:
1) `type` - type of request:
    1) use "collection" for delete collection
    2) use "object" for delete object from collection
2) `collection` - name of collection. If field is empty, object will set into `default` collection
3) `key` - key for delete object from collection. If you want to delete collection just leave empty.

## Response 
All request has one struct of response 
### Struct:
1) `data` - binary object data. ONLY for GET requests 
2) `success` - is request successful 
3) `error` - is request has some error
   1) `message` - error message of details 
   2) `code` - http code 

## Client
simple client (not fully functional)
>See example of using client in client/example

## Tests
Project has two integration tests, because I don't have enough time for full coverage of project. 
> All tests - PASS


## TODO or what can be added in future updates
1) Parallel serving of requests (two or more server's)
2) Add metrics like Prometheus for analysis how well the service works 
3) Add `set/get/delete` requests for many objects
4) Add refreshing objects (object should refresh his data after it expire)
5) Something else...  
