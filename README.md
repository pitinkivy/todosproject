# clean lib
go mod tidy

# clean cache
go clean -modcache

# git bash
select default shell

# code coverage in current package


json to go

# sql database
golang gorm


# request tracking
OpenCensus
OpenTracing
OpenTelemetry

# log
uber zap

# logfile
lumberjack

# httpClient
    client := &http.Client{
        Transport: &http.Transport{
            MaxIdleConnsPerHost : 100
            MaxConnsPerHost : 100
        }
    }

disable keep alive

# kafka
golang salama