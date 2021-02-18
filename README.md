## Project Evredika
#### Here you can find implementation of light-weight service, which is used for storing user data

### What is used ? 
* Since I decided to store all data in .json format I realized that data receiving 
would be much faster if it'll be stored in service-cache.  
* Logrus is used for event-logging.  
* There are several implementations of using Database : S3 (currently is not working due to network issues w/ local-minio but you can use others S3 storages) and files in system folder.

### API
#### Postman
You can find test API calls here : https://www.getpostman.com/collections/5151a4015339afab9255
#### Routes and description


```golang
// Used data format :
type User struct {
	ID string
	Name string
}
```

Route | Method | Description | Query Params | Request Body | Response Body
--- | --- | --- | --- |---   |--- 
/api/v1/user | POST | Create user in collection; If it already exist - return error | *Empty* | Json : User struct | *Empty; Only status*
/api/v1/user | PUT | Update user in collection; If it doesn't exist - return error | *Empty* | Json : User struct | *Empty; Only status*
/api/v1/user | DELETE | Delete user in collection | ID - string (ID of user) | *Empty*| *Empty; Only status*
/api/v1/user | GET | GET user from collection | ID - string (ID of user) | *Empty*| Json : User struct
/api/v1/user/list | GET | List users from collection | skip, limit - INT; used for search| *Empty*| Json : List of User struct


### How to use
To run server just type into console : 'docker-compose up'
* If you want to up your environment using file saver in local system as your Database, write to docker-compose.yml 
  (Is already set as default) : 
  ``` yaml
      version: '3.7'
      services:
        user_server:
            build:
              dockerfile: Dockerfile.user_server
              context: .
              args:
                - APP_PKG_NAME=project-evredika
                - GOOS=linux
            ports:
              - "10080:80"
            environment:
              - PORT=:80
              - SUFFIX=.json
              - BUCKET=common_data/
              - CHANNELS_SIZE=20
              - STORAGE_TYPE=OS
  ``` 
* If you want to up your environment using S3 as your Database, write to docker-compose.yml
  (I got stuck with this version because of network issues; i cant create bucket :/ ) :
  ``` yaml
  version: '3.7'
  services:
    storage:
      image: minio/minio
      ports:
       - "19000:9000"
      environment:
        MINIO_ROOT_USER: minio123
        MINIO_ROOT_PASSWORD: minio123
      command: server /data
      networks:
       - user_server
  
    user_server:
      build:
        dockerfile: Dockerfile.user_server
        context: .
        args:
          - APP_PKG_NAME=project-evredika
          - GOOS=linux
      ports:
        - "10080:80"  
      environment:
        - PORT=:80
        - SUFFIX=.json
        - BUCKET=commondata
        - CHANNELS_SIZE=20
        - STORAGE_TYPE=S3
        - S3_ENDPOINT=storage:9000
        - S3_REGION=ru-msk
        - ACCESS_KEY_ID=minio123
        - SECRET_ACCESS_KEY=minio123
        - DISABLE_SSL=true
        - S3_FORCE_PATH_STYLE=true
      depends_on:
       - storage
      networks:
       - user_server
  
  networks:
    user_server:
      driver: bridge
  ```   
