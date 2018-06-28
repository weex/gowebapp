# Lightning Payment Processor

Accept and forward Bitcoin lightning payments.

## REST API

* All requests via HTTP GET except where noted.
* Prepend https://<domain>/api/v1/ to generate full API endpoint URLs.
* Data returned as JSON.

### `invoice` - get a new invoice that will be processed as below
    Parameters
        `amount` - in satoshis
        `description` - what should go in memo field
    Returns
        `payment_request` - string to present to user for payment
        `payment_hash` - id used to look up payment status
        `expire_time` - seconds since the epoch in which time will exire
        `current_time` - seconds since the epoch for the processing server

### `check_invoice` - get status  of invoice previously created by `invoice`
    Parameters
        `payment_hash`- id used to lookup invoice status
    Returns
        `invoice` - invoice details including the boolean `settled`

## Future work
    Add to `invoice`, `balance`, and `send`:
        `api_key` - get this from /api_key

    For `invoice`, add parameters:
        `forward_type` - string
            accpted values:
                `immediate` - requires `address`
                `threshold` - requress `address` and `threshold`
                `hold` - funds held until onchain or lightning transfer triggered via /send
        `threshold` - optional - amount to hold in satoshis before sending onchain

### Add `account` - setup an account
    Parameters
        `email` - to receive payment notifications
        `notification_url` (optional) endpoint called with HTTP GET on payment, will have parameter r_hash filled out
    Returns
        `api_key` - id or username for the account
        `api_secret` - secret used to authenticate later calls

### Add `balance` - get your account balance
    Parameters: none
    Returns
        `amount` - amount in satoshis currently being held

### Add `send` - forward previously collected funds
    Parameters
        `destination` - Lightning payment request
        `amount` - satoshis to send, or `all` to send maximium available
    Returns
        `result` - `success` or `error`
        `fee` - amount paid to route payment
        `payment_preimage` - preimage for lightning payment
        `error_details` - reason for error
