namespace go api

struct Account {
    1: required string Address
    2: required i64 balance
    3: required i64 nonce
}

struct GetAccountDataRequest {
    1: required string address
}

struct GetAccountDataResponse {
    1: required string message
    2: required Account account
}

service kanBanDatabase {
    GetAccountDataResponse GetAccountData(1:  GetAccountDataRequest req)
}