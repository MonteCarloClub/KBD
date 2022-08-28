namespace go api

struct Account {
    1: required string Address
    2: required i64 balance
    3: required i64 nonce
}

struct GetDataRequest{
    1: required string key
}

struct GetDataResponse{
    1: required string value
}

struct PutDataRequest{
    1: required string key
    2: required string value
}

struct PutDataResponse{
    1: required bool success
}

struct GetAccountDataRequest {
    1: required string address
}

struct GetAccountDataResponse {
    1: required string message
    2: optional Account account
}

service kanBanDatabase {
    GetDataResponse GetData(1: GetDataRequest req)
    PutDataResponse PutData(1: PutDataRequest req)
    GetAccountDataResponse GetAccountData(1:  GetAccountDataRequest req)
}