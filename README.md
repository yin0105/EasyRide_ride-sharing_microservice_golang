# EasyRide_ride-sharing_microservice_golang

## Run
### Roster Microservice

#### join roster
- http://localhost:8088/api/v1/driver/add/{user}/{password}/{rate}


#### change rate
- http://localhost:8088/api/v1/driver/change/{user}/{password}/{new_rate}


#### leave roster
- http://localhost:8088/api/v1/driver/remove/{user}/{password}


#### get drivers count and min rate
- http://localhost:8088/api/v1/driver/get_info/


#### get roster
- http://localhost:8088/api/v1/driver/drivers

----------------------------------------------------

### Mapping Microservice

#### get a distance of the route and "A road"
- http://localhost:8088/api/v1/mapping/origin={start_position}&destination={end_position}

----------------------------------------------------

### Surge_Pricing Microservice

#### get surge_pricing
- http://localhost:8088/api/v1/surge_pricing/{distance}/{rate}/{a_road}/{drivers_count}/{start_time_hour}/{start_time_minute}

----------------------------------------------------

### Start Microservice

#### join roster
- http://localhost:8088/api/v1/start/{start_position}/{end_position}/{start_time_hour}/{start_time_minute}
