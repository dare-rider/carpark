# Carpark

## Application Architecture

The application architecture is mostly inclined towards `The Clean Architecture`, 
and used `Dependency injection` and `MVC - Model (Views)Presentors Controllers` design patterns. 
Following is the walkthrough of the application structure.
* All the core business logic resides under `app` package, like `controllers`, `models`, `tasks`, `presentors` etc
* `appmiddleware`, `cofig`, `db`, `router` packages are used to initialize application.
* `constant`, `types`, `utils` packages are used to support business logic of application.
* All the 3rd Party used plugins are already downloaded under `vendor` directory.
* `go mod` is used for package dependency management. 

## Code Overview

### Terminologies, external packages and their usage 
* `repo` packages that are used as adapter to 3rd Party, whether it's database or any 
external http call. 
* `usecase` are the wraper around `repo` which contains the business logic if any over the raw response
from `repo` and internally exposed to rest of application
* `utils/scy21` is the package(copied from some open source repo) that exports the functions to convert GeoLocation 
coordinates to lat, long
* `utils/geodist` is the package(copied from some open source repo) used to calculate the distance between 2 geo-locations. 
* `golang-migrate/migrate` package is used to execute & keep state of the database schema migrations.
* `validator.v9` package is used for request input validations.

### Application startup
* `main.go` is the entry point of application startup execution. What it does?
<pre> 
 generalConfig := config.LoadConfig(configFilePath)
 constant.InitConstants(generalConfig.MiscConfig)
 db.InitMysqlDb(generalConfig.DbConfig)
 db.InitMigrations(generalConfig.DbConfig)
 routes := router.InitRoutes(generalConfig, db.MysqlConn())

 // Starting server.
 http.ListenAndServe(":3005", routes)
</pre>

* `config.LoadConfig` is the function used to load the application config mentioned in `config/conf.yml` into a global 
config struct `generalConfig` that can be injected as dependency wherever required.
    * sample config.yml
    <pre>
    db:
        dsn: root:@tcp(127.0.0.1:3306)/carpark_dev?parseTime=true
        migration_path: "file://db/migrations/schema" # Relative
        seed_path: "db/migrations/seed" # Relative
    gov_sg_service:
        base_url: "https://api.data.gov.sg/"
    misc:
        environment: "development"
    </pre>   
* `constant.InitConstants` initializes the global constants, load environment specific variables and provide corresponding 
global getters.
* `db.InitMysqlDb` initializes the global db connection, that is then injected to all repos while routes initialization.
* `db.InitMigrations`, takes the `DbConfig` which has the migrations path, and `golang-migrate/migrate` executes all the 
pending migrations inside `db/migrations/schema` (current migrations path). `migrate` CLI is used to generate those 
migrations `up.sql` & `down.sql`. `schema_migrations` table is used to keep track of last executed migration.
* `router.InitRoutes` initializes the `controller` routes by injecting dependencies amongst their respective controller structs.
Will discuss about it in following section.

 ### `router` package, InitRoutes.
 * using `chi` router
 * initialises all the required middlewares default + custom.
 * prepare instances of all application dependencies and injecting then as required
    1. `app/models` & `app/extservices` `repo`
    2. `app/models` & `app/extservices` `usecase`, by injecting dependencies prepared in Step 1.
    3. `app/tasks`, by injecting dependencies prepared in Step 2.
    4. `app/controllers` by injecting dependencies prepared in Step 2 and Step 3.
 * Mounting controller routes - each controller has it's separate `Router` function which contains controller specific 
    routes, which is then mounted here under some default pattern, currently used as `/`
 
 ### Control flow
 <pre>
    routes ---> controller ---> usecase ---> repo ---> DB
                    |              ^           |
                    |              |           ------> GovSg
                    ----------> tasks  
 </pre> 
 
 ### API's, their payload, responses
 * GET `/tasks/carpark_upload`, used to upload carparks from `db/migrations/seed/carpark.csv`, Insert handled on 
 `create_if_not_exists` else `update` metadata.
    * Query/Path Parameters - None
    * Response Success - 200 OK
    <pre>
    {
        "status": true,
        "message": "Request Implemented Successfully"
    }
    </pre>
    
    * Response Failure - 422 Unprocessable Entity
    <pre>
    {
        "status": false,
        "message": "open db/migrations/seeds/carpark.csv: no such file or directory"
    }
    </pre>
 * GET `/tasks/carparkinfo_upload` used to pull recent carpark_infos from `https://data.gov.sg/dataset/carpark-availability` and
 persist locally, Insert handled on `create_if_not_exists` else `update` metadata.
    * Query/Path Parameters - None
    * Response Success - 200 OK
    <pre>
    {
        "status": true,
        "message": "Request Implemented Successfully"
    }
    </pre>
    
    * Response Failure - 422 Unprocessable Entity
    <pre>
    {
        "status": false,
        "message": "Get https://api.data.gov.sgs/v1/transport/carpark-availability: dial tcp: lookup api.data.gov.sgs: no such host"
    }
    </pre>
 *  GET `/carparks/nearest` to fetch the carparks sorted based on distance from input location. 
    * Algo: While uploading carparks, their distance in `KMs` from a central SG location is also populated in DB along with lat, long.
      While fetching the nearest carpark from a location L1, L1 distance is calculated from central SG location, and then sorted based 
      on absolute difference from the carpark distance from central SG location
    <pre>
    L0 - central SG location
    L0CP1, L0CP2 .... L0CPn - distance of carpark location from L0 in KMs
    L1 - current location
    Comparator: L0L1 <=> L0CPi 
    </pre>
    * Query Parameters
    <pre>
       Latitude  float64 `schema:"latitude" validate:"required"`
	   Longitude float64 `schema:"longitude" validate:"required"`
	   Page      int     `schema:"page" validate:"omitempty,min=1"`
	   PerPage   int     `schema:"per_page" validate:"omitempty,min=1"`
    </pre>
    * Response Success - 200 OK
    <pre>
    [
        {
            "address": "BLK 301-302,305-308 CLEMENTI AVENUE 4",
            "latitude": 1.201848090322625,
            "longitude": 103.88527398855013,
            "total_slots": 316,
            "available_slots": 134
        },
        {
            "address": "BLK 124/129 BEDOK NORTH STREET 2",
            "latitude": 1.375709138082525,
            "longitude": 103.89194725680792,
            "total_slots": 414,
            "available_slots": 238
        },
        {
            "address": "BLK 909A HOUGANG STREET 91",
            "latitude": 1.3200501516802696,
            "longitude": 103.9414272888625,
            "total_slots": 1052,
            "available_slots": 751
        }
    ]
    </pre>
    
    * Response Failure - 400 Bad Request
    <pre>
    {
        "status": false,
        "message": "Key: 'NearestCarparksRequest.Latitude' Error:Field validation for 'Latitude' failed on the 'required' tag"
    }
    </pre>
    * Default pagination - page: 1, per_page: 50

  