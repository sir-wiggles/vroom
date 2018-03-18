# User Schema

```sql
                                         Table "public.User"
    Column     |           Type           | Collation | Nullable |              Default
---------------+--------------------------+-----------+----------+------------------------------------
 id            | integer                  |           | not null | nextval('"User_id_seq"'::regclass)
 created_at    | timestamp with time zone |           | not null | now()
 updated_at    | timestamp with time zone |           | not null | now()
 phone         | character varying(15)    |           | not null |
 email         | character varying(50)    |           |          |
 password      | character varying(255)   |           |          |
 first_name    | character varying(255)   |           |          |
 last_name     | character varying(255)   |           |          |
 userLevel     | integer                  |           | not null |
 verified      | boolean                  |           | not null |
 employee_pin  | integer                  |           |          |
 phone_plus_cc | character varying(30)    |           |          |
 profile_url   | character varying(255)   |           |          |
 referred_by   | integer                  |           |          |
 referral_code | character varying(255)   |           |          |
 vip           | boolean                  |           |          | false
 src           | character varying        |           |          |
 properties    | json                     |           |          |
Indexes:
    "User_pkey" PRIMARY KEY, btree (id)
    "user_employee_pin_unique" UNIQUE CONSTRAINT, btree (employee_pin)
    "user_phone_unique" UNIQUE CONSTRAINT, btree (phone)
```

# Users and User States

We have two types of users: `Guest` and `Member`.  Each user can have three
states: `Registering`, `VerifiedRegistering` and `Verified`. The last state
`Verified` of `Guest` and `Member` is identical, i.e. the `Guest` has converted
over to a `Member` and represented in the database as a `Member`. So really we
have five states as shown in the following table

code | user                        | userLevel | verified | password
---- | --------------------------- |---------- | -------- | --------
1    | `MemberRegistering`         | 0         |  f       |  f
3    | `MemberVerifiedRegistering` | 0         |  t       |  f
4    | `MemberVerified`            | 0         |  t       |  t
6    | `GuestRegistering`          | 99        |  f       |  t
8    | `GuestVerifiedRegistering`  | 99        |  t       |  t

The columns here (`userLevel`, `verified`, `password`) represent the columns in
the database. `password` is represented here as a boolean to indicate weather
the field is set.

The `code` column is an int representing the state the user could be in where
`userLevel` is the MSB and `password` is the LSB


# User password

Good read on how to properly [hash a user
password](https://crackstation.net/hashing-security.htm) Basically we salt the
users password and then `bcrypt` it storing the has in the database

# Authorization

A user will send an `Authorization: Bearer <token>` in the request header for
routes required authorization.

The token is a JWT payload of the user object.  We should take care to remove
all sensitive information from the user object like password which is currently
being used the node implementation

# Creating users

## As a guest

Users that complete an order without logging in to a registered i.e. the phone
number doesn't exists in the database, is considered a guest user.

Guest users have their phone number where the server will append five random
digits to the end.

Actions required to create a guest user:

* Check that we don't already have an phone number registerd to another user
```
>>>
POST /user/phone?action=check
Authorization: no
Ratelimit: yes per IP
Body: {
    phone: \d{10}
}
<<<
Response: 200 | 404
```

* Create a guest user
```
>>>
Request:
POST /user?type=guest
Authorization: no
Body: {
    email:      string,
    first_name: string,
    last_name:  string,
    password:   string, // "random" generation
    phone:      string,
    referral:   int,    // optional
}
<<<
Response: 201 409
Body:{
    token: <jwt-token>
    user: { ... }
}
```

## As a member


Actions required to create a member user:

* Check that we don't already have an phone number registerd to another user
```
>>>
POST /user/phone?action=check
Authorization: no
Ratelimit: yes per IP
Body: {
    phone: \d{10}
}
<<<
Response: 200 | 404
```

* Send a code to the phone number
```
>>>
POST /user/phone?action=send
Authorization: no
Ratelimit: yes per IP
Body: {
    phone: \d{10}
}
<<<
Response: 200
```

* The user enters in the code from above to be verified
```
>>>
POST /user/phone?action=verify
Authorization: no
Ratelimit: yes per IP
Body: {
    phone: \d{10}
}
<<<
Response: 200 | 400
Body: {
    token: hash(phone + secret) // we know the secret
}
```

* We are now permitted to create a user with the token from above as our auth
```
>>>
Request:
POST /user?type=member
Authorization: Bearer <token>
Body: {
    phone: \d{10}
}
<<<
Response: 201 | 409
Body: {
    token: <jwt-token>
    user: { ... }
}
```



-------------------------------------------------------------------------------

# /user/{phonenumber}

## HEAD



## GET
### Auth required
Reads a user information

#### Status Codes
- 200
- 404

## PUT
### Auth required
Update a user, useful for converting a guest user to a verified user

#### Status Codes
- 200

## POST
### Auth required
Creates a new user

#### Status Codes
- 201

## PATCH
### Auth required
Update a property of the user

#### Status Codes
- 200
- 404

## DELETE
### Auth required
Remove a user from the system

-------------------------------------------------------------------------------

# /user/sms?action=[send,verify]

## POST
### Throttling required
Send a sms verification code to the user
Verify a sent verification code from a user

-------------------------------------------------------------------------------

# /user/orders?state=[any|active|inactive]?active

## GET
### Auth required
Get a users orders given an order state defaulting to active.  If `any` is
used then an object will be returned with keys `active` and `inactive` having
values as an array of orders sorted from newest to oldest
