// custom Error type for our application
class AppError extends Error {
  constructor(message, cause) {
    super(message + (cause ? ": " + cause.message : ""))
    this.stack = cause?.stack // useful because we can log error stack from cause with err.stack
  }
}

exports.AppError = AppError