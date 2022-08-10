namespace go api

struct GetAccountDataRequest {
    1: string account
}

struct GetAccountDataResponse {
    1: string message
}

service kanBanDatabase {
    GetAccountDataResponse GetAccountData(1:  GetAccountDataRequest req)
}