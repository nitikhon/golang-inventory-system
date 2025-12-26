export interface User {
    username: string,
    email: string,
    password: string,
    phone: string,
    first_name: string,
    last_name: string,
    is_admin: boolean,
    refresh_token: string,
}

export interface Token {
    access_token: string,
    refresh_token: string,
}