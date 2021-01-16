// declare and exports constants
const NATS_URL = 'nats://localhost:4222'
const NATS_USER = process.env.NATS_USER
const NATS_PASS = process.env.NATS_PASS
const NATS_SUBJECT = "init-segment"
const NATS_REQUEST_TIMEOUT = 5000

exports.NATS_USER = NATS_USER
exports.NATS_PASS = NATS_PASS
exports.NATS_SUBJECT = NATS_SUBJECT
exports.NATS_REQUEST_TIMEOUT = NATS_REQUEST_TIMEOUT